package store

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

const (
	boltDatabaseFilePathEnvVarName = "BOLT_DATABASE_FILE_PATH"
	etcdDatabaseUrlsEnvVarName     = "ETCD_DATABASE_URLS"
)

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

	// UpdateLastCrawlDatetime updates the date time at which the last crawling happened.
	// It is helpful to store this information so that the indexer doesn't systematically crawl everytime it is
	// restarted. Note though that to fully benefit from this, the indexer needs to be run with a persistent store (
	// either bolt with a persistent volume, or etcd)
	UpdateLastCrawlDatetime(ctx context.Context, lastCrawlTime time.Time) error

	// GetLastCrawlDatetime returns the datetime at which the last crawling happened. If no datetime is currently
	// stored (b/c no crawling has ever happened, or the store is not persistent), it returns time.Time{}
	// (i.e. the zero value for time)
	GetLastCrawlDatetime(ctx context.Context) (time.Time, error)
}

func InstantiateStoreFromEnvVar() (KurtosisIndexerStore, error) {
	etcdDatabaseUrls := os.Getenv(etcdDatabaseUrlsEnvVarName)
	if etcdDatabaseUrls != "" {
		logrus.Infof("Environment variable '%s' recognized ('%s'). Instanciating etcd database client",
			etcdDatabaseUrlsEnvVarName, etcdDatabaseUrls)
		var etcdUrls []string
		for _, url := range strings.Split(etcdDatabaseUrls, ",") {
			etcdUrls = append(etcdUrls, strings.TrimSpace(url))
		}
		etcdStore, err := createEtcdBackedStore(etcdUrls)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Error creating store backed by etcd")
		}
		return etcdStore, nil
	}

	boltDatabaseFilePath := os.Getenv(boltDatabaseFilePathEnvVarName)
	if boltDatabaseFilePath != "" {
		logrus.Infof("Environment variable '%s' recognized ('%s'). Instanciating Bbolt database",
			boltDatabaseFilePathEnvVarName, boltDatabaseFilePath)
		boltStore, err := createBboltStore(boltDatabaseFilePath)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Error creating Bbolt store")
		}
		return boltStore, nil
	}
	// env var not set, default to in-memory store
	logrus.Infof("No environment variable set for the store, defaulting to in-memory store")
	return newInMemoryStore(), nil
}
