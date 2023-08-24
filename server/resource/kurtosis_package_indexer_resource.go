package resource

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/api_constructors"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"google.golang.org/protobuf/types/known/emptypb"
)

type KurtosisPackageIndexer struct {
}

func NewKurtosisPackageIndexer() *KurtosisPackageIndexer {
	return &KurtosisPackageIndexer{}
}

func (resource *KurtosisPackageIndexer) Ping(_ context.Context, _ *generated.IndexerPing) (*generated.IndexerPong, error) {
	return &generated.IndexerPong{}, nil
}

func (resource *KurtosisPackageIndexer) GetPackages(_ context.Context, _ *emptypb.Empty) (*generated.GetPackagesResponse, error) {
	// TODO: implement
	hyperlanePackage := api_constructors.NewKurtosisPackage(
		"github.com/kurtosis-tech/hyperlane-package",
		api_constructors.NewUntypedPackageArg("origin_chain", true),
		api_constructors.NewUntypedPackageArg("remote_chains", true),
		api_constructors.NewUntypedPackageArg("aws_env", false),
	)
	return api_constructors.NewGetPackagesResponse(hyperlanePackage), nil
}
