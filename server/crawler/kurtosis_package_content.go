package crawler

import "github.com/kurtosis-tech/kurtosis-package-indexer/server/store"

type KurtosisPackageContent struct {
	Identifier store.KurtosisPackageIdentifier

	Description string

	Url string

	PackageArguments []*StarlarkFunctionArgument

	Stars uint64
}
