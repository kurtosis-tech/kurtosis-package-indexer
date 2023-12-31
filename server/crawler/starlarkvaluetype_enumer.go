// Code generated by "enumer -type=StarlarkValueType -trimprefix=StarlarkValueType_"; DO NOT EDIT.

package crawler

import (
	"fmt"
	"strings"
)

const _StarlarkValueTypeName = "StringBoolIntListDictJson"

var _StarlarkValueTypeIndex = [...]uint8{0, 6, 10, 13, 17, 21, 25}

const _StarlarkValueTypeLowerName = "stringboolintlistdictjson"

func (i StarlarkValueType) String() string {
	if i < 0 || i >= StarlarkValueType(len(_StarlarkValueTypeIndex)-1) {
		return fmt.Sprintf("StarlarkValueType(%d)", i)
	}
	return _StarlarkValueTypeName[_StarlarkValueTypeIndex[i]:_StarlarkValueTypeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _StarlarkValueTypeNoOp() {
	var x [1]struct{}
	_ = x[StarlarkValueType_String-(0)]
	_ = x[StarlarkValueType_Bool-(1)]
	_ = x[StarlarkValueType_Int-(2)]
	_ = x[StarlarkValueType_List-(3)]
	_ = x[StarlarkValueType_Dict-(4)]
	_ = x[StarlarkValueType_Json-(5)]
}

var _StarlarkValueTypeValues = []StarlarkValueType{StarlarkValueType_String, StarlarkValueType_Bool, StarlarkValueType_Int, StarlarkValueType_List, StarlarkValueType_Dict, StarlarkValueType_Json}

var _StarlarkValueTypeNameToValueMap = map[string]StarlarkValueType{
	_StarlarkValueTypeName[0:6]:        StarlarkValueType_String,
	_StarlarkValueTypeLowerName[0:6]:   StarlarkValueType_String,
	_StarlarkValueTypeName[6:10]:       StarlarkValueType_Bool,
	_StarlarkValueTypeLowerName[6:10]:  StarlarkValueType_Bool,
	_StarlarkValueTypeName[10:13]:      StarlarkValueType_Int,
	_StarlarkValueTypeLowerName[10:13]: StarlarkValueType_Int,
	_StarlarkValueTypeName[13:17]:      StarlarkValueType_List,
	_StarlarkValueTypeLowerName[13:17]: StarlarkValueType_List,
	_StarlarkValueTypeName[17:21]:      StarlarkValueType_Dict,
	_StarlarkValueTypeLowerName[17:21]: StarlarkValueType_Dict,
	_StarlarkValueTypeName[21:25]:      StarlarkValueType_Json,
	_StarlarkValueTypeLowerName[21:25]: StarlarkValueType_Json,
}

var _StarlarkValueTypeNames = []string{
	_StarlarkValueTypeName[0:6],
	_StarlarkValueTypeName[6:10],
	_StarlarkValueTypeName[10:13],
	_StarlarkValueTypeName[13:17],
	_StarlarkValueTypeName[17:21],
	_StarlarkValueTypeName[21:25],
}

// StarlarkValueTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func StarlarkValueTypeString(s string) (StarlarkValueType, error) {
	if val, ok := _StarlarkValueTypeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _StarlarkValueTypeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to StarlarkValueType values", s)
}

// StarlarkValueTypeValues returns all values of the enum
func StarlarkValueTypeValues() []StarlarkValueType {
	return _StarlarkValueTypeValues
}

// StarlarkValueTypeStrings returns a slice of all String values of the enum
func StarlarkValueTypeStrings() []string {
	strs := make([]string, len(_StarlarkValueTypeNames))
	copy(strs, _StarlarkValueTypeNames)
	return strs
}

// IsAStarlarkValueType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i StarlarkValueType) IsAStarlarkValueType() bool {
	for _, v := range _StarlarkValueTypeValues {
		if i == v {
			return true
		}
	}
	return false
}
