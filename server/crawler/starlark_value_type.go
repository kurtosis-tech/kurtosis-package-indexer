package crawler

//go:generate go run github.com/dmarkham/enumer -type=StarlarkValueType -trimprefix=StarlarkValueType_
type StarlarkValueType int

const (
	StarlarkValueType_String StarlarkValueType = iota
	StarlarkValueType_Bool
	StarlarkValueType_Int
	StarlarkValueType_Float
	StarlarkValueType_Dict
	StarlarkValueType_Json
)
