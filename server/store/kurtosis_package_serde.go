package store

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/stacktrace"
	"google.golang.org/protobuf/proto"
	"strings"
)

const (
	metadataKeyPrefix    = "metadata__"
	packageNameKeyPrefix = "package__"
)

func serializePackageInfo(packageName string, packageInfo *generated.KurtosisPackage) ([]byte, error) {
	serializedPackageInfo, err := proto.Marshal(packageInfo)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred serializing Kurtosis package info for package '%s'", packageName)
	}
	return serializedPackageInfo, nil
}

func deserializePackageInfo(packageName string, serializedPackageInfo []byte) (*generated.KurtosisPackage, error) {
	kurtosisPackageInfo := new(generated.KurtosisPackage)
	if err := proto.Unmarshal(serializedPackageInfo, kurtosisPackageInfo); err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred deserializing Kurtosis package info for package '%s'", packageName)
	}
	return kurtosisPackageInfo, nil
}

func compositePackageKey(packageName string) []byte {
	return []byte(fmt.Sprintf("%s%s", packageNameKeyPrefix, packageName))
}

func packageNameFromKey(key []byte) (string, error) {
	keyAsString := string(key)
	if !strings.HasPrefix(keyAsString, packageNameKeyPrefix) {
		return "", stacktrace.NewError("Kurtosis package key seems invalid: '%s'", keyAsString)
	}
	packageName := strings.TrimPrefix(keyAsString, packageNameKeyPrefix)
	return packageName, nil
}

func compositeMetadataKey(key string) []byte {
	return []byte(fmt.Sprintf("%s%s", metadataKeyPrefix, key))
}
