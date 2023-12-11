package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
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

func getPackagesRunCountKeyBytes() []byte {
	return []byte(packagesRunCountKey)
}

func serializePackagesRunCount(packagesRunCount types.PackagesRunCount) ([]byte, error) {

	packagesRunCountBytesBuffer := new(bytes.Buffer)

	encoder := gob.NewEncoder(packagesRunCountBytesBuffer)

	if err := encoder.Encode(packagesRunCount); err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred encoding the packages run count map")
	}

	serializedPackagesRunCount := packagesRunCountBytesBuffer.Bytes()

	return serializedPackagesRunCount, nil
}

func deserializePackagesRunCount(serializedPackagesRunCount []byte) (types.PackagesRunCount, error) {
	packagesRunCount := types.PackagesRunCount{}

	bytesReader := bytes.NewReader(serializedPackagesRunCount)

	decoder := gob.NewDecoder(bytesReader)

	if err := decoder.Decode(&packagesRunCount); err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred decoding the packages run count map")
	}

	return packagesRunCount, nil
}
