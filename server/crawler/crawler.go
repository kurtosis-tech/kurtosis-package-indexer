package crawler

import (
	"context"
	"fmt"
	"github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/api_constructors"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/catalog"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/consts"
	github_indexer_lib "github.com/kurtosis-tech/kurtosis-package-indexer/server/github"
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
	mainCrawlFrequency      = 5 * time.Minute
	secondaryCrawlFrequency = 30 * time.Minute
	crawlIntervalBuffer     = 15 * time.Second

	defaultSecondaryCrawlInitialDelay = 30 * time.Minute

	githubUrl = "github.com"

	successfulParsingText = "Parsed package content successfully"

	noStartsSet = uint64(0)

	noDefaultBranchSet = ""

	onlyOneCommit = 1
)

var supportedDockerComposeYmlFilenames = []string{
	"compose.yml",
	"compose.yaml",
	"docker-compose.yml",
	"docker-compose.yaml",
	"docker_compose.yml",
	"docker_compose.yaml",
}

var (
	zeroValueTime = time.Time{}
)

type GithubCrawler struct {
	store           store.KurtosisIndexerStore
	ctx             context.Context
	mainTicker      *ticker.Ticker
	secondaryTicker *ticker.Ticker
}

func NewGitHubCrawler(ctx context.Context, store store.KurtosisIndexerStore) (*GithubCrawler, error) {

	newCrawler := &GithubCrawler{
		store:           store,
		ctx:             ctx,
		mainTicker:      nil,
		secondaryTicker: nil,
	}

	return newCrawler, nil
}

func (crawler *GithubCrawler) Schedule(forceRunNow bool) error {
	if crawler.mainTicker != nil {
		logrus.Infof("Main crawler already scheduled - stopping it first")
		crawler.mainTicker.Stop()
	}

	if err := crawler.scheduleMainCrawler(forceRunNow); err != nil {
		return stacktrace.Propagate(err, "an error occurred scheduling the main crawler")
	}

	if crawler.secondaryTicker != nil {
		logrus.Infof("Secondary crawler already scheduled - stopping it first")
		crawler.secondaryTicker.Stop()
	}

	if err := crawler.scheduleSecondaryCrawler(forceRunNow); err != nil {
		return stacktrace.Propagate(err, "an error occurred scheduling the secondary crawler")
	}

	return nil
}

// scheduleMainCrawler schedule the main crawler which is in charge of adding or updating new packages in the store
// the main difference with the secondary crawler is that it's not update the repository stars and the last commit time
// for the packages that are already stored
func (crawler *GithubCrawler) scheduleMainCrawler(forceRunNow bool) error {
	lastCrawlDatetime, err := crawler.store.GetLastMainCrawlDatetime(crawler.ctx)
	if err != nil {
		return stacktrace.Propagate(err, "An unexpected error occurred retrieving last main crawl datetime from the store")
	}

	var initialDelay time.Duration
	if !forceRunNow && lastCrawlDatetime.Add(mainCrawlFrequency).After(time.Now()) {
		initialDelay = time.Until(lastCrawlDatetime.Add(mainCrawlFrequency))
	}
	logrus.Infof("Main crawler starting with an initial delay of '%v' and a period of '%v'", initialDelay, mainCrawlFrequency)

	crawler.mainTicker = ticker.NewTicker(initialDelay, mainCrawlFrequency)
	go func() {
		for {
			select {
			case <-crawler.ctx.Done():
				logrus.Info("Main crawler has been closed. Returning")
				crawler.mainTicker.Stop()
				return
			case tickerTime := <-crawler.mainTicker.C:
				crawler.doMainCrawlNoFailure(crawler.ctx, tickerTime)
			}
		}
	}()
	return nil
}

