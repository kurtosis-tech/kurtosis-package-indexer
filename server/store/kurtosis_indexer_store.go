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

	GetKurtosisPackages(ctx context.Context) ([]*generated.KurtosisPackage, error)

	UpsertPackage(ctx context.Context, kurtosisPackage *generated.KurtosisPackage) error

	DeletePackage(ctx context.Context, packageName KurtosisPackageIdentifier) error

	UpdateLastCrawlDatetime(ctx context.Context, lastCrawlTime time.Time) error

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
