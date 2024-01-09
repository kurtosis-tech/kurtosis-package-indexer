package catalog

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPackagesCatalog_Success(t *testing.T) {
	newPackagesCatalog, err := GetPackagesCatalog()
	require.NoError(t, err)

	require.NotNil(t, newPackagesCatalog)
	require.NotNil(t, newPackagesCatalog)
	require.Greater(t, len(newPackagesCatalog), 0)
}
