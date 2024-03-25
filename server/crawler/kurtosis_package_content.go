package crawler

import (
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type KurtosisPackageContent struct {
	// RepositoryMetadata are metadata about the GitHub repository hosting the package
	RepositoryMetadata *PackageRepositoryMetadata

	// Name corresponds to the `name` field in the kurtosis.yml file
	Name string

	// Description corresponds to the `description` field in the kurtosis.yml file
	Description string

	// EntrypointDescription is extracted from the docstring of the `run` function in the main.star file of the package
	// It corresponds to the text preceding the `Args:` block
	EntrypointDescription string

	// ReturnsDescription is extracted from the docstring of the `run` function in the main.star file of the package
	// It corresponds to the content of the `Returns:` block
	ReturnsDescription string

	// PackageArguments is extracted from the docstring of the `run` function in the main.star file of the package
	// It the parsed content of the `Args:` block
	PackageArguments []*StarlarkFunctionArgument

	// The result from parsing the package. This will hold any potential parsing errors
	ParsingResult string

	// The time this package was parsed
	ParsingTime *timestamppb.Timestamp

	// Commit SHA for the package
	Version string

	// The Kurtosis package icon's URL.
	IconURL string

	// RunCount represents the package run count provided by the Kurtosis metrics storage (currently Snowflake)
	RunCount uint32

	// The locator that can be used from within a Starlark script to run this package
	Locator string
}

func NewKurtosisPackageContent(
	repositoryMetadata *PackageRepositoryMetadata,
	name string,
	description string,
	entrypointDescription string,
	returnsDescription string,
	parsingResult string,
	parsingTime *timestamppb.Timestamp,
	version string,
	iconURL string,
	runCount uint32,
	locator string,
	packageArguments ...*StarlarkFunctionArgument,
) *KurtosisPackageContent {
	return &KurtosisPackageContent{
		RepositoryMetadata:    repositoryMetadata,
		Name:                  name,
		Description:           description,
		EntrypointDescription: entrypointDescription,
		ReturnsDescription:    returnsDescription,
		PackageArguments:      packageArguments,
		ParsingResult:         parsingResult,
		ParsingTime:           parsingTime,
		Version:               version,
		IconURL:               iconURL,
		RunCount:              runCount,
		Locator:               locator,
	}
}
