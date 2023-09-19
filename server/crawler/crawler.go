package crawler

import (
	"context"
	"fmt"
	"github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/api_constructors"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	crawlFrequency = 2 * time.Hour

	kurtosisYamlFileName        = "kurtosis.yml"
	starlarkMainDotStarFileName = "main.star"

	githubPageSize = 100 // that's the maximum allowed by GitHub API
)

type GithubCrawler struct {
	store store.KurtosisIndexerStore

	ctx    context.Context
	ticker *Ticker
}

func NewGithubCrawler(ctx context.Context, store store.KurtosisIndexerStore) *GithubCrawler {
	return &GithubCrawler{
		store:  store,
		ctx:    ctx,
		ticker: nil,
	}
}

func (crawler *GithubCrawler) Schedule(forceRunNow bool) error {
	if crawler.ticker != nil {
		logrus.Infof("Crawler already scheduled - stopping it first")
		crawler.ticker.Stop()
	}

	lastCrawlDatetime, err := crawler.store.GetLastCrawlDatetime(crawler.ctx)
	if err != nil {
		return stacktrace.Propagate(err, "An unexpected error occurred retrieving last crawl datetime from the store")
	}

	var initialDelay time.Duration
	if !forceRunNow && lastCrawlDatetime.Add(crawlFrequency).After(time.Now()) {
		initialDelay = time.Until(lastCrawlDatetime.Add(crawlFrequency))
	}
	logrus.Infof("Crawler starting with an initial delay of '%v' and a period of '%v'", initialDelay, crawlFrequency)

	crawler.ticker = NewTicker(initialDelay, crawlFrequency)
	go func() {
		for {
			select {
			case <-crawler.ctx.Done():
				logrus.Info("Crawler has been closed. Returning")
				crawler.ticker.Stop()
				return
			case tickerTime := <-crawler.ticker.C:
				crawler.doCrawlNoFailure(crawler.ctx, tickerTime)
			}
		}
	}()
	return nil
}

func (crawler *GithubCrawler) doCrawlNoFailure(ctx context.Context, tickerTime time.Time) {
	crawlSuccessful := false
	lastCrawlDatetime, err := crawler.store.GetLastCrawlDatetime(ctx)
	if err != nil {
		logrus.Errorf("Could not retrieve last crawl datetime from the store. This is not critical, crawling"+
			"will just continue even though the last crawl might be recent. Error was:\n%s", err.Error())
	}
	if time.Since(lastCrawlDatetime) < crawlFrequency {
		logrus.Infof("Last crawling happened as '%v' ('%v' ago), which is more recent than the crawling frequency "+
			"set to '%v', so crawling will be skipped. If the indexer is running with more than one node, it might be "+
			"that another node did the crawling in between, and this is totally expected.",
			lastCrawlDatetime, time.Since(lastCrawlDatetime), crawlFrequency)
		return
	}

	// we persist the crawl datetime before doing the crawl so that potential other nodes don't crawl at the same time
	currentCrawlDatetime := time.Now()
	if err = crawler.store.UpdateLastCrawlDatetime(crawler.ctx, currentCrawlDatetime); err != nil {
		logrus.Errorf("An error occurred persisting crawl time to database. This is not critical, but in case of "+
			"a service restart (or in a multiple nodes environment), crawling might happen more frequently than "+
			"expected. Error was was:\n%v", err.Error())
	} else {
		defer func() {
			if crawlSuccessful {
				return
			}
			logrus.Debugf("Reverting the last crawl datetime to '%s'. Current value is '%s'",
				lastCrawlDatetime, currentCrawlDatetime)
			// revert the crawl datetime to its previous value
			if err = crawler.store.UpdateLastCrawlDatetime(crawler.ctx, lastCrawlDatetime); err != nil {
				logrus.Errorf("An error occurred reverting the last crawl datetime to '%s'. Its value"+
					"will remain '%s' and no crawl will happen before '%s'. Error was was:\n%v",
					lastCrawlDatetime, currentCrawlDatetime, currentCrawlDatetime.Add(crawlFrequency), err.Error())
			}
		}()
	}

	authenticatedHttpClient, err := AuthenticatedHttpClientFromEnvVar(ctx)
	if err != nil {
		logrus.Warnf("Unable to build authenticated Github client from environment variable. It will now try "+
			"from AWS S3 bucket. Error was:\n%v", err.Error())
		authenticatedHttpClient, err = AuthenticatedHttpClientFromS3BucketContent(ctx)
		if err != nil {
			logrus.Warnf("Unable to build authenticated Github client from S3 bucket. Error was:\n%v", err.Error())
			logrus.Errorf("Unable to build authenticated Github client. This is required so that"+
				"the indexer can search Github to retrieve Kurtosis package content. Skipping indexing for now, will "+
				"retry in %v", crawlFrequency)
			return
		}
	}
	githubClient := github.NewClient(authenticatedHttpClient.Client)

	logrus.Info("Crawling Github for Kurtosis packages")
	repoUpdated, repoRemoved, err := crawler.crawlKurtosisPackages(ctx, githubClient)
	if err != nil {
		logrus.Errorf("An error occurred crawling Github for Kurtosis packages. The last crawl datetime"+
			"will be reverted to its previous value '%v'. This node will try crawling again in '%v'. "+
			"Error was:\n%s", lastCrawlDatetime, crawlFrequency, err.Error())
		crawlSuccessful = false
	} else {
		crawlSuccessful = true
	}
	logrus.Infof("Crawling finished in '%v'. Success: '%v'. Repository updated: %d - removed: %d",
		time.Since(tickerTime), crawlSuccessful, repoUpdated, repoRemoved)
}

