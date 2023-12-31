package crawler

type KurtosisMainDotStar struct {
	Description string

	Arguments []*StarlarkFunctionArgument

	ReturnDescription string
}

type StarlarkFunctionArgument struct {
	Name        string
	Description string
	Type        *StarlarkArgumentType
	IsRequired  bool

	// raw string of the default value raw if one exists, otherwise nil
	DefaultValue *string
}

type StarlarkArgumentType struct {
	Type StarlarkValueType

	InnerType1 *StarlarkValueType

	InnerType2 *StarlarkValueType
}
