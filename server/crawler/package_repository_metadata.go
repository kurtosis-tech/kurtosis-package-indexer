package crawler

import "fmt"

type PackageRepositoryMetadata struct {
	// Owner is the owner of the Github repository. It can be a Github organization or an individual user
	Owner string

	// Name is the name of the Github repository
	Name string

	// RootPath is the path relative to the Github repository root where the kurtosis.yml file can be found
	RootPath string

	// KurtosisYamlFileName is the name of the `kurtosis.yml` file. Right now it is hardcoded to `kurtosis.yml`, but
	// this might be useful if at some point we have packages using a variant such as `kurtosis.yaml`
	KurtosisYamlFileName string

	// SupportedDockerComposeYamlFileNames is a list of supported docker compose yaml file names.
	SupportedDockerComposeYamlFileNames []string

	// Stars is the number of stars of the repository hosting the Kurtosis package
	Stars uint64
}

func NewPackageRepositoryMetadata(
	owner string,
	name string,
	rootPath string,
	kurtosisYamlFileName string,
	supportedDockerComposeYamlFileNames []string,
	stars uint64,
) *PackageRepositoryMetadata {
	return &PackageRepositoryMetadata{
		Owner:                               owner,
		Name:                                name,
		RootPath:                            rootPath,
		KurtosisYamlFileName:                kurtosisYamlFileName,
		SupportedDockerComposeYamlFileNames: supportedDockerComposeYamlFileNames,
		Stars:                               stars,
	}
}

func (metadata *PackageRepositoryMetadata) GetLocator() string {
	return fmt.Sprintf("%s/%s/%s", metadata.Owner, metadata.Name, metadata.RootPath)
}
