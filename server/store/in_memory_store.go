package store

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"sort"
	"strings"
	"time"
)

type InMemoryStore struct {
	lastCrawlTime                    time.Time
	lastMetricsReporterQueryDatetime time.Time

	packages map[string]*generated.KurtosisPackage
}

func newInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		lastCrawlTime:                    time.Time{},
		lastMetricsReporterQueryDatetime: time.Time{},
		packages:                         map[string]*generated.KurtosisPackage{},
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

func (store *InMemoryStore) UpdateLastCrawlDatetime(_ context.Context, lastCrawlTime time.Time) error {
	store.lastCrawlTime = lastCrawlTime
	return nil
}

func (store *InMemoryStore) GetLastCrawlDatetime(_ context.Context) (time.Time, error) {
	return store.lastCrawlTime, nil
}

func (store *InMemoryStore) UpdateLastMetricsQueryDatetime(ctx context.Context, lastMetricsQueryTime time.Time) error {
	store.lastMetricsReporterQueryDatetime = lastMetricsQueryTime
	return nil
}

func (store *InMemoryStore) GetLastMetricsQueryDatetime(_ context.Context) (time.Time, error) {
	return store.lastMetricsReporterQueryDatetime, nil
}
