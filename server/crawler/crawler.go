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
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	githubUserTokenEnvVarName = "GITHUB_USER_TOKEN"
	crawlFrequency            = 2 * time.Hour

	kurtosisYamlFileName        = "kurtosis.yml"
	starlarkMainDotStarFileName = "main.star"

	githubPageSize = 100 // that's the maximum allowed by GitHub API
)

type GithubCrawler struct {
	store store.KurtosisIndexerStore

	done chan bool
}

func NewGithubCrawler(store store.KurtosisIndexerStore) *GithubCrawler {
	return &GithubCrawler{
		store: store,

		done: make(chan bool),
	}
}

func (crawler *GithubCrawler) Schedule(ctx context.Context) error {
	// ticker does not immediately tick, it waits for the duration to elapse before issuing its first tick.
	// So we call it once manually to populate the store
	go crawler.doCrawlNoFailure(ctx, time.Now())

	ticker := time.NewTicker(crawlFrequency)
	go func() {
		for {
			select {
			case <-crawler.done:
				logrus.Info("Crawler has been closed. Returning")
				ticker.Stop()
				return
			case tickerTime := <-ticker.C:
				crawler.doCrawlNoFailure(ctx, tickerTime)
			}
		}
	}()
	return nil
}

func (crawler *GithubCrawler) Close() error {
	close(crawler.done)
	return nil
}

func (crawler *GithubCrawler) mustBeAuthenticated(ctx context.Context) (*http.Client, error) {
	githubTokenMaybe := os.Getenv(githubUserTokenEnvVarName)
	if githubTokenMaybe != "" {
		tokenSource := oauth2.StaticTokenSource(
			// nolint: exhaustruct
			&oauth2.Token{
				AccessToken: githubTokenMaybe,
			},
		)
		return oauth2.NewClient(ctx, tokenSource), nil
	}
	return nil, stacktrace.NewError("Environment variable '%s' was empty", githubUserTokenEnvVarName)
}

func (crawler *GithubCrawler) doCrawlNoFailure(ctx context.Context, tickerTime time.Time) {
	authenticatedHttpClient, err := crawler.mustBeAuthenticated(ctx)
	if err != nil {
		logrus.Errorf("Unable to build authenticated Github client. This is required so that"+
			"the indexer can search Github to retrieve Kurtosis package content. Error was:\n%v", err.Error())
		return
	}
	githubClient := github.NewClient(authenticatedHttpClient)

	logrus.Info("Crawling Github for Kurtosis packages")
	repoUpdated, repoRemoved, err := crawler.crawlKurtosisPackages(ctx, githubClient)
	if err != nil {
		logrus.Errorf("An error occurred crawling Github for Kurtosis packages. Will try again in '%v'. "+
			"Error was:\n%v", crawlFrequency, err.Error())
	}
	logrus.Infof("Crawling finished in %v. Repository updated: %d - removed: %d",
		time.Since(tickerTime), repoUpdated, repoRemoved)
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
		if err = crawler.store.UpsertPackage(kurtosisPackageApi); err != nil {
			logrus.Errorf("An error occurred updating the package '%s' in the indexer store. Stored package data"+
				"might be out of date until next refresh.", kurtosisPackageContent.Identifier)
		}
		logrus.Debugf("Added Kurtosis package at '%s/%s'",
			kurtosisPackageLocator.Repository.GetFullName(), kurtosisPackageLocator.RootPath)
		existingPackages[kurtosisPackageContent.Identifier] = true
		repoUpdated += 1
	}

	// remove all packages from the store which don't exist anymore
	storedPackages, err := crawler.store.GetKurtosisPackages()
	if err != nil {
		return repoUpdated, repoRemoved, stacktrace.Propagate(err, "Unable to retrieve all Kurtosis packages currently stored by the indexer")
	}
	for _, storedPackage := range storedPackages {
		packageIdentifier := store.KurtosisPackageIdentifier(storedPackage.GetName())
		if _, found := existingPackages[packageIdentifier]; !found {
			repoRemoved += 1
			logrus.Infof("Removing package '%s' from the store as it doesn't seem to exist anymore", packageIdentifier)
			if err = crawler.store.DeletePackage(packageIdentifier); err != nil {
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
		convertedPackageArgType, ok := generated.PackageArgType_value[strings.ToUpper(arg.Type)]
		if ok {
			convertedPackageArg = api_constructors.NewPackageArg(
				arg.Name,
				arg.IsRequired,
				generated.PackageArgType(convertedPackageArgType))
		} else {
			if arg.Type != "" {
				logrus.Warnf("Argument type '%s' for argument '%s' in package '%s' is no known and will be dropped",
					arg.Type, arg.Name, kurtosisPackageContent.Identifier)
			}
			convertedPackageArg = api_constructors.NewUntypedPackageArg(arg.Name, arg.IsRequired)
		}
		kurtosisPackageArgsApi = append(kurtosisPackageArgsApi, convertedPackageArg)
	}
	return api_constructors.NewKurtosisPackage(
		string(kurtosisPackageContent.Identifier),
		kurtosisPackageContent.Url,
		kurtosisPackageContent.Stars,
		kurtosisPackageArgsApi...,
	)
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
	kurtosisPackageName, err := ParseKurtosisYaml(kurtosisYamlFileContentResult)
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
		Identifier:       kurtosisPackageName,
		Url:              kurtosisPackageUrl,
		Stars:            numberOfStars,
		PackageArguments: mainDotStarParsedContent.Arguments,
	}, true, nil
}
