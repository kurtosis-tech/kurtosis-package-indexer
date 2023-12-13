package crawler

import (
	"context"
	"fmt"
	"github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/api_constructors"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/ticker"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	crawlFrequency      = 5 * time.Minute
	crawlIntervalBuffer = 15 * time.Second

	defaultKurtosisYamlFilename = "kurtosis.yml"
	starlarkMainDotStarFileName = "main.star"

	githubPageSize = 100 // that's the maximum allowed by GitHub API

	githubUrl = "github.com"

	successfulParsingText = "Parsed package content successfully"

	noStartsSet = 0

	noDefaultBranchSet = ""

	onlyOneCommit = 1

	kurtosisPackageIconImgName = "kurtosis-package-icon.png"
)

var (
	zeroValueTime = time.Time{}
)

type GithubCrawler struct {
	store  store.KurtosisIndexerStore
	ctx    context.Context
	ticker *ticker.Ticker
}

func NewGitHubCrawler(ctx context.Context, store store.KurtosisIndexerStore) (*GithubCrawler, error) {

	newCrawler := &GithubCrawler{
		store:  store,
		ctx:    ctx,
		ticker: nil,
	}

	return newCrawler, nil
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

	crawler.ticker = ticker.NewTicker(initialDelay, crawlFrequency)
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

func (crawler *GithubCrawler) ReadPackage(
	ctx context.Context,
	apiRepositoryMetadata *generated.PackageRepository,
) (*generated.KurtosisPackage, error) {
	kurtosisPackageMetadata := NewPackageRepositoryMetadata(
		apiRepositoryMetadata.GetOwner(),
		apiRepositoryMetadata.GetName(),
		apiRepositoryMetadata.GetRootPath(),
		defaultKurtosisYamlFilename,
		noStartsSet,        // it will be updated extractKurtosisPackageContent below
		zeroValueTime,      // it will be updated extractKurtosisPackageContent below
		noDefaultBranchSet, // it will be updated extractKurtosisPackageContent below
	)

	githubClient, err := createGithubClient(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred while creating the github client")
	}

	packageRepositoryLocator := kurtosisPackageMetadata.GetLocator()

	packagesRunCount := crawler.getPackagesRunCountNotFailure(ctx)
	kurtosisPackageContent, packageFound, err := extractKurtosisPackageContent(ctx, githubClient, kurtosisPackageMetadata, packagesRunCount)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred extracting content for Kurtosis package repository '%s'", packageRepositoryLocator)
	}
	if !packageFound {
		logrus.Warn("No Kurtosis package found.") // don't want to log provided package repository locator bc it's a security risk (e.g. malicious data)
		return nil, stacktrace.NewError("No Kurtosis package found. Ensure that a package exists at '%v' with valid '%v' and '%v' files.", packageRepositoryLocator, defaultKurtosisYamlFilename, starlarkMainDotStarFileName)
	}

	kurtosisPackageApi := convertRepoContentToApi(kurtosisPackageContent)
	return kurtosisPackageApi, nil
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
					"will remain '%s' and no crawl will happen before '%s'. Error was:\n%v",
					lastCrawlDatetime, currentCrawlDatetime, currentCrawlDatetime.Add(crawlFrequency), err.Error())
			}
		}()
	}

	githubClient, err := createGithubClient(ctx)
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

