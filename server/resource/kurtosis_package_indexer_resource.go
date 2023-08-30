package resource

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/api_constructors"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/crawler"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	"github.com/kurtosis-tech/stacktrace"
	"google.golang.org/protobuf/types/known/emptypb"
)

type KurtosisPackageIndexer struct {
	store store.KurtosisIndexerStore

	crawler *crawler.GithubCrawler
}

func NewKurtosisPackageIndexer(store store.KurtosisIndexerStore, crawler *crawler.GithubCrawler) *KurtosisPackageIndexer {
	return &KurtosisPackageIndexer{
		store:   store,
		crawler: crawler,
	}
}

func (resource *KurtosisPackageIndexer) IsAvailable(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (resource *KurtosisPackageIndexer) GetPackages(_ context.Context, _ *emptypb.Empty) (*generated.GetPackagesResponse, error) {
	packages, err := resource.store.GetKurtosisPackages()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred fetching the packages from the store")
	}
	return api_constructors.NewGetPackagesResponse(packages...), nil
}

func (resource *KurtosisPackageIndexer) Reindex(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := resource.crawler.Schedule(ctx); err != nil {
		return nil, stacktrace.Propagate(err, "Error updating crawler schedule to kick off a reindex immediately")
	}
	return &emptypb.Empty{}, nil
}