func (crawler *GithubCrawler) crawlKurtosisPackages(ctx context.Context, githubClient *github.Client) (int, int, error) {
	var repoUpdated, repoRemoved int

	allKurtosisPackageLocators, err := getKurtosisPackageLocators(ctx, githubClient)
	if err != nil {
		return 0, 0, stacktrace.Propagate(err, "Error search for Kurtosis package repositories on Github")
	}

	existingPackages := map[store.KurtosisPackageIdentifier]bool{}
	for _, kurtosisPackageLocator := range allKurtosisPackageLocators {
		kurtosisPackageContent, ok, err := extractKurtosisPackageContent(ctx, githubClient, kurtosisPackageLocator)
		if err != nil {
			return repoUpdated, repoRemoved, stacktrace.Propagate(err, "An unexpected error occurred retrieving content for Kurtosis package repository '%s' with root '%s'",
				kurtosisPackageLocator.Repository.GetFullName(), kurtosisPackageLocator.RootPath)
		}
		if !ok {
			logrus.Warnf("Kurtosis package repository content '%s/%s' could not be retrieved. Does it contain a valid "+
				"'%s' and '%s' files at the root of the repository?",
				kurtosisPackageLocator.Repository.GetFullName(), kurtosisPackageLocator.RootPath,
				kurtosisYamlFileName, starlarkMainDotStarFileName)
			continue
		}

		kurtosisPackageApi := convertRepoContentToApi(kurtosisPackageContent)
		if err = crawler.store.UpsertPackage(crawler.ctx, kurtosisPackageApi); err != nil {
			logrus.Errorf("An error occurred updating the package '%s' in the indexer store. Stored package data"+
				"might be out of date until next refresh.", kurtosisPackageContent.Identifier)
		}
		logrus.Debugf("Added Kurtosis package at '%s/%s'",
			kurtosisPackageLocator.Repository.GetFullName(), kurtosisPackageLocator.RootPath)
		existingPackages[kurtosisPackageContent.Identifier] = true
		repoUpdated += 1
	}

	// remove all packages from the store which don't exist anymore
	storedPackages, err := crawler.store.GetKurtosisPackages(crawler.ctx)
	if err != nil {
		return repoUpdated, repoRemoved, stacktrace.Propagate(err, "Unable to retrieve all Kurtosis packages currently stored by the indexer")
	}
	for _, storedPackage := range storedPackages {
		packageIdentifier := store.KurtosisPackageIdentifier(storedPackage.GetName())
		if _, found := existingPackages[packageIdentifier]; !found {
			repoRemoved += 1
			logrus.Infof("Removing package '%s' from the store as it doesn't seem to exist anymore", packageIdentifier)
			if err = crawler.store.DeletePackage(crawler.ctx, packageIdentifier); err != nil {
				logrus.Errorf("Unable to remove package '%s' from the store. This package data will be outdated", packageIdentifier)
				continue
			}
		}
	}
	return repoUpdated, repoRemoved, nil
}

