package crawler

import "github.com/kurtosis-tech/kurtosis-package-indexer/server/store"

type KurtosisPackageContent struct {
	Identifier store.KurtosisPackageIdentifier

	PackageArguments []*StarlarkFunctionArgument

	Stars uint64
}
