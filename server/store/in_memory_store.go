package store

import "github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"

type InMemoryStore struct {
	packages map[KurtosisPackageIdentifier]*generated.KurtosisPackage
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		packages: map[KurtosisPackageIdentifier]*generated.KurtosisPackage{},
	}
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
