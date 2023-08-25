package crawler

type KurtosisMainDotStar struct {
	Arguments []*StarlarkFunctionArgument
}

type StarlarkFunctionArgument struct {
	Name       string
	Type       string
	IsRequired bool
}
