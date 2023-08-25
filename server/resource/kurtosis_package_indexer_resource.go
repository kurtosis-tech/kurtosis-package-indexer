package resource

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/api_constructors"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	"github.com/kurtosis-tech/stacktrace"
	"google.golang.org/protobuf/types/known/emptypb"
)

type KurtosisPackageIndexer struct {
	store store.KurtosisIndexerStore
}

func NewKurtosisPackageIndexer(store store.KurtosisIndexerStore) *KurtosisPackageIndexer {
	return &KurtosisPackageIndexer{
		store: store,
	}
}

func (resource *KurtosisPackageIndexer) Ping(_ context.Context, _ *generated.IndexerPing) (*generated.IndexerPong, error) {
	return &generated.IndexerPong{}, nil
}

func (resource *KurtosisPackageIndexer) GetPackages(_ context.Context, _ *emptypb.Empty) (*generated.GetPackagesResponse, error) {
	packages, err := resource.store.GetKurtosisPackages()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred fetching the packages from the store")
	}
	return api_constructors.NewGetPackagesResponse(packages...), nil
}