func convertRepoContentToApi(kurtosisPackageContent *KurtosisPackageContent) *generated.KurtosisPackage {
	var kurtosisPackageArgsApi []*generated.PackageArg
	for _, arg := range kurtosisPackageContent.PackageArguments {
		var convertedPackageArg *generated.PackageArg
		convertedPackageArgTypeV2Ptr, ok := convertArgumentType(arg.Type)
		if !ok {
			if arg.Type != nil {
				logrus.Warnf("Argument type '%s' for argument '%s' in package '%s' could not be parsed to the new format",
					arg.Type, arg.Name, kurtosisPackageContent.Identifier)
			}
		}

		var convertedPackageArgTypeV1Ptr *generated.ArgumentValueType
		if arg.Type != nil {
			// if arg.Type is not nil then arg.Type.Type cannot be nil as it is a required field
			convertedPackageArgTypeV1Raw := arg.Type.Type.toApiType()
			convertedPackageArgTypeV1Ptr = &convertedPackageArgTypeV1Raw
		}
		convertedPackageArg = api_constructors.NewPackageArg(
			arg.Name, arg.Description, arg.IsRequired, convertedPackageArgTypeV1Ptr, convertedPackageArgTypeV2Ptr)
		kurtosisPackageArgsApi = append(kurtosisPackageArgsApi, convertedPackageArg)
	}
	return api_constructors.NewKurtosisPackage(
		string(kurtosisPackageContent.Identifier),
		kurtosisPackageContent.Description,
		kurtosisPackageContent.Url,
		kurtosisPackageContent.Stars,
		kurtosisPackageContent.EntrypointDescription,
		kurtosisPackageContent.ReturnsDescription,
		kurtosisPackageArgsApi...,
	)
}

func convertArgumentType(argumentType *StarlarkArgumentType) (*generated.PackageArgumentType, bool) {
	if argumentType == nil {
		return nil, true
	}
	mainType := argumentType.Type.toApiType()
	packageArgumentType := &generated.PackageArgumentType{
		TopLevelType: mainType,
		InnerType_1:  nil,
		InnerType_2:  nil,
	}
	if argumentType.InnerType1 != nil {
		innerType1 := argumentType.InnerType1.toApiType()
		packageArgumentType.InnerType_1 = &innerType1
	}
	if argumentType.InnerType2 != nil {
		innerType2 := argumentType.InnerType2.toApiType()
		packageArgumentType.InnerType_2 = &innerType2
	}
	return packageArgumentType, true
}

func getKurtosisPackageLocators(ctx context.Context, client *github.Client) ([]*KurtosisPackageLocator, error) {
	// Pagination logic taken from https://github.com/google/go-github/tree/master#pagination

	var allKurtosisPackageLocator []*KurtosisPackageLocator
	searchOpts := &github.SearchOptions{
		Sort:      "stars",
		Order:     "desc",
		TextMatch: false,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: githubPageSize,
		},
	}
	searchQuery := fmt.Sprintf("filename:%s", kurtosisYamlFileName)
	for {
		repoSearchResult, rawResponse, err := client.Search.Code(ctx, searchQuery, searchOpts)
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred searching for Kurtosis package repositories - page number %d", searchOpts.Page)
		}

		for _, searchResult := range repoSearchResult.CodeResults {
			rootPath := strings.TrimSuffix(*searchResult.Path, kurtosisYamlFileName)
			if rootPath == *searchResult.Path {
				// it means nothing was trimmed. It's likely invalid
				logrus.Warnf("A search result was invalid because the path to the file did not finish with '%s'. "+
					"Resository was: '%s', path was: '%s'", kurtosisYamlFileName, *searchResult.Repository.FullName,
					*searchResult.Path)
				continue
			}
			newKurtosisPackageLocator := NewKurtosisPackageLocator(searchResult.Repository, rootPath, kurtosisYamlFileName)
			allKurtosisPackageLocator = append(allKurtosisPackageLocator, newKurtosisPackageLocator)
		}

		if rawResponse.NextPage == 0 {
			break
		}
		searchOpts.Page = rawResponse.NextPage
	}
	return allKurtosisPackageLocator, nil
}

