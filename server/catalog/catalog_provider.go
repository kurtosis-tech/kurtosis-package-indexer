package catalog

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/stacktrace"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
)

type packagesCatalogFileContent struct {
	Packages []map[string]string
}

func GetPackagesCatalog() (PackageCatalog, error) {
	fileContent, err := getPackagesCatalogFileContent()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the packages catalog file content")
	}
	newPackagesCatalogFileContent := &packagesCatalogFileContent{}
	if err = yaml.Unmarshal(fileContent, newPackagesCatalogFileContent); err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred unmarshalling the package catalog file content")
	}
	newPackagesCatalog, err := newPackageCatalogFromPackageCatalogFileContent(newPackagesCatalogFileContent)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred creating a new package catalog object from this catalog file content '%+v'", newPackagesCatalogFileContent)
	}
	return newPackagesCatalog, nil
}

func getPackagesCatalogFileContent() ([]byte, error) {
	//TODO put a constant and unify it with what we put inside the crawler
	//TODO point the URL to the main branch
	yamlFileURL := "https://raw.githubusercontent.com/kurtosis-tech/kurtosis-package-indexer/lporoli/refactor-get-packages/kurtosis-package-catalog.yml"
	response, getErr := http.Get(yamlFileURL)
	if getErr != nil {
		return nil, stacktrace.Propagate(getErr, "an error occurred getting the yaml file content from URL '%s'", yamlFileURL)
	}
	defer response.Body.Close()
	responseBodyBytes, readAllErr := io.ReadAll(response.Body)
	if readAllErr != nil {
		return nil, stacktrace.Propagate(readAllErr, "an error occurred reading the yaml file content")
	}
	return responseBodyBytes, nil
}

func newPackageCatalogFromPackageCatalogFileContent(catalogFileContent *packagesCatalogFileContent) (PackageCatalog, error) {
	packagesData := []packageDataFromCatalog{}

	for _, packageMap := range catalogFileContent.Packages {
		// TODO put a constant
		packageNameStr, found := packageMap["name"]
		if !found {
			// TODO add constants for the key and for the file name
			return nil, stacktrace.NewError("expected to find key '%s' in the package catalog file content, but it was not found. The kurtosis-package-catalog.yml file could be corrupted", "name")
		}
		packageName := types.PackageName(packageNameStr)
		packageDataFromCatalogObj, err := createPackageDataFromCatalogFromPackageName(packageName)
		if err != nil {
			return nil, stacktrace.Propagate(err, "an error occurred creating package catalog data from package name '%s'", packageName)
		}

		packagesData = append(packagesData, *packageDataFromCatalogObj)
	}

	return packagesData, nil
}
