package catalog

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreatePackageDataFromCatalogFromPackageName_Success(t *testing.T) {
	// TODO use constants
	packageNameTestList := []types.PackageName{
		types.PackageName("github.com/kurtosis-tech/postgres-package"),
		types.PackageName("github.com/kurtosis-tech/awesome-kurtosis/data-package"),
		types.PackageName("github.com/kurtosis-tech/awesome-kurtosis/flink-kafka-example/kurtosis-package"),
	}

	expectedRepositoryOwner := "kurtosis-tech"

	expectedRepositoryName := []string{
		"postgres-package",
		"awesome-kurtosis",
		"awesome-kurtosis",
	}

	expectedRootPath := []string{
		"/",
		"data-package/",
		"flink-kafka-example/kurtosis-package/",
	}

	for packageIndex, packageName := range packageNameTestList {
		newPackageRepositoryMainMetadata, err := createPackageDataFromCatalogFromPackageName(packageName)
		require.NoError(t, err)
		require.Equal(t, expectedRepositoryOwner, newPackageRepositoryMainMetadata.GetRepositoryOwner())
		require.Equal(t, expectedRepositoryName[packageIndex], newPackageRepositoryMainMetadata.GetRepositoryName())
		require.Equal(t, expectedRootPath[packageIndex], newPackageRepositoryMainMetadata.GetRepositoryPackageRootPath())
	}
}
