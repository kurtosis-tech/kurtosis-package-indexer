package crawler

import (
	"context"
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
	githubClient := github.NewClient(crawler.authenticateIfTokenProvided(ctx))

	// ticker does not immediately tick, it waits for the duration to elapse before issuing its first tick.
	// So we call it once manually to populate the store
	go crawler.doCrawlNoFailure(ctx, time.Now(), githubClient)

	ticker := time.NewTicker(crawlFrequency)
	go func() {
		for {
			select {
			case <-crawler.done:
				logrus.Info("Crawler has been closed. Returning")
				ticker.Stop()
				return
			case tickerTime := <-ticker.C:
				crawler.doCrawlNoFailure(ctx, tickerTime, githubClient)
			}
		}
	}()
	return nil
}

func (crawler *GithubCrawler) Close() error {
	close(crawler.done)
	return nil
}

func (crawler *GithubCrawler) authenticateIfTokenProvided(ctx context.Context) *http.Client {
	githubTokenMaybe := os.Getenv(githubUserTokenEnvVarName)
	if githubTokenMaybe != "" {
		tokenSource := oauth2.StaticTokenSource(
			// nolint: exhaustruct
			&oauth2.Token{
				AccessToken: githubTokenMaybe,
			},
		)
		return oauth2.NewClient(ctx, tokenSource)
	}
	return nil
}

func (crawler *GithubCrawler) doCrawlNoFailure(ctx context.Context, tickerTime time.Time, githubClient *github.Client) {
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

	allRepositories, err := getKurtosisPackageRepositories(ctx, githubClient)
	if err != nil {
		return 0, 0, stacktrace.Propagate(err, "Error search for Kurtosis package repositories on Github")
	}

	existingPackages := map[store.KurtosisPackageIdentifier]bool{}
	for _, repository := range allRepositories {
		kurtosisPackageContent, ok, err := extractKurtosisPackageContent(ctx, githubClient, repository)
		if err != nil {
			return 0, 0, stacktrace.Propagate(err, "An unexpected error occurred retrieving content for Kurtosis package repository '%s'",
				repository.GetName())
		}
		if !ok {
			logrus.Warnf("Kurtosis package repository '%s' content could not be retrieved. Does it contain a valid "+
				"'%s' and '%s' files at the root of the repository?", repository.GetFullName(), kurtosisYamlFileName, starlarkMainDotStarFileName)
			continue
		}

		kurtosisPackageApi := convertRepoContentToApi(kurtosisPackageContent)
		if err = crawler.store.UpsertPackage(kurtosisPackageApi); err != nil {
			logrus.Errorf("An error occurred updating the package '%s' in the indexer store. Stored package data"+
				"might be out of date until next refresh.", kurtosisPackageContent.Identifier)
		}
		existingPackages[kurtosisPackageContent.Identifier] = true
		repoUpdated += 1
	}

	// remove all packages from the store which don't exist anymore
	storedPackages, err := crawler.store.GetKurtosisPackages()
	if err != nil {
		return 0, 0, stacktrace.Propagate(err, "Unable to retrieve all Kurtosis packages currently stored by the indexer")
	}
	for _, storedPackage := range storedPackages {
		packageIdentifier := store.KurtosisPackageIdentifier(storedPackage.GetName())
		if _, found := existingPackages[packageIdentifier]; !found {
			repoRemoved += 1
			logrus.Infof("Removing package '%s' from the store as it doesn't seem to exist anymore", packageIdentifier)
			if err = crawler.store.DeletePackage(packageIdentifier); err != nil {
				logrus.Errorf("Unable to remove package '%s' from the store. This package data will be outdated", packageIdentifier)
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
		kurtosisPackageArgsApi...,
	)
}

func getKurtosisPackageRepositories(ctx context.Context, client *github.Client) ([]*github.Repository, error) {
	// See https://docs.github.com/en/rest/search/search?apiVersion=2022-11-28#search-repositories
	// Pagination logic taken from https://github.com/google/go-github/tree/master#pagination

	var allRepositories []*github.Repository
	repoSearchOpts := &github.SearchOptions{
		Sort:      "stars",
		Order:     "desc",
		TextMatch: false,
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: githubPageSize,
		},
	}
	for {
		repoSearchResult, rawResponse, err := client.Search.Repositories(ctx, "kurtosis-package in:topics template:false archived:false", repoSearchOpts)
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred searching for Kurtosis package repositories - page number %d", repoSearchOpts.Page)
		}
		allRepositories = append(allRepositories, repoSearchResult.Repositories...)
		if rawResponse.NextPage == 0 {
			break
		}
		repoSearchOpts.Page = rawResponse.NextPage
	}
	return allRepositories, nil
}

func extractKurtosisPackageContent(ctx context.Context, client *github.Client, githubRepository *github.Repository) (*KurtosisPackageContent, bool, error) {
	// It's not clear what to use in the `GetContents` function called below, as the `owner` field. We would prefer to
	// just pass in githubRepository.GetFullName() and that's it, but it doesn't work
	// So, to parse the owner name, we try both the `Login` and `Name` field, taking the first one that is not empty.
	// For repo owned by `kurtosis-tech` for example, it's the Login field that will be used, as the Name one is empty
	// If this is too fragile, worst case we can regexp parse `githubRepository.GetFullName()`
	var repoOwner string
	if githubRepository.GetOwner().GetLogin() != "" {
		repoOwner = githubRepository.GetOwner().GetLogin()
	} else if githubRepository.GetOwner().GetName() != "" {
		repoOwner = githubRepository.GetOwner().GetName()
	} else {
		return nil, false, stacktrace.NewError("Repository does not have an owner and is not owned by any organization, "+
			"this is unexpected:\n%v", githubRepository)
	}
	repoName := githubRepository.GetName()
	repoGetContentOpts := &github.RepositoryContentGetOptions{
		Ref: "",
	}

	kurtosisYamlFileContentResult, _, resp, err := client.Repositories.GetContents(ctx, repoOwner, repoName, kurtosisYamlFileName, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == 404 {
		logrus.Debugf("No '%s' file in repo '%s'", kurtosisYamlFileName, githubRepository.GetFullName())
		return nil, false, nil
	} else if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s'", repoName)
	}
	kurtosisPackageName, err := ParseKurtosisYaml(kurtosisYamlFileContentResult)
	if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred parsing '%s' YAML file", kurtosisYamlFileName)
	}

	starlarkMainDotStartContentResult, _, resp, err := client.Repositories.GetContents(ctx, repoOwner, repoName, starlarkMainDotStarFileName, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == 404 {
		logrus.Debugf("No '%s' file in repo '%s'", starlarkMainDotStarFileName, githubRepository.GetFullName())
		return nil, false, nil
	} else if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s'", repoName)
	}
	mainDotStarParsedContent, err := ParseStarlarkMainDoStar(starlarkMainDotStartContentResult)
	if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred parsing '%s' YAML file", kurtosisYamlFileName)
	}

	return &KurtosisPackageContent{
		Identifier:       kurtosisPackageName,
		Stars:            *githubRepository.StargazersCount,
		PackageArguments: mainDotStarParsedContent.Arguments,
	}, true, nil
}
