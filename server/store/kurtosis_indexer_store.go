package store

import "github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"

type KurtosisIndexerStore interface {
	GetKurtosisPackages() ([]*generated.KurtosisPackage, error)

	UpsertPackage(kurtosisPackage *generated.KurtosisPackage) error

	DeletePackage(packageName KurtosisPackageIdentifier) error
}
