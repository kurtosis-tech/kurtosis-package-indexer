package crawler

import (
	"github.com/google/go-github/v54/github"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	stringValueType = StarlarkValueType_String
	intValueType    = StarlarkValueType_Int

	stringType = StarlarkArgumentType{
		Type:       stringValueType,
		InnerType1: nil,
		InnerType2: nil,
	}

	intType = StarlarkArgumentType{
		Type:       intValueType,
		InnerType1: nil,
		InnerType2: nil,
	}

	boolType = StarlarkArgumentType{
		Type:       StarlarkValueType_Bool,
		InnerType1: nil,
		InnerType2: nil,
	}

	jsonType = StarlarkArgumentType{
		Type:       StarlarkValueType_Json,
		InnerType1: nil,
		InnerType2: nil,
	}

	dictStringStringType = StarlarkArgumentType{
		Type:       StarlarkValueType_Dict,
		InnerType1: &stringValueType,
		InnerType2: &stringValueType,
	}

	listIntType = StarlarkArgumentType{
		Type:       StarlarkValueType_List,
		InnerType1: &intValueType,
		InnerType2: nil,
	}
)

func TestParseStarlarkMainDoStar_Minimal(t *testing.T) {
	contentStr := `
def run(plan, args):
	plan.print("Hello World")
`
	// nolint: exhaustruct
	content := &github.RepositoryContent{
		Content: &contentStr,
	}

	result, err := ParseStarlarkMainDoStar(content)
	require.NoError(t, err)

	require.Len(t, result.Arguments, 2)

	require.Equal(t, "plan", result.Arguments[0].Name)
	require.True(t, result.Arguments[0].IsRequired)
	require.Equal(t, "", result.Arguments[0].Description)
	require.Nil(t, result.Arguments[0].Type)

	require.Equal(t, "args", result.Arguments[1].Name)
	require.True(t, result.Arguments[1].IsRequired)
	require.Equal(t, "", result.Arguments[1].Description)
	require.Nil(t, result.Arguments[1].Type)
}

func TestParseStarlarkMainDoStar_ArgsIsOptional(t *testing.T) {
	contentStr := `
def run(plan, args={}):
	plan.print("Hello World")
`
	// nolint: exhaustruct
	content := &github.RepositoryContent{
		Content: &contentStr,
	}

	result, err := ParseStarlarkMainDoStar(content)
	require.NoError(t, err)

	require.Len(t, result.Arguments, 2)

	require.Equal(t, "plan", result.Arguments[0].Name)
	require.True(t, result.Arguments[0].IsRequired)
	require.Equal(t, "", result.Arguments[0].Description)
	require.Nil(t, result.Arguments[0].Type)

	require.Equal(t, "args", result.Arguments[1].Name)
	require.False(t, result.Arguments[1].IsRequired)
	require.Equal(t, "", result.Arguments[1].Description)
	require.Nil(t, result.Arguments[1].Type)
}

func TestParseStarlarkMainDoStar_WithType(t *testing.T) {
	contentStr := `
def run(
		plan,
		first_arg,         #type: string
        second_arg,        # type:int   
        third_arg=True,    # type: bool
		fourth_arg={},
		fifth_arg=[],
		untyped_arg={},
	):
	"""This is the run function

	Args:
		first_arg (string): This is the first arg
		second_arg (int):
		third_arg (bool): yup, bool works too!
		fourth_arg (dict[string, string]): this is a dict
		fifth_arg (list[int]): this is a list of integers 
		untyped_arg: no type info here, will default to JSON
	Returns: 
		Returns nothing for now
	"""
	plan.print("Hello World")
`
	// nolint: exhaustruct
	content := &github.RepositoryContent{
		Content: &contentStr,
	}

	result, err := ParseStarlarkMainDoStar(content)
	require.NoError(t, err)

	require.Len(t, result.Arguments, 7)

	require.Equal(t, "plan", result.Arguments[0].Name)
	require.True(t, result.Arguments[0].IsRequired)
	require.Equal(t, "", result.Arguments[0].Description)
	require.Nil(t, result.Arguments[0].Type)

	require.Equal(t, "first_arg", result.Arguments[1].Name)
	require.True(t, result.Arguments[1].IsRequired)
	require.Equal(t, "This is the first arg", result.Arguments[1].Description)
	require.Equal(t, &stringType, result.Arguments[1].Type)

	require.Equal(t, "second_arg", result.Arguments[2].Name)
	require.True(t, result.Arguments[2].IsRequired)
	require.Equal(t, "", result.Arguments[2].Description)
	require.Equal(t, &intType, result.Arguments[2].Type)

	require.Equal(t, "third_arg", result.Arguments[3].Name)
	require.False(t, result.Arguments[3].IsRequired)
	require.Equal(t, "yup, bool works too!", result.Arguments[3].Description)
	require.Equal(t, &boolType, result.Arguments[3].Type)

	require.Equal(t, "fourth_arg", result.Arguments[4].Name)
	require.False(t, result.Arguments[4].IsRequired)
	require.Equal(t, "this is a dict", result.Arguments[4].Description)
	require.Equal(t, &dictStringStringType, result.Arguments[4].Type)

	require.Equal(t, "fifth_arg", result.Arguments[5].Name)
	require.False(t, result.Arguments[5].IsRequired)
	require.Equal(t, "this is a list of integers", result.Arguments[5].Description)
	require.Equal(t, &listIntType, result.Arguments[5].Type)

	require.Equal(t, "untyped_arg", result.Arguments[6].Name)
	require.False(t, result.Arguments[6].IsRequired)
	require.Equal(t, "no type info here, will default to JSON", result.Arguments[6].Description)
	require.Equal(t, &jsonType, result.Arguments[6].Type)
}

func TestParseStarlarkMainDoStar_WithImports(t *testing.T) {
	contentStr := `
mongodb = import_module("github.com/kurtosis-tech/mongodb-package/main.star")

def run(plan, args):
	plan.print("Hello World")
`
	// nolint: exhaustruct
	content := &github.RepositoryContent{
		Content: &contentStr,
	}

	result, err := ParseStarlarkMainDoStar(content)
	require.NoError(t, err)

	require.Len(t, result.Arguments, 2)

	require.Equal(t, "plan", result.Arguments[0].Name)
	require.True(t, result.Arguments[0].IsRequired)
	require.Nil(t, result.Arguments[0].Type)

	require.Equal(t, "args", result.Arguments[1].Name)
	require.True(t, result.Arguments[1].IsRequired)
	require.Nil(t, result.Arguments[1].Type)
}