// scheduleSecondaryCrawler schedule the secondary crawler which is in charge of updating the package's repository stars
// and last commit time for each stored package
func (crawler *GithubCrawler) scheduleSecondaryCrawler(forceRunNow bool) error {
	lastCrawlDatetime, err := crawler.store.GetLastSecondaryCrawlDatetime(crawler.ctx)
	if err != nil {
		return stacktrace.Propagate(err, "An unexpected error occurred retrieving last secondary crawl datetime from the store")
	}

	// It has an initial delay the first time it runs because the main crawler is going to get the stars and the last commit time for
	// all the packages the first time it get them
	initialDelay := defaultSecondaryCrawlInitialDelay
	if !forceRunNow && lastCrawlDatetime.Add(secondaryCrawlFrequency).After(time.Now()) {
		initialDelay = time.Until(lastCrawlDatetime.Add(secondaryCrawlFrequency))
	}
	logrus.Infof("Secondary crawler starting with an initial delay of '%v' and a period of '%v'", initialDelay, secondaryCrawlFrequency)

	crawler.secondaryTicker = ticker.NewTicker(initialDelay, secondaryCrawlFrequency)
	go func() {
		for {
			select {
			case <-crawler.ctx.Done():
				logrus.Info("Secondary crawler has been closed. Returning")
				crawler.secondaryTicker.Stop()
				return
			case tickerTime := <-crawler.secondaryTicker.C:
				crawler.doSecondaryCrawlNoFailure(crawler.ctx, tickerTime)
			}
		}
	}()
	return nil
}

// ReadPackage returns a KurtosisPackage requested by its repository owner, name and root path
// the KurtosisPackage returned should we filled with fresh data expect for some fields
// we decided to cache in the storage, like the repository stars and the last commit time
func (crawler *GithubCrawler) ReadPackage(
	ctx context.Context,
	apiRepositoryMetadata *generated.PackageRepository,
) (*generated.KurtosisPackage, error) {

	packageMetadata := crawler.getPackageMetadataFromStorageOrDefault(ctx, apiRepositoryMetadata.GetOwner(), apiRepositoryMetadata.GetName(), apiRepositoryMetadata.GetRootPath())

	githubClient, err := createGithubClient(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred while creating the github client")
	}

	packagesRunCount := crawler.getPackagesRunCountNotFailure(ctx)
	kurtosisPackageContent, packageFound, err := extractKurtosisPackageContent(ctx, githubClient, packageMetadata, packagesRunCount)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred extracting content for Kurtosis package repository '%s'", packageMetadata.GetLocator())
	}
	if !packageFound {
		logrus.Warn("No Kurtosis package found.") // don't want to log provided package repository locator bc it's a security risk (e.g. malicious data)
		return nil, stacktrace.NewError("No Kurtosis package found. Ensure that a package exists at '%v' with valid '%v' and '%v' files or one of the following support docker compose files '%v'.", packageMetadata.GetLocator(), consts.DefaultKurtosisYamlFilename, consts.StarlarkMainDotStarFileName, supportedDockerComposeYmlFilenames)
	}

	kurtosisPackageApi := convertRepoContentToApi(kurtosisPackageContent)
	return kurtosisPackageApi, nil
}

func (crawler *GithubCrawler) getPackageMetadataFromStorageOrDefault(ctx context.Context, repositoryOwner string, repositoryName string, repositoryPackageRootPath string) *PackageRepositoryMetadata {

	// setting default values for the cached data
	repositoryStars := noStartsSet
	repositoryLastCommitTime := zeroValueTime
	repositoryDefaultBranch := noDefaultBranchSet

	// creating this first package metadata just to call the GetLocator method
	incompletePackageMetadata := NewPackageRepositoryMetadata(
		repositoryOwner,
		repositoryName,
		repositoryPackageRootPath,
		consts.DefaultKurtosisYamlFilename,
		repositoryStars,          // it will be filled below when we get this from the storage to avoid shooting an extra GitHub request
		repositoryLastCommitTime, // it will be filled below when we get this from the storage to avoid shooting an extra GitHub request
		repositoryDefaultBranch,  // it will be filled below when we get this from the storage to avoid shooting an extra GitHub request
	)

	packageLocator := incompletePackageMetadata.GetLocator()

	// Doing the best effort to get the cached package info from the storage but, it could fail because the package
	// was not already saved in the storage or, so, default values will be returned instead
	kurtosisPackage, err := crawler.store.GetKurtosisPackage(ctx, packageLocator)
	if err != nil {
		logrus.Errorf("an error occurred getting Kurtosis package '%s' from the storage, this could be because the packages was not already saved in the storage. Error was:\n%s", packageLocator, err.Error())
	} else {
		repositoryStars = kurtosisPackage.GetStars()
		repositoryLastCommitTime = kurtosisPackage.GetRepositoryMetadata().GetLastCommitTime().AsTime()
		repositoryDefaultBranch = kurtosisPackage.GetRepositoryMetadata().GetDefaultBranch()
	}

	newPackageRepositoryMetadata := NewPackageRepositoryMetadata(
		repositoryOwner,
		repositoryName,
		repositoryPackageRootPath,
		consts.DefaultKurtosisYamlFilename,
		repositoryStars,
		repositoryLastCommitTime,
		repositoryDefaultBranch,
	)

	return newPackageRepositoryMetadata
}

