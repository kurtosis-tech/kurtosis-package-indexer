package store

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/stacktrace"
	"go.etcd.io/bbolt"
	"sort"
	"strings"
	"time"
)

const (
	backingFilePerm = 0600

	lastCrawlMetadataBucketName = "lastCrawlMetadataBucketName"
	lastCrawlDatetimeKey        = "datetime"

	kurtosisPackagesBucketName = "kurtosisPackages"
)

type BboltStore struct {
	db *bbolt.DB
}

func createBboltStore(backingFilePath string) (*BboltStore, error) {
	opts := new(bbolt.Options)
	db, err := bbolt.Open(backingFilePath, backingFilePerm, opts)
	if err != nil {
		return nil, stacktrace.Propagate(err, "Unable to init Bbolt db with at '%s'", backingFilePath)
	}
	return &BboltStore{
		db: db,
	}, nil
}

func (store *BboltStore) Close() error {
	if err := store.db.Close(); err != nil {
		return stacktrace.Propagate(err, "An error occurred closing BBolt database")
	}
	return nil
}

func (store *BboltStore) GetKurtosisPackages(_ context.Context) ([]*generated.KurtosisPackage, error) {
	var packages []*generated.KurtosisPackage
	err := store.db.View(func(tx *bbolt.Tx) error {
		bucketInstance := tx.Bucket([]byte(kurtosisPackagesBucketName))
		if bucketInstance == nil {
			return nil
		}
		bucketCursor := bucketInstance.Cursor()
		for packageNameKey, serializedPackageInfo := bucketCursor.First(); packageNameKey != nil; packageNameKey, serializedPackageInfo = bucketCursor.Next() {
			packageName, err := packageNameFromKey(packageNameKey)
			if err != nil {
				return stacktrace.Propagate(err, "Unable to compute package name from the key: '%s'", string(packageNameKey))
			}
			packageInfo, err := deserializePackageInfo(packageName, serializedPackageInfo)
			if err != nil {
				// TODO: maybe try the other packages before exiting, who knows...
				return stacktrace.Propagate(err, "An error occurred deserializing Kurtosis package information store in Bbolt database for package '%s'", string(packageName))
			}
			packages = append(packages, packageInfo)
		}
		return nil
	})
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred retrieving package from Bbolt database")
	}
	sort.SliceStable(packages, func(i, j int) bool {
		return strings.ToLower(packages[i].GetUrl()) < strings.ToLower(packages[j].GetUrl())
	})
	return packages, nil
}

func (store *BboltStore) UpsertPackage(_ context.Context, kurtosisPackageLocator string, kurtosisPackage *generated.KurtosisPackage) error {
	err := store.db.Update(func(tx *bbolt.Tx) error {
		serializedPackageInfo, err := serializePackageInfo(kurtosisPackageLocator, kurtosisPackage)
		if err != nil {
			return stacktrace.Propagate(err, "An error occurred serializing package information '%s' prior to persisting it to the Bbolt database", kurtosisPackageLocator)
		}
		bucketInstance, err := tx.CreateBucketIfNotExists([]byte(kurtosisPackagesBucketName))
		if err != nil {
			return stacktrace.Propagate(err, "Unable to create or get bucket from Bbolt database")
		}
		packageNameKey := compositePackageKey(kurtosisPackageLocator)
		if err := bucketInstance.Put(packageNameKey, serializedPackageInfo); err != nil {
			return stacktrace.Propagate(err, "An error occurred persisting serialized package info to Bbolt database for package '%s'", kurtosisPackageLocator)
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred updating package in Bbolt database")
	}
	return nil
}

func (store *BboltStore) DeletePackage(_ context.Context, kurtosisPackageLocator string) error {
	err := store.db.Update(func(tx *bbolt.Tx) error {
		bucketInstance := tx.Bucket([]byte(kurtosisPackagesBucketName))
		if bucketInstance == nil {
			// no bucket means not data, nothing to delete
			return nil
		}
		packageNameKey := compositePackageKey(kurtosisPackageLocator)
		if err := bucketInstance.Delete(packageNameKey); err != nil {
			return stacktrace.Propagate(err, "An error occurred deleting package info from Bbolt database for package '%s'", kurtosisPackageLocator)
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred deleting package '%s' from Bbolt database", kurtosisPackageLocator)
	}
	return nil
}

func (store *BboltStore) UpdateLastCrawlDatetime(_ context.Context, lastCrawlDatetime time.Time) error {
	err := store.db.Update(func(tx *bbolt.Tx) error {
		serializedDatetime, err := lastCrawlDatetime.MarshalBinary()
		if err != nil {
			return stacktrace.Propagate(err, "An error occurred serializing time information '%v' prior to persisting it to the Bbolt database", lastCrawlDatetime)
		}

		bucketInstance, err := tx.CreateBucketIfNotExists([]byte(lastCrawlMetadataBucketName))
		if err != nil {
			return stacktrace.Propagate(err, "Unable to create or get bucket from Bbolt database")
		}
		lastCrawlDatetimeKeyBytes := compositeMetadataKey(lastCrawlDatetimeKey)
		if err := bucketInstance.Put(lastCrawlDatetimeKeyBytes, serializedDatetime); err != nil {
			return stacktrace.Propagate(err, "An error occurred persisting serialized time info to Bbolt database")
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred updating the last crawl datetime in Bbolt database")
	}
	return nil
}

func (store *BboltStore) GetLastCrawlDatetime(_ context.Context) (time.Time, error) {
	lastCrawlDatetime := time.Time{}
	err := store.db.View(func(tx *bbolt.Tx) error {
		bucketInstance := tx.Bucket([]byte(lastCrawlMetadataBucketName))
		if bucketInstance == nil {
			// no bucket means not data
			return nil
		}
		lastCrawlDatetimeKeyBytes := compositeMetadataKey(lastCrawlDatetimeKey)
		serializedDatetime := bucketInstance.Get(lastCrawlDatetimeKeyBytes)
		if serializedDatetime == nil {
			// no data stored yet
			return nil
		}

		if err := lastCrawlDatetime.UnmarshalBinary(serializedDatetime); err != nil {
			return stacktrace.Propagate(err, "An error occurred unmarshalling datetime retrieved from Bbolt database")
		}
		return nil
	})
	if err != nil {
		return lastCrawlDatetime, stacktrace.Propagate(err, "An error occurred retrieving the last crawl datetime from Bbolt database")
	}
	return lastCrawlDatetime, nil
}
