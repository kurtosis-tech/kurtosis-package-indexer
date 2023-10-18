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

func (resource *KurtosisPackageIndexer) GetPackages(ctx context.Context, _ *emptypb.Empty) (*generated.GetPackagesResponse, error) {
	packages, err := resource.store.GetKurtosisPackages(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred fetching the packages from the store")
	}
	return api_constructors.NewGetPackagesResponse(packages...), nil
}

func (resource *KurtosisPackageIndexer) Reindex(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	if err := resource.crawler.Schedule(true); err != nil {
		return nil, stacktrace.Propagate(err, "Error updating crawler schedule to kick off a reindex immediately")
	}
	return &emptypb.Empty{}, nil
}

func (resource *KurtosisPackageIndexer) ReadPackage(ctx context.Context, input *generated.ReadPackageRequest) (*generated.ReadPackageResponse, error) {
	if input == nil {
		return nil, stacktrace.NewError("an empty input was provided")
	}
	githubClient, err := crawler.CreateGithubClient(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred while creating the github client")
	}
	pack, err := crawler.ReadPackage(ctx, githubClient, input.GetRepositoryMetadata())
	return api_constructors.NewReadPackageResponse(pack), nil
}
