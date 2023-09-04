package store

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"time"
)

type InMemoryStore struct {
	lastCrawlTime time.Time

	packages map[KurtosisPackageIdentifier]*generated.KurtosisPackage
}

func newInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		lastCrawlTime: time.Time{},
		packages:      map[KurtosisPackageIdentifier]*generated.KurtosisPackage{},
	}
}

func (store *InMemoryStore) Close() error {
	return nil
}

func (store *InMemoryStore) GetKurtosisPackages() ([]*generated.KurtosisPackage, error) {
	var packages []*generated.KurtosisPackage
	for _, kurtosisPackage := range store.packages {
		packages = append(packages, kurtosisPackage)
	}
	return packages, nil
}

func (store *InMemoryStore) UpsertPackage(kurtosisPackage *generated.KurtosisPackage) error {
	store.packages[KurtosisPackageIdentifier(kurtosisPackage.GetName())] = kurtosisPackage
	return nil
}

func (store *InMemoryStore) DeletePackage(packageName KurtosisPackageIdentifier) error {
	delete(store.packages, packageName)
	return nil
}

func (store *InMemoryStore) UpdateLastCrawlDatetime(lastCrawlTime time.Time) error {
	store.lastCrawlTime = lastCrawlTime
	return nil
}

func (store *InMemoryStore) GetLastCrawlDatetime() (time.Time, error) {
	return store.lastCrawlTime, nil
}
