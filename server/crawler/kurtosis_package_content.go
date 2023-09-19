package crawler

type KurtosisPackageContent struct {
	// RepositoryMetadata are metadata about the Github repository hosting the package
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
}

func NewKurtosisPackageContent(repositoryMetadata *PackageRepositoryMetadata, name string, description string, entrypointDescription string, returnsDescription string, packageArguments ...*StarlarkFunctionArgument) *KurtosisPackageContent {
	return &KurtosisPackageContent{
		RepositoryMetadata:    repositoryMetadata,
		Name:                  name,
		Description:           description,
		EntrypointDescription: entrypointDescription,
		ReturnsDescription:    returnsDescription,
		PackageArguments:      packageArguments,
	}
}
