package store

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/stacktrace"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
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

func (store *BboltStore) GetKurtosisPackages() ([]*generated.KurtosisPackage, error) {
	var packages []*generated.KurtosisPackage
	err := store.db.View(func(tx *bbolt.Tx) error {
		bucketInstance := tx.Bucket([]byte(kurtosisPackagesBucketName))
		if bucketInstance == nil {
			return nil
		}
		bucketCursor := bucketInstance.Cursor()
		for packageName, serializedPackageInfo := bucketCursor.First(); packageName != nil; packageName, serializedPackageInfo = bucketCursor.Next() {
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
	return packages, nil
}

func (store *BboltStore) UpsertPackage(kurtosisPackage *generated.KurtosisPackage) error {
	err := store.db.Update(func(tx *bbolt.Tx) error {
		kurtosisPackageName := kurtosisPackage.GetName()
		serializedPackageInfo, err := serializePackageInfo(kurtosisPackage.GetName(), kurtosisPackage)
		if err != nil {
			return stacktrace.Propagate(err, "An error occurred serializing package information '%s' prior to persisting it to the Bbolt database", kurtosisPackageName)
		}
		bucketInstance, err := tx.CreateBucketIfNotExists([]byte(kurtosisPackagesBucketName))
		if err != nil {
			return stacktrace.Propagate(err, "Unable to create or get bucket from Bbolt database")
		}
		if err := bucketInstance.Put([]byte(kurtosisPackageName), serializedPackageInfo); err != nil {
			return stacktrace.Propagate(err, "An error occurred persisting serialized package info to Bbolt database for package '%s'", kurtosisPackageName)
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred updating package in Bbolt database")
	}
	return nil
}

func (store *BboltStore) DeletePackage(packageName KurtosisPackageIdentifier) error {
	err := store.db.Update(func(tx *bbolt.Tx) error {
		bucketInstance := tx.Bucket([]byte(kurtosisPackagesBucketName))
		if bucketInstance == nil {
			// no bucket means not data, nothing to delete
			return nil
		}
		if err := bucketInstance.Delete([]byte(packageName)); err != nil {
			return stacktrace.Propagate(err, "An error occurred deleting package info from Bbolt database for package '%s'", packageName)
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred deleting package '%s' from Bbolt database", packageName)
	}
	return nil
}

func (store *BboltStore) UpdateLastCrawlDatetime(lastCrawlDatetime time.Time) error {
	err := store.db.Update(func(tx *bbolt.Tx) error {
		serializedDatetime, err := lastCrawlDatetime.MarshalBinary()
		if err != nil {
			return stacktrace.Propagate(err, "An error occurred serializing time information '%v' prior to persisting it to the Bbolt database", lastCrawlDatetime)
		}

		bucketInstance, err := tx.CreateBucketIfNotExists([]byte(lastCrawlMetadataBucketName))
		if err != nil {
			return stacktrace.Propagate(err, "Unable to create or get bucket from Bbolt database")
		}
		if err := bucketInstance.Put([]byte(lastCrawlDatetimeKey), serializedDatetime); err != nil {
			return stacktrace.Propagate(err, "An error occurred persisting serialized time info to Bbolt database")
		}
		return nil
	})
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred updating the last crawl datetime in Bbolt database")
	}
	return nil
}

func (store *BboltStore) GetLastCrawlDatetime() (time.Time, error) {
	lastCrawlDatetime := time.Time{}
	err := store.db.View(func(tx *bbolt.Tx) error {
		bucketInstance := tx.Bucket([]byte(lastCrawlMetadataBucketName))
		if bucketInstance == nil {
			// no bucket means not data
			return nil
		}
		serializedDatetime := bucketInstance.Get([]byte(lastCrawlDatetimeKey))
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

func serializePackageInfo(packageName string, packageInfo *generated.KurtosisPackage) ([]byte, error) {
	serializedPackageInfo, err := proto.Marshal(packageInfo)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred serializing Kurtosis package info for package '%s'", packageName)
	}
	return serializedPackageInfo, nil
}

func deserializePackageInfo(packageName []byte, serializedPackageInfo []byte) (*generated.KurtosisPackage, error) {
	kurtosisPackageInfo := new(generated.KurtosisPackage)
	if err := proto.Unmarshal(serializedPackageInfo, kurtosisPackageInfo); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred deserializing Kurtosis package info for package '%s'", string(packageName))
	}
	return kurtosisPackageInfo, nil
}