// TODO review this because it could be optimized, right now the ReadPackages method is calling it
// TODO and creating a new client on each request
func createGithubClient(ctx context.Context) (*github.Client, error) {
	authenticatedHttpClient, err := AuthenticatedHttpClientFromEnvVar(ctx)
	if err != nil {
		logrus.Warnf("Unable to build authenticated Github client from environment variable. It will now try "+
			"from AWS S3 bucket. Error was:\n%v", err.Error())
		authenticatedHttpClient, err = AuthenticatedHttpClientFromS3BucketContent(ctx)
		if err != nil {
			logrus.Warnf("Unable to build authenticated Github client from S3 bucket. Error was:\n%v", err.Error())
			return nil, stacktrace.NewError("Unable to build authenticated Github client. This is required so that "+
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
	// First, we refresh the packages that are currently stored, and remove them if they don't exist anymore
	kurtosisPackageUpdated, kurtosisPackageRemoved, err := crawler.updateOrDeleteStoredPackages(ctx, githubClient)
	if err != nil {
		return 0, 0, 0, stacktrace.Propagate(err, "an error occurred updating or deleting the stored packages")
	}

	// Then we search for potential new packages and add they into the store
	kurtosisPackageAdded, err := crawler.addNewPackages(ctx, githubClient, kurtosisPackageUpdated)
	if err != nil {
		return 0, 0, 0,
			stacktrace.Propagate(err, "an error occurred adding new packages")
	}
	return len(kurtosisPackageUpdated), len(kurtosisPackageAdded), len(kurtosisPackageRemoved), nil
}

func (crawler *GithubCrawler) updateOrDeleteStoredPackages(ctx context.Context, githubClient *github.Client) (map[string]bool, map[string]bool, error) {
	kurtosisPackageUpdated := map[string]bool{}
	kurtosisPackageRemoved := map[string]bool{}

	logrus.Debugf("Updating currently stored packages...")
	currentlyStoredPackages, err := crawler.store.GetKurtosisPackages(ctx)
	if err != nil {
		return nil, nil, stacktrace.Propagate(err, "An error occurred retrieving currently stored packages")
	}
	for _, storedPackage := range currentlyStoredPackages {
		apiRepositoryMetadata := storedPackage.GetRepositoryMetadata()
		kurtosisPackageMetadata := NewPackageRepositoryMetadata(
			apiRepositoryMetadata.GetOwner(),
			apiRepositoryMetadata.GetName(),
			apiRepositoryMetadata.GetRootPath(),
			defaultKurtosisYamlFilename,
			storedPackage.GetStars(),                           // this is optional here as it will be updated extractKurtosisPackageContent below
			apiRepositoryMetadata.GetLastCommitTime().AsTime(), // this is optional here as it will be updated extractKurtosisPackageContent below
			apiRepositoryMetadata.GetDefaultBranch(),           // this is optional here as it will be updated extractKurtosisPackageContent below
		)
		packageRepositoryLocator := kurtosisPackageMetadata.GetLocator()
		logrus.Debugf("Trying to update content of package '%s'", packageRepositoryLocator) // TODO: remove log line

		packagesRunCount := crawler.getPackagesRunCountNotFailure(ctx)
		kurtosisPackageContent, packageFound, err := extractKurtosisPackageContent(ctx, githubClient, kurtosisPackageMetadata, packagesRunCount)
		if err != nil {
			return nil, nil, stacktrace.Propagate(err, "An error occurred extracting content for Kurtosis package repository '%s'", packageRepositoryLocator)
		}
		if !packageFound {
			logrus.Warnf("Kurtosis package repository content '%s' could not be retrieved. It will be removed from the store", packageRepositoryLocator)
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
	logrus.Debugf("...finished updating stored packages.")
	return kurtosisPackageUpdated, kurtosisPackageRemoved, nil
}

func (crawler *GithubCrawler) addNewPackages(ctx context.Context, githubClient *github.Client, kurtosisPackageUpdated map[string]bool) (map[string]bool, error) {
	logrus.Debugf("Going to read the package catalog for potential new packages now...")
	allKurtosisPackagesRepositoryMetadata, err := getPackagesRepositoryMetadata(ctx, githubClient)
	if err != nil {
		return nil, stacktrace.Propagate(err, "and error occurred getting the packages repository metadata")
	}

	kurtosisPackageAdded := map[string]bool{}
	for _, kurtosisPackageMetadata := range allKurtosisPackagesRepositoryMetadata {
		packageRepositoryLocator := kurtosisPackageMetadata.GetLocator()
		if _, found := kurtosisPackageUpdated[packageRepositoryLocator]; found {
			// package was already stored prior to this reading. Its content has been refreshed above. Skipping it here
			continue
		}

		packagesRunCount := crawler.getPackagesRunCountNotFailure(ctx)
		kurtosisPackageContent, packageFound, err := extractKurtosisPackageContent(ctx, githubClient, kurtosisPackageMetadata, packagesRunCount)
		if err != nil {
			logrus.Warnf("An error occurred extracting content for Kurtosis package at '%s'", packageRepositoryLocator)
		}
		if !packageFound {
			logrus.Warnf("Kurtosis package repository content '%s' could not be retrieved as it was invalid.", packageRepositoryLocator)
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
	logrus.Debugf("...added '%v' new packages in the store", len(kurtosisPackageAdded))

	return kurtosisPackageAdded, nil
}

func (crawler *GithubCrawler) getPackagesRunCountNotFailure(ctx context.Context) types.PackagesRunCount {
	packagesRunCount, err := crawler.store.GetPackagesRunCount(ctx)
	if err != nil {
		logrus.Errorf("an error occurred when trying to get the packages run count from the store "+
			"it won't stop the crawler because its prefered to have a degraded experience in the indexer (no returning the packages run metrics) "+
			"than not running it at all. Check if necessary the metrics storage env vars are set, this is the most probably failure. Error was:\n%s", err.Error())
	}
	return packagesRunCount
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
		kurtosisPackageContent.RepositoryMetadata.LastCommitTime,
		kurtosisPackageContent.RepositoryMetadata.DefaultBranch,
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
		kurtosisPackageContent.IconURL,
		kurtosisPackageContent.RunCount,
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

func getPackagesRepositoryMetadata(ctx context.Context, client *github.Client) ([]*PackageRepositoryMetadata, error) {

	var allPackageRepositoryMetadatas []*PackageRepositoryMetadata
	packagesCatalog, err := catalog.GetPackagesCatalog()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the packages catalog")
	}

	for _, packageData := range packagesCatalog {
		repositoryOwner := packageData.GetRepositoryOwner()
		repositoryName := packageData.GetRepositoryName()
		repositoryPackageRootPath := packageData.GetRepositoryPackageRootPath()

		repository, _, err := client.Repositories.Get(ctx, repositoryOwner, repositoryName)
		if err != nil {
			return nil, stacktrace.Propagate(err, "an error occurred getting repository '%s' owned by '%s'", repositoryName, repositoryOwner)
		}

		var numberOfStars uint64
		if repository.GetStargazersCount() > 0 {
			numberOfStars = uint64(repository.GetStargazersCount())
		}

		defaultBranch := noDefaultBranchSet
		if repository.GetDefaultBranch() != noDefaultBranchSet {
			defaultBranch = repository.GetDefaultBranch()
		}

		newPackageRepositoryMetadata := NewPackageRepositoryMetadata(repositoryOwner, repositoryName, repositoryPackageRootPath, defaultKurtosisYamlFilename, numberOfStars, zeroValueTime, defaultBranch)
		allPackageRepositoryMetadatas = append(allPackageRepositoryMetadatas, newPackageRepositoryMetadata)
	}

	return allPackageRepositoryMetadatas, nil
}

// If the package doesn't exist, returns an empty package content, false, and empty error
// If package exists, but error occurred while extracting package content returns empty package content, false, with an error
func extractKurtosisPackageContent(
	ctx context.Context,
	client *github.Client,
	packageRepositoryMetadata *PackageRepositoryMetadata,
	packagesRunCount types.PackagesRunCount,
) (*KurtosisPackageContent, bool, error) {
	repositoryOwner := packageRepositoryMetadata.Owner
	repositoryName := packageRepositoryMetadata.Name
	packageRootPath := packageRepositoryMetadata.RootPath
	kurtosisYamlFilename := packageRepositoryMetadata.KurtosisYamlFileName

	repositoryFullName := fmt.Sprintf("%s/%s/%s", repositoryOwner, repositoryName, packageRootPath)

	repoGetContentOpts := &github.RepositoryContentGetOptions{
		Ref: "",
	}

	nowAsUTC := getTimeProtobufInUTC()

	kurtosisYamlFilePath := fmt.Sprintf("%s%s", packageRootPath, kurtosisYamlFilename)

	// get contents of kurtosis yaml file from GitHub
	kurtosisYamlFileContentResult, _, resp, err := client.Repositories.GetContents(ctx, repositoryOwner, repositoryName, kurtosisYamlFilePath, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == http.StatusNotFound {
		logrus.Debugf("No '%s' file in repo '%s'", kurtosisYamlFilePath, repositoryFullName)
		return nil, false, nil
	} else if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s' - file '%s'", repositoryFullName, kurtosisYamlFilePath)
	}

	kurtosisPackageName, kurtosisPackageDescription, commitSHA, err := ParseKurtosisYaml(kurtosisYamlFileContentResult)
	if err != nil {
		logrus.Warnf("An error occurred parsing '%s' file in repository '%s'"+
			"Error was:\n%v", kurtosisYamlFilePath, repositoryFullName, err.Error())
		return nil, false, stacktrace.Propagate(err, "An error occurred parsing the '%v' file.", defaultKurtosisYamlFilename)
	}

	// we check that the name set in kurtosis.yml matches the location on GitHub. If not, we exclude it from the
	// indexer b/c users won't be able to run it with `kurtosis run <PACKAGE_NAME>`
	expectedPackageName := normalizeName(
		fmt.Sprintf(
			"%s/%s/%s/%s",
			githubUrl,
			repositoryOwner,
			repositoryName,
			packageRootPath,
		),
	)
	if normalizeName(kurtosisPackageName) != expectedPackageName {
		nameError := stacktrace.NewError(
			"The package '%s' is invalid because the name set in its '%s' doesn't match its Github URL (name set: '%s' - expected: '%s').",
			repositoryFullName,
			defaultKurtosisYamlFilename,
			kurtosisPackageName,
			expectedPackageName,
		)
		logrus.Warn("package name must match the repo name", nameError)
		return nil, false, nameError
	}

	// get contents of main.star file from GitHub
	kurtosisMainDotStarFilePath := fmt.Sprintf("%s%s", packageRootPath, starlarkMainDotStarFileName)
	starlarkMainDotStartContentResult, _, resp, err := client.Repositories.GetContents(ctx, repositoryOwner, repositoryName, kurtosisMainDotStarFilePath, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == http.StatusNotFound {
		logrus.Debugf("No '%s' file in repo '%s'", kurtosisMainDotStarFilePath, repositoryFullName)
		return nil, false, nil
	} else if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s' - file '%s'", repositoryFullName, kurtosisMainDotStarFilePath)
	}
	mainDotStarParsedContent, err := ParseStarlarkMainDotStar(starlarkMainDotStartContentResult)
	if err != nil {
		logrus.Warnf("An error occurred parsing '%s' file in repository '%s'. This Kurtosis package will not be indexed. "+
			"Error was:\n%v", kurtosisMainDotStarFilePath, repositoryFullName, err.Error())
		return nil, false, stacktrace.Propagate(err, "An error occurred parsing the '%v' file.", starlarkMainDotStarFileName)
	}

	if err != addOrUpdatePackageRepositoryMetadataWithStarsAndLastCommitDate(ctx, client, repositoryOwner, repositoryName, packageRepositoryMetadata) {
		logrus.Warnf("an error occurred while adding or updating the repo starts and last commit date for repository '%s'. Error was:\n%v", repositoryName, err.Error())
	}

	// check if the Kurtosis package icon exist in the package's root and get the URL
	var packageIconURL string
	packageIconURL, err = getPackageIconURLIfExists(packageRepositoryMetadata)
	if err != nil {
		logrus.Warnf("an error occurred while getting and checking if the package icon exists in repository '%s'. Error was:\n%v", repositoryName, err.Error())
	}

	// add the run count metrics
	runCount, found := packagesRunCount[kurtosisPackageName]
	if !found {
		logrus.Infof("no package run metrics found for '%s'", kurtosisPackageName)
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
		packageIconURL,
		runCount,
		mainDotStarParsedContent.Arguments...,
	), true, nil
}

func addOrUpdatePackageRepositoryMetadataWithStarsAndLastCommitDate(
	ctx context.Context,
	client *github.Client,
	repositoryOwner string,
	repositoryName string,
	packageRepositoryMetadata *PackageRepositoryMetadata,
) error {

	repository, _, err := client.Repositories.Get(ctx, repositoryOwner, repositoryName)
	if err != nil {
		return stacktrace.Propagate(err, "an error occurred getting repository '%s' owned by '%s'", repositoryName, repositoryOwner)
	}
	packageRepositoryMetadata.DefaultBranch = repository.GetDefaultBranch()

	// add or update the stars
	// is necessary to check for nil because if it's not returned in the response the number will be 0 when we call GetStargazersCount()
	// and this could overwrite a previous value
	if repository.StargazersCount == nil {
		return stacktrace.NewError("no stars received when calling GitHub for repository '%s' it should return at least 0 stars", repositoryName)
	}
	packageRepositoryMetadata.Stars = uint64(repository.GetStargazersCount())

	// add or update the latest commit date
	// nolint:exhaustruct
	commitListOptions := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: onlyOneCommit, // it will return the latest one from the main branch
		},
	}
	repositoryCommits, _, err := client.Repositories.ListCommits(ctx, repositoryOwner, repositoryName, commitListOptions)
	if err != nil {
		return stacktrace.Propagate(err, "an error occurred getting the commit list for repository '%s'", repositoryName)
	}
	if len(repositoryCommits) == 0 {
		return stacktrace.NewError("zero commits were received when calling GitHub for getting the last commit for repository '%s'", repositoryName)
	}

	latestCommit := repositoryCommits[0]
	if latestCommit == nil {
		return stacktrace.NewError("an error occurred while trying to get the last commit form the received repository commits response '%+v'", repositoryCommits)
	}
	latestCommitGitHubTimestamp := latestCommit.GetCommit().GetCommitter().GetDate()
	packageRepositoryMetadata.LastCommitTime = *latestCommitGitHubTimestamp.GetTime()

	return nil
}

func getRepositoryOwner(repository *github.Repository) (string, error) {
	// It's not clear what to use in the `GetContents` function called below, as the `owner` field. We would prefer to
	// just pass in githubRepository.GetFullName() and that's it, but it doesn't work
	// So, to parse the owner name, we try both the `Login` and `Name` field, taking the first one that is not empty.
	// For repo owned by `kurtosis-tech` for example, it's the Login field that will be used, as the Name one is empty
	// If this is too fragile, worst case we can regexp parse `githubRepository.GetFullName()`
	if repository.GetOwner().GetLogin() != "" {
		return repository.GetOwner().GetLogin(), nil
	} else if repository.GetOwner().GetName() != "" {
		return repository.GetOwner().GetName(), nil
	}

	return "", stacktrace.NewError("impossible to get owner from Github repository '%+v'", repository)

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

func getPackageIconURLIfExists(packageRepositoryMetadata *PackageRepositoryMetadata) (string, error) {
	repositoryDownloadRootURL, err := packageRepositoryMetadata.GetDownloadRootURL()
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred getting the package repository download root URL")
	}
	packageIconURL, err := url.JoinPath(repositoryDownloadRootURL, kurtosisPackageIconImgName)
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred generating the package icon URL for repository '%s'", packageRepositoryMetadata.Name)
	}
	response, err := http.Head(packageIconURL)
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred while sending a HEAD request to '%s'", packageIconURL)
	}
	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusPartialContent {
		return packageIconURL, nil
	}
	return "", nil
}