func extractKurtosisPackageContent(ctx context.Context, client *github.Client, kurtosisPackageLocator *KurtosisPackageLocator) (*KurtosisPackageContent, bool, error) {
	// It's not clear what to use in the `GetContents` function called below, as the `owner` field. We would prefer to
	// just pass in githubRepository.GetFullName() and that's it, but it doesn't work
	// So, to parse the owner name, we try both the `Login` and `Name` field, taking the first one that is not empty.
	// For repo owned by `kurtosis-tech` for example, it's the Login field that will be used, as the Name one is empty
	// If this is too fragile, worst case we can regexp parse `githubRepository.GetFullName()`
	var repoOwner string
	if kurtosisPackageLocator.Repository.GetOwner().GetLogin() != "" {
		repoOwner = kurtosisPackageLocator.Repository.GetOwner().GetLogin()
	} else if kurtosisPackageLocator.Repository.GetOwner().GetName() != "" {
		repoOwner = kurtosisPackageLocator.Repository.GetOwner().GetName()
	} else {
		return nil, false, stacktrace.NewError("Repository does not have an owner and is not owned by any organization, "+
			"this is unexpected:\n%v", kurtosisPackageLocator.Repository)
	}
	repoName := kurtosisPackageLocator.Repository.GetName()
	repoGetContentOpts := &github.RepositoryContentGetOptions{
		Ref: "",
	}

	kurtosisYamlFilePath := fmt.Sprintf("%s%s", kurtosisPackageLocator.RootPath, kurtosisPackageLocator.KurtosisYamlFileName)
	kurtosisYamlFileContentResult, _, resp, err := client.Repositories.GetContents(ctx, repoOwner, repoName, kurtosisYamlFilePath, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == 404 {
		logrus.Debugf("No '%s' file in repo '%s'", kurtosisYamlFilePath, kurtosisPackageLocator.Repository.GetFullName())
		return nil, false, nil
	} else if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s' - file '%s'", repoName, kurtosisYamlFilePath)
	}
	kurtosisPackageName, kurtosisPackageDescription, err := ParseKurtosisYaml(kurtosisYamlFileContentResult)
	if err != nil {
		logrus.Warnf("An error occurred parsing '%s' YAML file in repository '%s'. This Kurtosis package will not be indexed. "+
			"Error was:\n%v", kurtosisYamlFilePath, kurtosisPackageLocator.Repository.GetFullName(), err.Error())
		return nil, false, nil
	}

	var kurtosisPackageUrl string
	if kurtosisYamlFileContentResult.HTMLURL != nil {
		kurtosisPackageUrl = strings.TrimSuffix(kurtosisYamlFileContentResult.GetHTMLURL(), kurtosisPackageLocator.KurtosisYamlFileName)
	} else {
		kurtosisPackageUrl = kurtosisPackageLocator.Repository.GetURL()
	}

	kurtosisMainDotStarFilePath := fmt.Sprintf("%s%s", kurtosisPackageLocator.RootPath, starlarkMainDotStarFileName)
	starlarkMainDotStartContentResult, _, resp, err := client.Repositories.GetContents(ctx, repoOwner, repoName, kurtosisMainDotStarFilePath, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == 404 {
		logrus.Debugf("No '%s' file in repo '%s'", kurtosisMainDotStarFilePath, kurtosisPackageLocator.Repository.GetFullName())
		return nil, false, nil
	} else if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s' - file '%s'", repoName, kurtosisMainDotStarFilePath)
	}
	mainDotStarParsedContent, err := ParseStarlarkMainDoStar(starlarkMainDotStartContentResult)
	if err != nil {
		logrus.Warnf("An error occurred parsing '%s' YAML file in repository '%s'. This Kurtosis package will not be indexed. "+
			"Error was:\n%v", kurtosisMainDotStarFilePath, kurtosisPackageLocator.Repository.GetFullName(), err.Error())
		return nil, false, nil
	}

	var numberOfStars uint64
	if kurtosisPackageLocator.Repository.StargazersCount != nil {
		numberOfStars = uint64(*kurtosisPackageLocator.Repository.StargazersCount)
	}
	return &KurtosisPackageContent{
		Identifier:            kurtosisPackageName,
		Description:           kurtosisPackageDescription,
		Url:                   kurtosisPackageUrl,
		EntrypointDescription: mainDotStarParsedContent.Description,
		ReturnsDescription:    mainDotStarParsedContent.ReturnDescription,
		Stars:                 numberOfStars,
		PackageArguments:      mainDotStarParsedContent.Arguments,
	}, true, nil
}
