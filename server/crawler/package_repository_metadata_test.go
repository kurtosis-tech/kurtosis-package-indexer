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
