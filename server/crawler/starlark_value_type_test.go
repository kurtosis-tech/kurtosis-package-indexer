package crawler

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPerfectMAtchBetweenApiAndInternalType(t *testing.T) {
	starlarkValueTypeValues := StarlarkValueTypeValues()
	apiArgumentValueTypeValues := generated.ArgumentValueType_name

	require.Equalf(t, len(starlarkValueTypeValues), len(apiArgumentValueTypeValues),
		"The ArgumentValueType API enum does not contain the same number of values than the internal "+
			"StarlarkValueType enum")

	for _, starlarkValueType := range starlarkValueTypeValues {
		convertedType := starlarkValueType.toApiType()
		_, ok := apiArgumentValueTypeValues[int32(convertedType)]
		require.Truef(t, ok, "Type '%v' couldn't be converted to an API type", starlarkValueType)
	}
}
