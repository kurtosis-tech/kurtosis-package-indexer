package store

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSerializeAndDeserializePackagesRunCount_Success(t *testing.T) {
	packagesRunCount := types.PackagesRunCount{
		"github.com/kurtosis-tech/ethereum-package": 1765,
		"github.com/kurtosis-tech/postgres-package": 84,
		"github.com/kurtosis-tech/redis-package":    23,
	}

	serializedPackagesRunCount, err := serializePackagesRunCount(packagesRunCount)
	require.NoError(t, err)

	actualPackagesRunCount, err := deserializePackagesRunCount(serializedPackagesRunCount)
	require.NoError(t, err)

	require.EqualValues(t, packagesRunCount, actualPackagesRunCount)
}
