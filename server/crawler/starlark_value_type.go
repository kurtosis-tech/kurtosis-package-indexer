package crawler

import "github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"

//go:generate go run github.com/dmarkham/enumer -type=StarlarkValueType -trimprefix=StarlarkValueType_
type StarlarkValueType int

// The name of each value here should correspond to the type name used in Starlark docstrings. For example, users
// will write `my_type (int)` so `int` should be in this enum.
// When updating this enum, the below toApiType() function needs to be updated as well
const (
	StarlarkValueType_String StarlarkValueType = iota
	StarlarkValueType_Bool
	StarlarkValueType_Int
	StarlarkValueType_Float
	StarlarkValueType_List
	StarlarkValueType_Dict
	StarlarkValueType_Json
)

// The consistency between the above enum and this conversion function is guaranteed by the
// `TestPerfectMAtchBetweenApiAndInternalType` unit test
func (starlarkValueType StarlarkValueType) toApiType() generated.ArgumentValueType {
	switch starlarkValueType {
	case StarlarkValueType_String:
		return generated.ArgumentValueType_STRING
	case StarlarkValueType_Bool:
		return generated.ArgumentValueType_BOOL
	case StarlarkValueType_Int:
		return generated.ArgumentValueType_INTEGER
	case StarlarkValueType_Float:
		return generated.ArgumentValueType_FLOAT
	case StarlarkValueType_List:
		return generated.ArgumentValueType_LIST
	case StarlarkValueType_Dict:
		return generated.ArgumentValueType_DICT
	case StarlarkValueType_Json:
		return generated.ArgumentValueType_JSON
	}
	return -1
}
