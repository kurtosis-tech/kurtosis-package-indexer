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
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"time"
)

const (
	crawlFrequency      = 5 * time.Minute
	crawlIntervalBuffer = 15 * time.Second

	kurtosisYamlFileName        = "kurtosis.yml"
	starlarkMainDotStarFileName = "main.star"

	githubPageSize = 100 // that's the maximum allowed by GitHub API

	githubUrl = "github.com"

	successfulParsingText = "Parsed package content successfully"
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

	// Add a small buffer to avoid false positive when checking if the new crawl is sooner than the crawl frequency interval
	lastCrawlDatetimeWithBuffer := lastCrawlDatetime.Add(-crawlIntervalBuffer)
	if time.Since(lastCrawlDatetimeWithBuffer) < crawlFrequency {
		logrus.Infof("Last crawling happened as '%v' ('%v' ago), which is more recent than the crawling frequency "+
			"set to '%v', so crawling will be skipped. If the indexer is running with more than one node, it might be "+
			"that another node did the crawling in between, and this is totally expected.",
			lastCrawlDatetimeWithBuffer, time.Since(lastCrawlDatetimeWithBuffer), crawlFrequency)
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

	githubClient, err := CreateGithubClient(ctx)
	if err != nil {
		logrus.Error("An error occurred while creating the github client. Aborting", err)
		return
	}

	logrus.Info("Crawling Github for Kurtosis packages")
	packageUpdated, packageAdded, packageRemoved, err := crawler.crawlKurtosisPackages(ctx, githubClient)
	if err != nil {
		logrus.Errorf("An error occurred crawling Github for Kurtosis packages. The last crawl datetime"+
			"will be reverted to its previous value '%v'. This node will try crawling again in '%v'. "+
			"Error was:\n%s", lastCrawlDatetime, crawlFrequency, err.Error())
		crawlSuccessful = false
	} else {
		crawlSuccessful = true
	}
	logrus.Infof("Crawling finished in %v. Success: '%v'. Repository updated: %d - added: %d - removed: %d",
		time.Since(tickerTime), crawlSuccessful, packageUpdated, packageAdded, packageRemoved)
}

func CreateGithubClient(ctx context.Context) (*github.Client, error) {
	authenticatedHttpClient, err := AuthenticatedHttpClientFromEnvVar(ctx)
	if err != nil {
		logrus.Warnf("Unable to build authenticated Github client from environment variable. It will now try "+
			"from AWS S3 bucket. Error was:\n%v", err.Error())
		authenticatedHttpClient, err = AuthenticatedHttpClientFromS3BucketContent(ctx)
		if err != nil {
			logrus.Warnf("Unable to build authenticated Github client from S3 bucket. Error was:\n%v", err.Error())
			return nil, stacktrace.NewError("Unable to build authenticated Github client. This is required so that"+
				"the indexer can search Github to retrieve Kurtosis package content. Skipping indexing for now, will "+
				"retry in %v", crawlFrequency)
		}
	}
	githubClient := github.NewClient(authenticatedHttpClient.Client)
	return githubClient, nil
}

func (crawler *GithubCrawler) crawlKurtosisPackages(
	ctx context.Context,
	githubClient *github.Client,
) (int, int, int, error) {
	kurtosisPackageUpdated := map[string]bool{}
	kurtosisPackageRemoved := map[string]bool{}
	kurtosisPackageAdded := map[string]bool{}

	// First, we refresh the packages that are currently stored, and remove them if they don't exist anymore
	currentlyStoredPackages, err := crawler.store.GetKurtosisPackages(ctx)
	if err != nil {
		return 0, 0, 0, stacktrace.Propagate(err, "An error occurred retrieving currently stored packages")
	}
	for _, storedPackage := range currentlyStoredPackages {
		apiRepositoryMetadata := storedPackage.GetRepositoryMetadata()
		kurtosisPackageMetadata := NewPackageRepositoryMetadata(
			apiRepositoryMetadata.GetOwner(),
			apiRepositoryMetadata.GetName(),
			apiRepositoryMetadata.GetRootPath(),
			kurtosisYamlFileName,
			storedPackage.GetStars(), // this is optional here as it will be updated extractKurtosisPackageContent below
		)
		packageRepositoryLocator := kurtosisPackageMetadata.GetLocator()
		logrus.Debugf("Trying to update content of package '%s'", packageRepositoryLocator) // TODO: remove log line
		kurtosisPackageContent, ok, err := extractKurtosisPackageContent(ctx, githubClient, kurtosisPackageMetadata)
		if err != nil {
			return 0, 0, 0, stacktrace.Propagate(err, "An unexpected error occurred retrieving content for Kurtosis package repository '%s'",
				packageRepositoryLocator)
		}
		if !ok {
			logrus.Warnf("Kurtosis package repository content '%s' could not be retrieved. It will be "+
				"removed from the store", packageRepositoryLocator)
			kurtosisPackageRemoved[packageRepositoryLocator] = true
			if err := crawler.store.DeletePackage(ctx, packageRepositoryLocator); err != nil {
				logrus.Warnf("An error occurred removing package '%s' from repository '%s' from the store. It will remain but the "+
					"package doesn't exist anymore.", kurtosisPackageContent.Name, packageRepositoryLocator)
			}
		} else {
			logrus.Debugf("Updating package content '%s' from repository '%s'",
				kurtosisPackageContent.Name, packageRepositoryLocator)
			kurtosisPackageUpdated[packageRepositoryLocator] = true
			if err := crawler.store.UpsertPackage(ctx, packageRepositoryLocator, convertRepoContentToApi(kurtosisPackageContent)); err != nil {
				logrus.Errorf("An error occurred updating the package '%s' in the indexer store. Stored package data"+
					"might be out of date until next refresh.", packageRepositoryLocator)
			}
		}
	}

	logrus.Debugf("Finished updating currently stored packages. Going to search for potential new packages now")

	// Then we search for potential new packages
	allKurtosisPackageLocatorsFromSearch, err := searchForKurtosisPackageRepositories(ctx, githubClient)
	if err != nil {
		return 0, 0, 0, stacktrace.Propagate(err, "Error search for Kurtosis package repositories on Github")
	}
	for _, kurtosisPackageMetadata := range allKurtosisPackageLocatorsFromSearch {
		packageRepositoryLocator := kurtosisPackageMetadata.GetLocator()
		if _, found := kurtosisPackageUpdated[packageRepositoryLocator]; found {
			// package was already stored prior to this crawling. Its content has been refreshed above. Skipping it here
			continue
		}

		kurtosisPackageContent, ok, err := extractKurtosisPackageContent(ctx, githubClient, kurtosisPackageMetadata)
		if err != nil {
			logrus.Warnf("An error occurred parsing content for Kurtosis package repository '%s'", packageRepositoryLocator)
		}
		if !ok {
			logrus.Warnf("Kurtosis package repository content '%s' could not be retrieved as it was invalid.",
				packageRepositoryLocator)
			continue
		}

		kurtosisPackageApi := convertRepoContentToApi(kurtosisPackageContent)
		if err = crawler.store.UpsertPackage(crawler.ctx, packageRepositoryLocator, kurtosisPackageApi); err != nil {
			logrus.Errorf("An error occurred updating the package '%s' in the indexer store. Stored package data"+
				"might be out of date until next refresh.", kurtosisPackageContent.Name)
		}
		kurtosisPackageAdded[packageRepositoryLocator] = true
		logrus.Debugf("Added new Kurtosis package to the store: '%s'", packageRepositoryLocator)
	}
	return len(kurtosisPackageUpdated), len(kurtosisPackageAdded), len(kurtosisPackageRemoved), nil
}

func ReadPackage(
	ctx context.Context,
	githubClient *github.Client,
	apiRepositoryMetadata *generated.PackageRepository,
) (*generated.KurtosisPackage, error) {
	kurtosisPackageMetadata := NewPackageRepositoryMetadata(
		apiRepositoryMetadata.GetOwner(),
		apiRepositoryMetadata.GetName(),
		apiRepositoryMetadata.GetRootPath(),
		kurtosisYamlFileName,
		0, // We don't know (or care) what the star count is
	)
	packageRepositoryLocator := kurtosisPackageMetadata.GetLocator()
	kurtosisPackageContent, ok, err := extractKurtosisPackageContent(ctx, githubClient, kurtosisPackageMetadata)
	if !ok {
		parsingError := stacktrace.NewError(fmt.Sprintf("Kurtosis package repository content '%s' could not be retrieved as it was invalid.", packageRepositoryLocator), err)
		logrus.Warn(parsingError)
		return nil, parsingError
	}
	if err != nil {
		return nil, stacktrace.Propagate(err, "an unexpected error occurred retrieving content for Kurtosis package repository '%s'",
			packageRepositoryLocator)
	}

	kurtosisPackageApi := convertRepoContentToApi(kurtosisPackageContent)
	return kurtosisPackageApi, nil
}

func convertRepoContentToApi(kurtosisPackageContent *KurtosisPackageContent) *generated.KurtosisPackage {
	var kurtosisPackageArgsApi []*generated.PackageArg
	for _, arg := range kurtosisPackageContent.PackageArguments {
		var convertedPackageArg *generated.PackageArg
		convertedPackageArgTypeV2Ptr, ok := convertArgumentType(arg.Type)
		if !ok {
			if arg.Type != nil {
				logrus.Warnf("Argument type '%s' for argument '%s' in package '%s' could not be parsed to the new format",
					arg.Type, arg.Name, kurtosisPackageContent.Name)
			}
		}
		convertedPackageArg = api_constructors.NewPackageArg(
			arg.Name, arg.Description, arg.IsRequired, convertedPackageArgTypeV2Ptr, arg.DefaultValue)
		kurtosisPackageArgsApi = append(kurtosisPackageArgsApi, convertedPackageArg)
	}
	packageRepository := api_constructors.NewPackageRepository(
		githubUrl,
		kurtosisPackageContent.RepositoryMetadata.Owner,
		kurtosisPackageContent.RepositoryMetadata.Name,
		kurtosisPackageContent.RepositoryMetadata.RootPath,
	)

	return api_constructors.NewKurtosisPackage(
		kurtosisPackageContent.Name,
		kurtosisPackageContent.Description,
		packageRepository,
		kurtosisPackageContent.RepositoryMetadata.Stars,
		kurtosisPackageContent.EntrypointDescription,
		kurtosisPackageContent.ReturnsDescription,
		kurtosisPackageContent.ParsingResult,
		kurtosisPackageContent.ParsingTime,
		kurtosisPackageContent.Version,
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

func searchForKurtosisPackageRepositories(ctx context.Context, client *github.Client) ([]*PackageRepositoryMetadata, error) {
	// Pagination logic taken from https://github.com/google/go-github/tree/master#pagination

	var allPackageRepositoryMetadatas []*PackageRepositoryMetadata
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

			// It's not clear what to use in the `GetContents` function called below, as the `owner` field. We would prefer to
			// just pass in githubRepository.GetFullName() and that's it, but it doesn't work
			// So, to parse the owner name, we try both the `Login` and `Name` field, taking the first one that is not empty.
			// For repo owned by `kurtosis-tech` for example, it's the Login field that will be used, as the Name one is empty
			// If this is too fragile, worst case we can regexp parse `githubRepository.GetFullName()`
			repository := searchResult.Repository
			var repoOwner string
			if repository.GetOwner().GetLogin() != "" {
				repoOwner = repository.GetOwner().GetLogin()
			} else if repository.GetOwner().GetName() != "" {
				repoOwner = repository.GetOwner().GetName()
			} else {
				logrus.Warnf("Search result '%s' was invalid b/c the Github repository has no owner",
					repository.GetFullName())
				continue
			}

			var numberOfStars uint64
			if repository.GetStargazersCount() > 0 {
				numberOfStars = uint64(repository.GetStargazersCount())
			}

			newPackageRepositoryMetadata := NewPackageRepositoryMetadata(repoOwner, repository.GetName(), rootPath, kurtosisYamlFileName, numberOfStars)
			allPackageRepositoryMetadatas = append(allPackageRepositoryMetadatas, newPackageRepositoryMetadata)
		}

		if rawResponse.NextPage == 0 {
			break
		}
		searchOpts.Page = rawResponse.NextPage
	}
	return allPackageRepositoryMetadatas, nil
}

func extractKurtosisPackageContent(
	ctx context.Context,
	client *github.Client,
	packageRepositoryMetadata *PackageRepositoryMetadata,
) (*KurtosisPackageContent, bool, error) {
	repositoryFullName := fmt.Sprintf("%s/%s/%s", packageRepositoryMetadata.Owner, packageRepositoryMetadata.Name, packageRepositoryMetadata.RootPath)

	repoGetContentOpts := &github.RepositoryContentGetOptions{
		Ref: "",
	}

	nowAsUTC := getTimeProtobufInUTC()

	kurtosisYamlFilePath := fmt.Sprintf("%s%s", packageRepositoryMetadata.RootPath, packageRepositoryMetadata.KurtosisYamlFileName)
	kurtosisYamlFileContentResult, _, resp, err := client.Repositories.GetContents(ctx, packageRepositoryMetadata.Owner, packageRepositoryMetadata.Name, kurtosisYamlFilePath, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == 404 {
		logrus.Debugf("No '%s' file in repo '%s'", kurtosisYamlFilePath, repositoryFullName)
		return nil, false, nil
	} else if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s' - file '%s'", repositoryFullName, kurtosisYamlFilePath)
	}
	kurtosisPackageName, kurtosisPackageDescription, commitSHA, err := ParseKurtosisYaml(kurtosisYamlFileContentResult)
	if err != nil {
		logrus.Warnf("An error occurred parsing '%s' YAML file in repository '%s'"+
			"Error was:\n%v", kurtosisYamlFilePath, repositoryFullName, err.Error())
		return nil, false, err
	}

	// we check that the name set in kurtosis.yml matches the location on Github. If not, we exclude it from the
	// indexer b/c users won't be able to run it with `kurtosis run <PACKAGE_NAME>`
	expectedPackageName := normalizeName(
		fmt.Sprintf(
			"%s/%s/%s/%s",
			githubUrl,
			packageRepositoryMetadata.Owner,
			packageRepositoryMetadata.Name,
			packageRepositoryMetadata.RootPath,
		),
	)
	if normalizeName(kurtosisPackageName) != expectedPackageName {
		nameError := stacktrace.NewError(
			"The package '%s' is invalid because the name set in its '%s' doesn't match its Github URL (name set: '%s' - expected: '%s').",
			repositoryFullName,
			kurtosisYamlFileName,
			kurtosisPackageName,
			expectedPackageName,
		)
		logrus.Warn("package name must match the repo name", nameError)
		return nil, false, nameError
	}

	kurtosisMainDotStarFilePath := fmt.Sprintf("%s%s", packageRepositoryMetadata.RootPath, starlarkMainDotStarFileName)
	starlarkMainDotStartContentResult, _, resp, err := client.Repositories.GetContents(ctx, packageRepositoryMetadata.Owner, packageRepositoryMetadata.Name, kurtosisMainDotStarFilePath, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == 404 {
		logrus.Debugf("No '%s' file in repo '%s'", kurtosisMainDotStarFilePath, repositoryFullName)
		return nil, false, nil
	} else if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s' - file '%s'", repositoryFullName, kurtosisMainDotStarFilePath)
	}
	mainDotStarParsedContent, err := ParseStarlarkMainDotStar(starlarkMainDotStartContentResult)
	if err != nil {
		logrus.Warnf("An error occurred parsing '%s' YAML file in repository '%s'. This Kurtosis package will not be indexed. "+
			"Error was:\n%v", kurtosisMainDotStarFilePath, repositoryFullName, err.Error())
		return nil, false, nil
	}

	return NewKurtosisPackageContent(
		packageRepositoryMetadata,
		kurtosisPackageName,
		kurtosisPackageDescription,
		mainDotStarParsedContent.Description,
		mainDotStarParsedContent.ReturnDescription,
		successfulParsingText,
		nowAsUTC,
		commitSHA,
		mainDotStarParsedContent.Arguments...,
	), true, nil
}

func getTimeInUTC() time.Time {
	now := time.Now()
	return now.UTC()
}

func getTimeProtobufInUTC() *timestamppb.Timestamp {
	return timestamppb.New(getTimeInUTC())
}

func normalizeName(name string) string {
	return strings.ToLower(strings.Trim(name, " /"))
}
