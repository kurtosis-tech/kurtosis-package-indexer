package store

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"time"
)

// KurtosisIndexerStore is the contract for the indexer store
// I left a comment to remember that there were implementations for etcd and boltdb but, both were removed
// because they were not being used in production, simply to remember that both can be added again by checking the commit history
type KurtosisIndexerStore interface {
	Close() error

	// GetKurtosisPackages returns the entire list of Kurtosis packages currently stored by the indexer
	GetKurtosisPackages(ctx context.Context) ([]*generated.KurtosisPackage, error)

	// UpsertPackage either insert or update the Kurtosis package information stored with the locator
	// `kurtosisPackageLocator`
	UpsertPackage(ctx context.Context, kurtosisPackageLocator string, kurtosisPackage *generated.KurtosisPackage) error

	// DeletePackage deletes the Kurtosis package information stored with the locator `kurtosisPackageLocator`.
	// It no-ops if no packages is stored for this locator
	DeletePackage(ctx context.Context, kurtosisPackageLocator string) error

	// UpdateLastMainCrawlDatetime updates the date time at which the last main crawling happened.
	// It is helpful to store this information so that the indexer doesn't systematically crawl everytime it is
	// restarted. Note though that to fully benefit from this, the indexer needs to be run with a persistent store (
	// either bolt with a persistent volume, or etcd)
	UpdateLastMainCrawlDatetime(ctx context.Context, lastCrawlTime time.Time) error

	// GetLastMainCrawlDatetime returns the datetime at which the last crawling happened. If no datetime is currently
	// stored (b/c no crawling has ever happened, or the store is not persistent), it returns time.Time{}
	// (i.e. the zero value for time)
	GetLastMainCrawlDatetime(ctx context.Context) (time.Time, error)

	// UpdateLastSecondaryCrawlDatetime updates the date time at which the last secondary crawling happened.
	UpdateLastSecondaryCrawlDatetime(ctx context.Context, lastCrawlTime time.Time) error

	// GetLastSecondaryCrawlDatetime returns the datetime at which the last crawling happened. If no datetime is currently
	// stored (b/c no crawling has ever happened, or the store is not persistent), it returns time.Time{}
	// (i.e. the zero value for time)
	GetLastSecondaryCrawlDatetime(ctx context.Context) (time.Time, error)

	// UpdateLastMetricsQueryDatetime updates the date time at which the last metrics query happened.
	// It is helpful to store this information so that the metrics reporter doesn't systematically query everytime it is
	// restarted. Note though that to fully benefit from this, the indexer needs to be run with a persistent store (
	// either bolt with a persistent volume, or etcd)
	UpdateLastMetricsQueryDatetime(ctx context.Context, lastMetricsQueryTime time.Time) error

	// GetLastMetricsQueryDatetime returns the datetime at which the last metrics query request happened. If no datetime is currently
	// stored (b/c no query has ever happened, or the store is not persistent), it returns time.Time{}
	// (i.e. the zero value for time)
	GetLastMetricsQueryDatetime(ctx context.Context) (time.Time, error)

	// GetPackagesRunCount returns the latest metrics package run count stored value
	// this value is generated and stored by the metrics reporter
	GetPackagesRunCount(ctx context.Context) (types.PackagesRunCount, error)

	// UpdatePackagesRunCount updates the packages run count value to the latest obtained from the metrics reporter
	UpdatePackagesRunCount(ctx context.Context, newPackagesRunCount types.PackagesRunCount) error
}
