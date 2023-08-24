package resource

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
)

type KurtosisPackageIndexer struct {
}

func NewKurtosisPackageIndexer() *KurtosisPackageIndexer {
	return &KurtosisPackageIndexer{}
}

func (resource KurtosisPackageIndexer) Ping(ctx context.Context, ping *generated.IndexerPing) (*generated.IndexerPong, error) {
	return &generated.IndexerPong{}, nil
}
