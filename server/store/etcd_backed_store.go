package store

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/stacktrace"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sort"
	"strings"
	"time"
)

const (
	etcdDialTimeout = 5 * time.Second
)

type EtcdBackedStore struct {
	client *clientv3.Client
}

func createEtcdBackedStore(etcdUrls []string) (*EtcdBackedStore, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:            etcdUrls,
		AutoSyncInterval:     0,
		DialTimeout:          etcdDialTimeout,
		DialKeepAliveTime:    0,
		DialKeepAliveTimeout: 0,
		MaxCallSendMsgSize:   0,
		MaxCallRecvMsgSize:   0,
		TLS:                  nil,
		Username:             "",
		Password:             "",
		RejectOldCluster:     false,
		DialOptions:          nil,
		Context:              nil,
		Logger:               nil,
		LogConfig:            nil,
		PermitWithoutStream:  false,
	})
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred creating etcd client")
	}
	return &EtcdBackedStore{
		client: client,
	}, nil
}

func (store *EtcdBackedStore) Close() error {
	if err := store.client.Close(); err != nil {
		return stacktrace.Propagate(err, "An error occurred closing etcd client")
	}
	return nil
}

func (store *EtcdBackedStore) GetKurtosisPackages(ctx context.Context) ([]*generated.KurtosisPackage, error) {
	var packages []*generated.KurtosisPackage
	resp, err := store.client.Get(ctx, packageNameKeyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, stacktrace.Propagate(err, "Error retrieving values from etcd")
	}

	for _, keyValue := range resp.Kvs {
		kurtosisPackageName, err := packageNameFromKey(keyValue.Key)
		if err != nil {
			return nil, stacktrace.Propagate(err, "The key couldn't be converted to a package name. This is unexpected")
		}
		kurtosisPackage, err := deserializePackageInfo(kurtosisPackageName, keyValue.Value)
		if err != nil {
			return nil, stacktrace.Propagate(err, "An error occurred deserializing package info")
		}
		packages = append(packages, kurtosisPackage)
	}
	sort.SliceStable(packages, func(i, j int) bool {
		return strings.ToLower(packages[i].GetUrl()) < strings.ToLower(packages[j].GetUrl())
	})
	return packages, nil
}

func (store *EtcdBackedStore) UpsertPackage(ctx context.Context, kurtosisPackageLocator string, kurtosisPackage *generated.KurtosisPackage) error {
	kurtosisPackageKey := compositePackageKey(kurtosisPackageLocator)
	serializedKurtosisPackage, err := serializePackageInfo(kurtosisPackageLocator, kurtosisPackage)
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred serializing Kurtosis package info")
	}

	_, err = store.client.Put(ctx, string(kurtosisPackageKey), string(serializedKurtosisPackage))
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred updating package in etcd database: '%s'",
			kurtosisPackage.GetName())
	}
	return nil
}

func (store *EtcdBackedStore) DeletePackage(ctx context.Context, kurtosisPackageLocator string) error {
	kurtosisPackageKey := compositePackageKey(kurtosisPackageLocator)
	_, err := store.client.Delete(ctx, string(kurtosisPackageKey))
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred deleting package in etcd database: '%s'",
			kurtosisPackageLocator)
	}
	return nil
}

func (store *EtcdBackedStore) UpdateLastCrawlDatetime(ctx context.Context, lastCrawlDatetime time.Time) error {
	serializedDatetime, err := lastCrawlDatetime.MarshalBinary()
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred serializing time information '%v' prior to persisting it to the etcd database", lastCrawlDatetime)
	}

	lastCrawlDatetimeKeyBytes := compositeMetadataKey(lastCrawlDatetimeKey)
	_, err = store.client.Put(ctx, string(lastCrawlDatetimeKeyBytes), string(serializedDatetime))
	if err != nil {
		return stacktrace.Propagate(err, "An error occurred persisting serialized time info to etcd database")
	}
	return nil
}

func (store *EtcdBackedStore) GetLastCrawlDatetime(ctx context.Context) (time.Time, error) {
	lastCrawlDatetime := time.Time{}
	lastCrawlDatetimeKeyBytes := compositeMetadataKey(lastCrawlDatetimeKey)
	resp, err := store.client.Get(ctx, string(lastCrawlDatetimeKeyBytes))
	if err != nil {
		return lastCrawlDatetime, stacktrace.Propagate(err, "An error occurred retrieving metadata from etcd database")
	}

	if len(resp.Kvs) == 0 {
		// no data stored yet
		return lastCrawlDatetime, nil
	}
	if len(resp.Kvs) > 1 {
		return lastCrawlDatetime, stacktrace.Propagate(err, "More than one value is stored with the key '%s' in etcd. This is unexpected", string(lastCrawlDatetimeKeyBytes))
	}
	serializedDatetime := resp.Kvs[0].Value
	if err = lastCrawlDatetime.UnmarshalBinary(serializedDatetime); err != nil {
		return lastCrawlDatetime, stacktrace.Propagate(err, "An error occurred unmarshalling datetime retrieved from etcd database")
	}
	return lastCrawlDatetime, nil
}
