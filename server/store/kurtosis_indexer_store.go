package store

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

const (
	boltDatabaseFilePathEnvVarName = "BOLT_DATABASE_FILE_PATH"
)

type KurtosisIndexerStore interface {
	Close() error

	GetKurtosisPackages() ([]*generated.KurtosisPackage, error)

	UpsertPackage(kurtosisPackage *generated.KurtosisPackage) error

	DeletePackage(packageName KurtosisPackageIdentifier) error

	UpdateLastCrawlDatetime(lastCrawlTime time.Time) error

	GetLastCrawlDatetime() (time.Time, error)
}

func InstantiateStoreFromEnvVar() (KurtosisIndexerStore, error) {
	boltDatabaseFilePath := os.Getenv(boltDatabaseFilePathEnvVarName)
	if boltDatabaseFilePath == "" {
		// env var not set, default to in-memory store
		logrus.Infof("Environment variable '%s' was not set, defaulting to in-memory store", boltDatabaseFilePathEnvVarName)
		return newInMemoryStore(), nil
	}

	boltStore, err := createBboltStore(boltDatabaseFilePath)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Error creating Bbolt store")
	}
	return boltStore, nil
}
