package catalog

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/stacktrace"
	"net/url"
	"path"
	"strings"
)

const (
	githubDomain              = "github.com"
	httpsSchema               = "https"
	httpSchemaDomainSeparator = "://"
	httpPathSeparator         = "/"
	rootPathSeparator         = "/"

	minElementsInPackageNme = 2
)

type packageDataFromCatalog struct {
	// packageName is the Kurtosis package name
	packageName types.PackageName

	// repositoryOwner is the repositoryOwner of the GitHub repository. It can be a GitHub organization or an individual user
	repositoryOwner string

	// repositoryName is the name of the GitHub repository
	repositoryName string

	// repositoryPackageRootPath is the path relative to the GitHub repository root where the kurtosis.yml file can be found
	repositoryPackageRootPath string
}

func (mainMetadata *packageDataFromCatalog) GetPackageName() types.PackageName {
	return mainMetadata.packageName
}

func (mainMetadata *packageDataFromCatalog) GetRepositoryOwner() string {
	return mainMetadata.repositoryOwner
}

func (mainMetadata *packageDataFromCatalog) GetRepositoryName() string {
	return mainMetadata.repositoryName
}

func (mainMetadata *packageDataFromCatalog) GetRepositoryPackageRootPath() string {
	return mainMetadata.repositoryPackageRootPath
}

func createPackageDataFromCatalogFromPackageName(packageName types.PackageName) (*packageDataFromCatalog, error) {
	// we expect something like github.com/repository-repositoryOwner/repository-name
	// we don't want schemas
	parsedURL, err := url.Parse(string(packageName))
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred parsing the URL of package '%s'", packageName)
	}
	if parsedURL.Scheme != "" {
		return nil, stacktrace.NewError("an error occurred parsing the URL of module '%v'. Expected schema to be empty got '%v'", packageName, parsedURL.Scheme)
	}

	// we prefix schema and make sure that the URL still parses
	packageURLPrefixedWithHttps := fmt.Sprintf("%s%s%s", httpsSchema, httpSchemaDomainSeparator, packageName)
	parsedURL, err = url.Parse(packageURLPrefixedWithHttps)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred parsing the URL with scheme '%s' for package '%s'", packageURLPrefixedWithHttps, packageName)
	}
	if strings.ToLower(parsedURL.Host) != githubDomain {
		return nil, stacktrace.NewError("an error occurred parsing the URL of package. We only support packages on Github for now but got '%v'", packageName)
	}

	splitPath := cleanPathAndSplit(parsedURL.Path)

	if len(splitPath) < minElementsInPackageNme {
		return nil, stacktrace.Propagate(err, "expected to find the repository name and owner in package name '%s' but at least one of them is missing, this is a bug in the Kurtosis indexer", packageName)
	}

	repositoryOwner := splitPath[0]
	repositoryName := splitPath[1]
	rootPathArray := splitPath[2:]

	rootPath := rootPathSeparator
	if len(rootPathArray) > 0 {
		rootPath = path.Join(rootPathArray...) + rootPathSeparator
	}

	newPackageRepositoryMainMetadata := &packageDataFromCatalog{
		packageName:               packageName,
		repositoryOwner:           repositoryOwner,
		repositoryName:            repositoryName,
		repositoryPackageRootPath: rootPath,
	}

	return newPackageRepositoryMainMetadata, nil
}

func cleanPathAndSplit(urlPath string) []string {
	cleanPath := path.Clean(urlPath)
	splitPath := strings.Split(cleanPath, httpPathSeparator)
	var sliceWithoutEmptyStrings []string
	for _, subPath := range splitPath {
		if subPath != "" {
			sliceWithoutEmptyStrings = append(sliceWithoutEmptyStrings, subPath)
		}
	}
	return sliceWithoutEmptyStrings
}
