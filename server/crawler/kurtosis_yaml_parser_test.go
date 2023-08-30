package crawler

import (
	"fmt"
	"github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDecodeKurtosisYaml_Minimal(t *testing.T) {
	packageName := store.KurtosisPackageIdentifier("github.com/kurtosis-tech/nginx-package")
	contentStr := fmt.Sprintf(`name: %q
`, packageName)
	// nolint: exhaustruct
	content := &github.RepositoryContent{
		Content: &contentStr,
	}
	parsedPackageName, parsedPackageDescription, err := ParseKurtosisYaml(content)
	require.NoError(t, err)
	require.Equal(t, packageName, parsedPackageName)
	require.Equal(t, "", parsedPackageDescription)
}

func TestDecodeKurtosisYaml_Full(t *testing.T) {
	packageName := store.KurtosisPackageIdentifier("github.com/kurtosis-tech/nginx-package")
	packageDescription := "Cool package"
	contentStr := fmt.Sprintf(`name: %q
description: %q
`, packageName, packageDescription)
	// nolint: exhaustruct
	content := &github.RepositoryContent{
		Content: &contentStr,
	}
	parsedPackageName, parsedPackageDescription, err := ParseKurtosisYaml(content)
	require.NoError(t, err)
	require.Equal(t, packageName, parsedPackageName)
	require.Equal(t, packageDescription, parsedPackageDescription)
}

func TestDecodeKurtosisYaml_UnknownField(t *testing.T) {
	packageName := store.KurtosisPackageIdentifier("github.com/kurtosis-tech/nginx-package")
	contentStr := fmt.Sprintf(`name: %q
unknownField: hello world
`, packageName)
	// nolint: exhaustruct
	content := &github.RepositoryContent{
		Content: &contentStr,
	}
	parsedPackageName, parsedPackageDescription, err := ParseKurtosisYaml(content)
	require.NoError(t, err)
	require.Equal(t, packageName, parsedPackageName)
	require.Equal(t, "", parsedPackageDescription)
}

func TestDecodeKurtosisYaml_WithComments(t *testing.T) {
	packageName := store.KurtosisPackageIdentifier("github.com/kurtosis-tech/nginx-package")
	contentStr := fmt.Sprintf(`# TODO replace this with your Github username and this repo's name!
name: %q
`, packageName)
	// nolint: exhaustruct
	content := &github.RepositoryContent{
		Content: &contentStr,
	}
	parsedPackageName, parsedPackageDescription, err := ParseKurtosisYaml(content)
	require.NoError(t, err)
	require.Equal(t, packageName, parsedPackageName)
	require.Equal(t, "", parsedPackageDescription)
}
