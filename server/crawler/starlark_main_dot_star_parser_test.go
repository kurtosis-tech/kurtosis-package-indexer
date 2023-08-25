package crawler

import (
	"github.com/google/go-github/v54/github"
	"github.com/stretchr/testify/require"
	"testing"
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
	require.Equal(t, "", result.Arguments[0].Type)

	require.Equal(t, "args", result.Arguments[1].Name)
	require.True(t, result.Arguments[1].IsRequired)
	require.Equal(t, "", result.Arguments[1].Type)
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
	require.Equal(t, "", result.Arguments[0].Type)

	require.Equal(t, "args", result.Arguments[1].Name)
	require.False(t, result.Arguments[1].IsRequired)
	require.Equal(t, "", result.Arguments[1].Type)
}

func TestParseStarlarkMainDoStar_WithType(t *testing.T) {
	contentStr := `
def run(
		plan,
		first_arg,         #type: string
        second_arg,        # type:int   
        third_arg=True,    # type: bool
		untyped_arg={},
	):
	plan.print("Hello World")
`
	// nolint: exhaustruct
	content := &github.RepositoryContent{
		Content: &contentStr,
	}

	result, err := ParseStarlarkMainDoStar(content)
	require.NoError(t, err)

	require.Len(t, result.Arguments, 5)

	require.Equal(t, "plan", result.Arguments[0].Name)
	require.True(t, result.Arguments[0].IsRequired)
	require.Equal(t, "", result.Arguments[0].Type)

	require.Equal(t, "first_arg", result.Arguments[1].Name)
	require.True(t, result.Arguments[1].IsRequired)
	require.Equal(t, "string", result.Arguments[1].Type)

	require.Equal(t, "second_arg", result.Arguments[2].Name)
	require.True(t, result.Arguments[2].IsRequired)
	require.Equal(t, "int", result.Arguments[2].Type)

	require.Equal(t, "third_arg", result.Arguments[3].Name)
	require.False(t, result.Arguments[3].IsRequired)
	require.Equal(t, "bool", result.Arguments[3].Type)

	require.Equal(t, "untyped_arg", result.Arguments[4].Name)
	require.False(t, result.Arguments[4].IsRequired)
	require.Equal(t, "", result.Arguments[4].Type)
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
	require.Equal(t, "", result.Arguments[0].Type)

	require.Equal(t, "args", result.Arguments[1].Name)
	require.True(t, result.Arguments[1].IsRequired)
	require.Equal(t, "", result.Arguments[1].Type)
}
