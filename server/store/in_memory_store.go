package store

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"sort"
	"strings"
	"time"
)

type InMemoryStore struct {
	lastMainCrawlTime                time.Time
	lastSecondaryCrawlTime           time.Time
	lastMetricsReporterQueryDatetime time.Time

	packages         map[string]*generated.KurtosisPackage
	packagesRunCount types.PackagesRunCount
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		lastMainCrawlTime:                time.Time{},
		lastMetricsReporterQueryDatetime: time.Time{},
		packages:                         map[string]*generated.KurtosisPackage{},
		packagesRunCount:                 types.PackagesRunCount{},
	}
}

func (store *InMemoryStore) Close() error {
	return nil
}

func (store *InMemoryStore) GetKurtosisPackages(_ context.Context) ([]*generated.KurtosisPackage, error) {
	var packages []*generated.KurtosisPackage
	for _, kurtosisPackage := range store.packages {
		packages = append(packages, kurtosisPackage)
	}
	sort.SliceStable(packages, func(i, j int) bool {
		return strings.ToLower(packages[i].GetUrl()) < strings.ToLower(packages[j].GetUrl())
	})
	return packages, nil
}

func (store *InMemoryStore) UpsertPackage(_ context.Context, kurtosisPackageLocator string, kurtosisPackage *generated.KurtosisPackage) error {
	store.packages[kurtosisPackageLocator] = kurtosisPackage
	return nil
}

func (store *InMemoryStore) DeletePackage(_ context.Context, packageLocator string) error {
	delete(store.packages, packageLocator)
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
