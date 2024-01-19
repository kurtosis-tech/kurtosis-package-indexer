package crawler

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetLocator_Success(t *testing.T) {

	packageRepositoryMetadata := NewPackageRepositoryMetadata(
		"kurtosis-tech",
		"datastore-army-package",
		"/",
		"",
		0,
		time.Time{},
		"",
	)

	expectedLocator := "kurtosis-tech/datastore-army-package/"

	require.Equal(t, expectedLocator, packageRepositoryMetadata.GetLocator())

	packageRepositoryMetadata = NewPackageRepositoryMetadata(
		"kurtosis-tech",
		"awesome-kurtosis",
		"flink-kafka-example/kurtosis-package/",
		"",
		0,
		time.Time{},
		"",
	)

	expectedLocator = "kurtosis-tech/awesome-kurtosis/flink-kafka-example/kurtosis-package/"

	require.Equal(t, expectedLocator, packageRepositoryMetadata.GetLocator())

}

func TestGetDownloadRootURL_Success(t *testing.T) {
	packageRepositoryMetadata := NewPackageRepositoryMetadata(
		"kurtosis-tech",
		"datastore-army-package",
		"/",
		"",
		0,
		time.Time{},
		"",
	)

	expectedPackageRepositoryDownloadURL := "https://raw.githubusercontent.com/kurtosis-tech/datastore-army-package/"

	downloadRootURL, err := packageRepositoryMetadata.GetDownloadRootURL()
	require.NoError(t, err)
	require.Equal(t, expectedPackageRepositoryDownloadURL, downloadRootURL)

	packageRepositoryMetadata = NewPackageRepositoryMetadata(
		"kurtosis-tech",
		"web3-tools",
		"/package",
		"",
		0,
		time.Time{},
		"master",
	)

	expectedPackageRepositoryDownloadURL = "https://raw.githubusercontent.com/kurtosis-tech/web3-tools/master/package"

	downloadRootURL, err = packageRepositoryMetadata.GetDownloadRootURL()
	require.NoError(t, err)
	require.Equal(t, expectedPackageRepositoryDownloadURL, downloadRootURL)
}

func TestGetDownloadRootURL_WrongBranch(t *testing.T) {
	packageRepositoryMetadata := NewPackageRepositoryMetadata(
		"kurtosis-tech",
		"datastore-army-package",
		"/",
		"",
		0,
		time.Time{},
		"/wrong-branch",
	)

	expectedPackageRepositoryDownloadURL := "https://raw.githubusercontent.com/kurtosis-tech/datastore-army-package/"

	downloadRootURL, err := packageRepositoryMetadata.GetDownloadRootURL()
	require.NoError(t, err)
	require.NotEqual(t, expectedPackageRepositoryDownloadURL, downloadRootURL)
}