func (crawler *GithubCrawler) doMainCrawlNoFailure(ctx context.Context, tickerTime time.Time) {
	crawlSuccessful := false
	lastCrawlDatetime, err := crawler.store.GetLastMainCrawlDatetime(ctx)
	if err != nil {
		logrus.Errorf("Could not retrieve last main crawl datetime from the store. This is not critical, crawling"+
			"will just continue even though the last main crawl might be recent. Error was:\n%s", err.Error())
	}

	// Add a small buffer to avoid false positive when checking if the new crawl is sooner than the crawl frequency interval
	lastCrawlDatetimeWithBuffer := lastCrawlDatetime.Add(-crawlIntervalBuffer)
	if time.Since(lastCrawlDatetimeWithBuffer) < mainCrawlFrequency {
		logrus.Infof("Last main crawling happened as '%v' ('%v' ago), which is more recent than the crawling frequency "+
			"set to '%v', so crawling will be skipped. If the indexer is running with more than one node, it might be "+
			"that another node did the crawling in between, and this is totally expected.",
			lastCrawlDatetimeWithBuffer, time.Since(lastCrawlDatetimeWithBuffer), mainCrawlFrequency)
		return
	}

	// we persist the crawl datetime before doing the crawl so that potential other nodes don't crawl at the same time
	currentCrawlDatetime := time.Now()
	if err = crawler.store.UpdateLastMainCrawlDatetime(crawler.ctx, currentCrawlDatetime); err != nil {
		logrus.Errorf("An error occurred persisting main crawl time to database. This is not critical, but in case of "+
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
			if err = crawler.store.UpdateLastMainCrawlDatetime(crawler.ctx, lastCrawlDatetime); err != nil {
				logrus.Errorf("An error occurred reverting the last main crawl datetime to '%s'. Its value"+
					"will remain '%s' and no main crawl will happen before '%s'. Error was:\n%v",
					lastCrawlDatetime, currentCrawlDatetime, currentCrawlDatetime.Add(mainCrawlFrequency), err.Error())
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
		logrus.Errorf("An error occurred crawling Github for Kurtosis packages. The last main crawl datetime"+
			"will be reverted to its previous value '%v'. This node will try crawling again in '%v'. "+
			"Error was:\n%s", lastCrawlDatetime, mainCrawlFrequency, err.Error())
		crawlSuccessful = false
	} else {
		crawlSuccessful = true
	}
	logrus.Infof("Crawling finished in %v. Success: '%v'. Repository updated: %d - added: %d - removed: %d",
		time.Since(tickerTime), crawlSuccessful, packageUpdated, packageAdded, packageRemoved)
}

func (crawler *GithubCrawler) doSecondaryCrawlNoFailure(ctx context.Context, tickerTime time.Time) {
	crawlSuccessful := false
	lastCrawlDatetime, err := crawler.store.GetLastSecondaryCrawlDatetime(ctx)
	if err != nil {
		logrus.Errorf("Could not retrieve last secondary crawl datetime from the store. This is not critical, crawling"+
			"will just continue even though the last secondary crawl might be recent. Error was:\n%s", err.Error())
	}

	// Add a small buffer to avoid false positive when checking if the new crawl is sooner than the crawl frequency interval
	lastCrawlDatetimeWithBuffer := lastCrawlDatetime.Add(-crawlIntervalBuffer)
	if time.Since(lastCrawlDatetimeWithBuffer) < secondaryCrawlFrequency {
		logrus.Infof("Last secondary crawling happened as '%v' ('%v' ago), which is more recent than the crawling frequency "+
			"set to '%v', so crawling will be skipped. If the indexer is running with more than one node, it might be "+
			"that another node did the crawling in between, and this is totally expected.",
			lastCrawlDatetimeWithBuffer, time.Since(lastCrawlDatetimeWithBuffer), secondaryCrawlFrequency)
		return
	}

	// we persist the crawl datetime before doing the crawl so that potential other nodes don't crawl at the same time
	currentCrawlDatetime := time.Now()
	if err = crawler.store.UpdateLastSecondaryCrawlDatetime(crawler.ctx, currentCrawlDatetime); err != nil {
		logrus.Errorf("An error occurred persisting secondary crawl time to database. This is not critical, but in case of "+
			"a service restart (or in a multiple nodes environment), crawling might happen more frequently than "+
			"expected. Error was was:\n%v", err.Error())
	} else {
		defer func() {
			if crawlSuccessful {
				return
			}
			logrus.Debugf("Reverting the last main crawl datetime to '%s'. Current value is '%s'",
				lastCrawlDatetime, currentCrawlDatetime)
			// revert the crawl datetime to its previous value
			if err = crawler.store.UpdateLastSecondaryCrawlDatetime(crawler.ctx, lastCrawlDatetime); err != nil {
				logrus.Errorf("An error occurred reverting the last secondary crawl datetime to '%s'. Its value"+
					"will remain '%s' and no secondary crawl will happen before '%s'. Error was:\n%v",
					lastCrawlDatetime, currentCrawlDatetime, currentCrawlDatetime.Add(secondaryCrawlFrequency), err.Error())
			}
		}()
	}

	githubClient, err := createGithubClient(ctx)
	if err != nil {
		logrus.Error("An error occurred while creating the github client. Aborting", err)
		return
	}

	logrus.Info("Crawling Github for updating Kurtosis packages stars and last commit time")
	packageUpdated, packagesFailed, err := crawler.updateAllPackagesStarsDefaultBranchAndLastCommitDate(ctx, githubClient)
	if err != nil {
		logrus.Errorf("An error occurred crawling GitHub for Kurtosis packages stars and last commit time. "+
			"The last secondary crawl datetime will be reverted to its previous value '%v'. This node will try crawling again in '%v'. "+
			"Error was:\n%s", lastCrawlDatetime, secondaryCrawlFrequency, err.Error())
		crawlSuccessful = false
	} else {
		crawlSuccessful = true
	}
	logrus.Infof("Crawling for stars and last commit time finished in %v. Success: '%v'. Repository updated: %d - failed: %d ",
		time.Since(tickerTime), crawlSuccessful, packageUpdated, packagesFailed)
}

func (crawler *GithubCrawler) updateAllPackagesStarsDefaultBranchAndLastCommitDate(ctx context.Context, client *github.Client) (uint32, uint32, error) {
	allStoredPackages, err := crawler.store.GetKurtosisPackages(ctx)
	if err != nil {
		return 0, 0, stacktrace.Propagate(err, "An error occurred retrieving currently stored packages")
	}

	packagesUpdated := uint32(0)
	packagesFailed := uint32(len(allStoredPackages))
	for _, kurtosisPackage := range allStoredPackages {
		repositoryOwner := kurtosisPackage.RepositoryMetadata.GetOwner()
		repositoryName := kurtosisPackage.RepositoryMetadata.GetName()
		repositoryPackageRootPath := kurtosisPackage.RepositoryMetadata.GetRootPath()

		defaultBranch, stars, lastCommitTime, err := getRepositoryStarsDefaultBranchAndLastCommitDate(ctx, client, repositoryOwner, repositoryName)
		if err != nil {
			return 0, 0, stacktrace.Propagate(err, "an error occurred while getting the repo stars, default branch and last commit time for repository '%s'", repositoryName)
		}

		kurtosisPackage.Stars = stars
		kurtosisPackage.RepositoryMetadata.DefaultBranch = defaultBranch
		kurtosisPackage.RepositoryMetadata.LastCommitTime = timestamppb.New(*lastCommitTime)

		packageRepositoryMetadata := NewPackageRepositoryMetadata(
			repositoryOwner,
			repositoryName,
			repositoryPackageRootPath,
			consts.DefaultKurtosisYamlFilename,
			stars,
			*lastCommitTime,
			defaultBranch,
		)

		kurtosisPackageLocator := packageRepositoryMetadata.GetLocator()

		if err := crawler.store.UpsertPackage(ctx, kurtosisPackageLocator, kurtosisPackage); err != nil {
			return 0, 0, stacktrace.Propagate(err, "an error occurred updating Kurtosis package '%s'", kurtosisPackageLocator)
		}
		logrus.Debugf("Kurtosis package stars and last commit time updated for package '%s'", kurtosisPackage.Name)
		packagesFailed--
		packagesUpdated++
	}

	return packagesUpdated, packagesFailed, nil
}

// TODO review this because it could be optimized, right now the ReadPackages method is calling it
// TODO and creating a new client on each request
func createGithubClient(ctx context.Context) (*github.Client, error) {
	githubClient, err := github_indexer_lib.CreateGithubClient(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Unable to build authenticated Github client. This is required so that "+
			"the indexer can search Github to retrieve Kurtosis package content. Skipping indexing for now, will "+
			"retry in %v", mainCrawlFrequency)
	}
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
		storedApiRepositoryMetadata := storedPackage.GetRepositoryMetadata()
		kurtosisPackageMetadata := NewPackageRepositoryMetadata(
			storedApiRepositoryMetadata.GetOwner(),
			storedApiRepositoryMetadata.GetName(),
			storedApiRepositoryMetadata.GetRootPath(),
			consts.DefaultKurtosisYamlFilename,
			storedPackage.GetStars(),                                 // this will be updated by the secondary crawler
			storedApiRepositoryMetadata.GetLastCommitTime().AsTime(), // this will be updated by the secondary crawler
			storedApiRepositoryMetadata.GetDefaultBranch(),           // this will be updated by the secondary crawler
		)
		packageRepositoryLocator := kurtosisPackageMetadata.GetLocator()
		logrus.Debugf("Trying to update content of package '%s'", packageRepositoryLocator) // TODO: remove log line

		packagesRunCount := crawler.getPackagesRunCountNotFailure(ctx)
		kurtosisPackageContent, packageFound, err := extractKurtosisPackageContent(ctx, githubClient, kurtosisPackageMetadata, packagesRunCount)
		if err != nil {
			logrus.Warnf("An error occurred extracting content for Kurtosis package repository '%s'. Error was:\n%s", packageRepositoryLocator, err.Error())
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
	allKurtosisPackagesRepositoryMetadata, err := readCatalogAndGetPackagesRepositoryMetadata(ctx, githubClient)
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

// readCatalogAndGetPackagesRepositoryMetadata imports and reads the packages catalog for getting the latest curated package list
// from the kurtosis-package-catalog repository and, it sends two requests to GitHub in order to create the PackageRepositoryMetadata object
// for each package, for this reason this function should be called only when adding new packages for avoid GitHub rate limits
func readCatalogAndGetPackagesRepositoryMetadata(ctx context.Context, client *github.Client) ([]*PackageRepositoryMetadata, error) {

	packagesCatalog, err := catalog.GetPackagesCatalog()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the packages catalog")
	}

	allPackageRepositoryMetadatas := make([]*PackageRepositoryMetadata, len(packagesCatalog))
	for packageIndex, packageData := range packagesCatalog {
		repositoryOwner := packageData.GetRepositoryOwner()
		repositoryName := packageData.GetRepositoryName()
		repositoryPackageRootPath := packageData.GetRepositoryPackageRootPath()

		defaultBranch, stars, lastCommitTime, getErr := getRepositoryStarsDefaultBranchAndLastCommitDate(ctx, client, repositoryOwner, repositoryName)
		if getErr != nil {
			return nil, stacktrace.Propagate(getErr, "an error occurred getting the repository stars, default branch and last commit date for owner '%s' and name '%s'", repositoryOwner, repositoryName)
		}

		newPackageRepositoryMetadata := NewPackageRepositoryMetadata(repositoryOwner, repositoryName, repositoryPackageRootPath, consts.DefaultKurtosisYamlFilename, stars, *lastCommitTime, defaultBranch)

		allPackageRepositoryMetadatas[packageIndex] = newPackageRepositoryMetadata
	}

	return allPackageRepositoryMetadatas, nil
}

// If the package doesn't exist, returns an empty package content, false, and empty error
// If package exists, but error occurred while extracting package content, returns empty package content, false, with an error
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

	repositoryFullName := packageRepositoryMetadata.GetLocator()

	repoGetContentOpts := &github.RepositoryContentGetOptions{
		Ref: "",
	}

	nowAsUTC := getTimeProtobufInUTC()

	kurtosisYamlFilePath := fmt.Sprintf("%s%s", packageRootPath, kurtosisYamlFilename)

	// get contents of kurtosis yaml file from GitHub
	kurtosisYamlFileContentResult, _, resp, err := client.Repositories.GetContents(ctx, repositoryOwner, repositoryName, kurtosisYamlFilePath, repoGetContentOpts)
	if err != nil && resp != nil && resp.StatusCode == http.StatusNotFound {
		logrus.Debugf("No '%s' file in repo '%s'", kurtosisYamlFilePath, repositoryFullName)
		logrus.Debugf("Attempting to extract a docker compose based package...")
		composePackageContent, composePackageFound, err := extractDockerComposePackageContent(ctx, client, packageRepositoryMetadata)
		if err != nil {
			return nil, false, err
		}
		if !composePackageFound {
			logrus.Debugf("No docker compose based package found in repo %v.", repositoryFullName)
			return nil, false, nil
		}
		return composePackageContent, true, nil
	} else if err != nil {
		return nil, false, stacktrace.Propagate(err, "An error occurred reading content of Kurtosis Package '%s' - file '%s'", repositoryFullName, kurtosisYamlFilePath)
	}

	kurtosisPackageName, kurtosisPackageDescription, commitSHA, err := ParseKurtosisYaml(kurtosisYamlFileContentResult)
	if err != nil {
		logrus.Warnf("An error occurred parsing '%s' file in repository '%s'"+
			"Error was:\n%v", kurtosisYamlFilePath, repositoryFullName, err.Error())
		return nil, false, stacktrace.Propagate(err, "An error occurred parsing the '%v' file.", consts.DefaultKurtosisYamlFilename)
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
			consts.DefaultKurtosisYamlFilename,
			kurtosisPackageName,
			expectedPackageName,
		)
		logrus.Warn("package name must match the repo name", nameError)
		return nil, false, nameError
	}

	// get contents of main.star file from GitHub
	kurtosisMainDotStarFilePath := fmt.Sprintf("%s%s", packageRootPath, consts.StarlarkMainDotStarFileName)
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
		return nil, false, stacktrace.Propagate(err, "An error occurred parsing the '%v' file.", consts.StarlarkMainDotStarFileName)
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

// If the docker compose yaml doesn't exist, returns an empty package content, false, and empty error
// If docker compose yaml exists, but error occurred while extracting package content, returns empty package content, false, with an error
func extractDockerComposePackageContent(
	ctx context.Context,
	client *github.Client,
	packageRepositoryMetadata *PackageRepositoryMetadata) (*KurtosisPackageContent, bool, error) {

	repoGetContentOpts := &github.RepositoryContentGetOptions{
		Ref: "",
	}

	nowAsUTC := getTimeProtobufInUTC()

	var dockerComposeYamlFileContentResult *github.RepositoryContent
	var dockerComposeYamlFilePath string
	for _, dockerComposeYamlVariation := range supportedDockerComposeYmlFilenames {
		dockerComposeYamlFilePath = fmt.Sprintf("%s/%s", packageRepositoryMetadata.RootPath, dockerComposeYamlVariation)

		yamlFileContentResult, _, resp, err := client.Repositories.GetContents(ctx, packageRepositoryMetadata.Owner, packageRepositoryMetadata.Name, dockerComposeYamlFilePath, repoGetContentOpts)
		if err == nil && resp != nil && resp.StatusCode != 404 {
			dockerComposeYamlFileContentResult = yamlFileContentResult
			break
		}
	}
	if dockerComposeYamlFileContentResult == nil {
		return nil, false, nil
	}

	// TODO: Parse dockerComposeYamlFileContentResult for metadata about the compose file (similar to main dot star)

	// no notion of main dot star in docker compose so leave main.star specific fields blank for now
	return NewKurtosisPackageContent(
		packageRepositoryMetadata,
		normalizeName(fmt.Sprintf("%s/%s/%s/%s", githubUrl, packageRepositoryMetadata.Owner, packageRepositoryMetadata.Name, packageRepositoryMetadata.RootPath)),
		"",
		"",
		"", // p
		successfulParsingText,
		nowAsUTC,
		"",
		"",
		0, // getting run count for docker compose packages
		nil,
	), true, nil
}

func getRepositoryStarsDefaultBranchAndLastCommitDate(
	ctx context.Context,
	client *github.Client,
	repositoryOwner string,
	repositoryName string,
) (string, uint64, *time.Time, error) {

	repository, _, err := client.Repositories.Get(ctx, repositoryOwner, repositoryName)
	if err != nil {
		return "", 0, nil, stacktrace.Propagate(err, "an error occurred getting repository '%s' owned by '%s'", repositoryName, repositoryOwner)
	}

	// get default branch
	defaultBranch := noDefaultBranchSet
	if repository.GetDefaultBranch() != noDefaultBranchSet {
		defaultBranch = repository.GetDefaultBranch()
	}

	// get the stars
	// is necessary to check for nil because if it's not returned in the response the number will be 0 when we call GetStargazersCount()
	// and this could overwrite a previous value
	if repository.StargazersCount == nil {
		return "", 0, nil, stacktrace.NewError("no stars received when calling GitHub for repository '%s' it should return at least 0 stars", repositoryName)
	}
	stars := uint64(repository.GetStargazersCount())

	// get the latest commit date
	lastCommitTime, err := getLastCommitTimeFromGitHub(ctx, client, repositoryOwner, repositoryName)
	if err != nil {
		return "", 0, nil, stacktrace.Propagate(err, "an error occurred getting the last commit time from GitHub")
	}

	return defaultBranch, stars, lastCommitTime, nil
}

func getLastCommitTimeFromGitHub(
	ctx context.Context,
	client *github.Client,
	repositoryOwner string,
	repositoryName string,
) (*time.Time, error) {

	// nolint:exhaustruct
	commitListOptions := &github.CommitsListOptions{
		ListOptions: github.ListOptions{
			PerPage: onlyOneCommit, // it will return the latest one from the main branch
		},
	}
	repositoryCommits, _, err := client.Repositories.ListCommits(ctx, repositoryOwner, repositoryName, commitListOptions)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the commit list for repository '%s'", repositoryName)
	}
	if len(repositoryCommits) == 0 {
		return nil, stacktrace.NewError("zero commits were received when calling GitHub for getting the last commit for repository '%s'", repositoryName)
	}

	latestCommit := repositoryCommits[0]
	if latestCommit == nil {
		return nil, stacktrace.NewError("an error occurred while trying to get the last commit form the received repository commits response '%+v'", repositoryCommits)
	}
	latestCommitGitHubTimestamp := latestCommit.GetCommit().GetCommitter().GetDate()
	return latestCommitGitHubTimestamp.GetTime(), nil
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
	packageIconURL, err := url.JoinPath(repositoryDownloadRootURL, consts.KurtosisPackageIconImgName)
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
