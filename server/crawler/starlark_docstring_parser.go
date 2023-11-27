package crawler

import (
	"github.com/kurtosis-tech/starlark-lsp/pkg/docstring"
	"regexp"
	"strings"
)

var (
	argumentNameAndTypeRegexp = regexp.MustCompile(`\s*(?P<name>[a-zA-Z0-9_]+)\s*(?P<type>\(.+\)|)`)
	dictTypeRegexp            = regexp.MustCompile(`dict\s*\[(?P<keyType>[a-zA-Z]+)\s*,\s*(?P<valueType>[a-zA-Z]+)\s*\]`)
	listTypeRegexp            = regexp.MustCompile(`list\s*\[(?P<valueType>[a-zA-Z]+)\s*\]`)
)

func ParseRunFunctionDocstring(rawDocstring string) (*KurtosisMainDotStar, error) {
	var description string
	var arguments []*StarlarkFunctionArgument
	var returns string

	parsedDocstring := docstring.Parse(rawDocstring)

	// Description
	description = parsedDocstring.Description

	// Arguments
	for _, arg := range parsedDocstring.Args() {
		argName, argType := parseNameAndType(arg.Name)

		if argName != "" {
			arguments = append(arguments, &StarlarkFunctionArgument{
				Name:         argName,
				Description:  arg.Desc,
				Type:         argType,
				IsRequired:   false,
				DefaultValue: nil,
			})
		}
	}

	// Returns
	returns = parsedDocstring.Returns()

	return &KurtosisMainDotStar{
		Description:       description,
		Arguments:         arguments,
		ReturnDescription: returns,
	}, nil
}

func parseNameAndType(argNameAndTypeFromParser string) (string, *StarlarkArgumentType) {
	var argName string
	var rawArgType string
	if argumentNameAndTypeRegexp.MatchString(argNameAndTypeFromParser) {
		nameAndTypeGroups := argumentNameAndTypeRegexp.FindStringSubmatch(argNameAndTypeFromParser)
		// nolint: gomnd
		if len(nameAndTypeGroups) >= 2 {
			argName = nameAndTypeGroups[1]
		}
		// nolint: gomnd
		if len(nameAndTypeGroups) >= 3 {
			rawArgType = strings.Trim(nameAndTypeGroups[2], "() ")
		}
	}

	if rawArgType == "" {
		rawArgType = "json" // we default to JSON, which is the most generic
	}
	argType := parseType(rawArgType)
	return argName, argType
}

func parseType(rawType string) *StarlarkArgumentType {
	trimmedType := strings.Trim(rawType, " ()")

	parsedTypeFromEnumMaybe, err := StarlarkValueTypeString(trimmedType)
	if err == nil {
		return &StarlarkArgumentType{
			Type:       parsedTypeFromEnumMaybe,
			InnerType1: nil,
			InnerType2: nil,
		}
		// it's not a string, int or json. It might be a composite type (i.e. dict)
	}

	if dictTypeRegexp.MatchString(trimmedType) {
		nestedTypeGroups := dictTypeRegexp.FindStringSubmatch(trimmedType)
		// nolint: gomnd
		if len(nestedTypeGroups) >= 3 {
			parsedInnerType1, err := StarlarkValueTypeString(nestedTypeGroups[1])
			if err != nil {
				// unknown type, unable to parse it
				return nil
			}
			parsedInnerType2, err := StarlarkValueTypeString(nestedTypeGroups[2])
			if err != nil {
				// unknown type, unable to parse it
				return nil
			}
			return &StarlarkArgumentType{
				Type:       StarlarkValueType_Dict,
				InnerType1: &parsedInnerType1,
				InnerType2: &parsedInnerType2,
			}
		}
	}

	if listTypeRegexp.MatchString(trimmedType) {
		nestedTypeGroups := listTypeRegexp.FindStringSubmatch(trimmedType)
		// nolint: gomnd
		if len(nestedTypeGroups) >= 2 {
			parsedInnerType1, err := StarlarkValueTypeString(nestedTypeGroups[1])
			if err != nil {
				// unknown type, unable to parse it
				return nil
			}
			return &StarlarkArgumentType{
				Type:       StarlarkValueType_List,
				InnerType1: &parsedInnerType1,
				InnerType2: nil, // no inner type 2 for lists
			}
		}
	}
	return nil
}
