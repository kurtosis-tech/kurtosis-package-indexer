package store

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/stacktrace"
	"sort"
	"strings"
	"sync"
	"time"
)

type InMemoryStore struct {
	lastMainCrawlTime                time.Time
	lastSecondaryCrawlTime           time.Time
	lastMetricsReporterQueryDatetime time.Time

	// using sync.Map because we have several sub-routines (the scheduled jobs) writing on it and this map supports concurrency
	packages         sync.Map
	packagesRunCount types.PackagesRunCount
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		lastMainCrawlTime:                time.Time{},
		lastSecondaryCrawlTime:           time.Time{},
		lastMetricsReporterQueryDatetime: time.Time{},
		packages:                         sync.Map{},
		packagesRunCount:                 types.PackagesRunCount{},
	}
}

func (store *InMemoryStore) Close() error {
	return nil
}

func (store *InMemoryStore) GetKurtosisPackages(_ context.Context) ([]*generated.KurtosisPackage, error) {
	var packages []*generated.KurtosisPackage
	store.packages.Range(func(k, kurtosisPackage interface{}) bool {
		packages = append(packages, kurtosisPackage.(*generated.KurtosisPackage))
		return true
	})
	sort.SliceStable(packages, func(i, j int) bool {
		return strings.ToLower(packages[i].GetUrl()) < strings.ToLower(packages[j].GetUrl())
	})
	return packages, nil
}

func (store *InMemoryStore) GetKurtosisPackage(_ context.Context, kurtosisPackageLocator string) (*generated.KurtosisPackage, error) {
	kurtosisPackage, found := store.packages.Load(kurtosisPackageLocator)
	if !found {
		return nil, stacktrace.NewError("Expected to find the kurtosis package '%s' in the storage but it was not found", kurtosisPackageLocator)
	}
	return kurtosisPackage.(*generated.KurtosisPackage), nil
}

func (store *InMemoryStore) UpsertPackage(_ context.Context, kurtosisPackageLocator string, kurtosisPackage *generated.KurtosisPackage) error {
	store.packages.Store(kurtosisPackageLocator, kurtosisPackage)
	return nil
}

func (store *InMemoryStore) DeletePackage(_ context.Context, packageLocator string) error {
	store.packages.Delete(packageLocator)
	return nil
}

func (store *InMemoryStore) UpdateLastMainCrawlDatetime(_ context.Context, lastCrawlTime time.Time) error {
	store.lastMainCrawlTime = lastCrawlTime
	return nil
}

func (store *InMemoryStore) GetLastMainCrawlDatetime(_ context.Context) (time.Time, error) {
	return store.lastMainCrawlTime, nil
}

func (store *InMemoryStore) UpdateLastSecondaryCrawlDatetime(_ context.Context, lastCrawlTime time.Time) error {
	store.lastSecondaryCrawlTime = lastCrawlTime
	return nil
}

func (store *InMemoryStore) GetLastSecondaryCrawlDatetime(_ context.Context) (time.Time, error) {
	return store.lastSecondaryCrawlTime, nil
}

func (store *InMemoryStore) UpdateLastMetricsQueryDatetime(_ context.Context, lastMetricsQueryTime time.Time) error {
	store.lastMetricsReporterQueryDatetime = lastMetricsQueryTime
	return nil
}

func (store *InMemoryStore) GetLastMetricsQueryDatetime(_ context.Context) (time.Time, error) {
	return store.lastMetricsReporterQueryDatetime, nil
}

func (store *InMemoryStore) GetPackagesRunCount(_ context.Context) (types.PackagesRunCount, error) {
	return store.packagesRunCount, nil
}

func (store *InMemoryStore) UpdatePackagesRunCount(_ context.Context, newPackagesRunCount types.PackagesRunCount) error {
	store.packagesRunCount = newPackagesRunCount
	return nil
}
