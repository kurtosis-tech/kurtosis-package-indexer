package catalog

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/utils"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"net/url"
)

const (
	defaultKurtosisPackageCatalogYamlFileName           = "kurtosis-package-catalog.yml"
	defaultKurtosisPackageCatalogYamlFileHostAndDirpath = "https://raw.githubusercontent.com/kurtosis-tech/kurtosis-package-catalog/main/"

	kurtosisPackageCatalogYamlFileNameEnvVarKey           = "KURTOSIS_PACKAGE_CATALOG_YAML_FILENAME"
	kurtosisPackageCatalogYamlFileHostAndDirpathEnvVarKey = "KURTOSIS_PACKAGE_CATALOG_YAML_FILE_HOST_AND_DIRPATH"

	packageNameKey = "name"
)

type packagesCatalogFileContent struct {
	Packages []map[string]string
}

func GetPackagesCatalog() (PackageCatalog, error) {
	fileContent, err := getPackagesCatalogFileContent()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the packages catalog file content")
	}

	newPackagesCatalog, err := GetPackageCatalogFromYamlFileContent(fileContent)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the package catalog")
	}

	return newPackagesCatalog, nil
}

// GetPackageCatalogFromYamlFileContent receives the kurtosis-package-catalog.yml file content and returns a PackageCatalog object
// it's public because it's also used by the kurtosis-package-catalog project
func GetPackageCatalogFromYamlFileContent(fileContent []byte) (PackageCatalog, error) {
	newPackagesCatalogFileContent := &packagesCatalogFileContent{
		Packages: []map[string]string{},
	}
	if err := yaml.Unmarshal(fileContent, newPackagesCatalogFileContent); err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred unmarshalling the package catalog file content")
	}
	newPackagesCatalog, err := newPackageCatalogFromPackageCatalogFileContent(newPackagesCatalogFileContent)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred creating a new package catalog object from this catalog file content '%+v'", newPackagesCatalogFileContent)
	}

	return newPackagesCatalog, nil
}

func getPackagesCatalogFileContent() ([]byte, error) {
	kurtosisPackageCatalogYamlFileURL, err := getKurtosisPackageCatalogYamlFileURL()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the Kurtosis package catalog YAML file URL")
	}

	response, getErr := http.Get(kurtosisPackageCatalogYamlFileURL)
	if getErr != nil {
		return nil, stacktrace.Propagate(getErr, "an error occurred getting the yaml file content from URL '%s'", kurtosisPackageCatalogYamlFileURL)
	}
	defer response.Body.Close()
	responseBodyBytes, readAllErr := io.ReadAll(response.Body)
	if readAllErr != nil {
		return nil, stacktrace.Propagate(readAllErr, "an error occurred reading the yaml file content")
	}
	return responseBodyBytes, nil
}

func newPackageCatalogFromPackageCatalogFileContent(catalogFileContent *packagesCatalogFileContent) (PackageCatalog, error) {
	catalogPackages := catalogFileContent.Packages

	packagesData := make([]packageDataFromCatalog, len(catalogPackages))

	for packageIndex, packageMap := range catalogPackages {
		packageNameStr, found := packageMap[packageNameKey]
		if !found {
			return nil, stacktrace.NewError("expected to find key '%s' in the package catalog file content, but it was not found. The %s file could be corrupted", packageNameKey, defaultKurtosisPackageCatalogYamlFileName)
		}
		packageName := types.PackageName(packageNameStr)
		packageDataFromCatalogObj, err := createPackageDataFromCatalogFromPackageName(packageName)
		if err != nil {
			return nil, stacktrace.Propagate(err, "an error occurred creating package catalog data from package name '%s'", packageName)
		}

		packagesData[packageIndex] = *packageDataFromCatalogObj
	}

	return packagesData, nil
}

func getKurtosisPackageCatalogYamlFileURL() (string, error) {
	var (
		filename       string
		hostAndDirpath string
		err            error
	)
	filename, err = getKurtosisPackageCatalogYamlFileNameFromEnvVar()
	if err != nil {
		logrus.Debugf("Kurtosis package catalog YAML filename env var was not set, using default value '%s'", defaultKurtosisPackageCatalogYamlFileName)
		filename = defaultKurtosisPackageCatalogYamlFileName
	}

	hostAndDirpath, err = getKurtosisPackageCatalogYamlFileHostAndDirpathFromEnvVar()
	if err != nil {
		logrus.Debugf("Kurtosis package catalog YAML file host and dirpath env var was not set, using default value '%s'", defaultKurtosisPackageCatalogYamlFileHostAndDirpath)
		hostAndDirpath = defaultKurtosisPackageCatalogYamlFileHostAndDirpath
	}

	kurtosisPackageCatalogYamlFileURL, err := url.JoinPath(hostAndDirpath, filename)
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred while creating the Kurtosis package catalog YAML file URL using '%s' and '%s'", hostAndDirpath, filename)
	}

	return kurtosisPackageCatalogYamlFileURL, nil
}

func getKurtosisPackageCatalogYamlFileNameFromEnvVar() (string, error) {
	return utils.GetFromEnvVar(kurtosisPackageCatalogYamlFileNameEnvVarKey, "Kurtosis package catalog YAML filename")
}

func getKurtosisPackageCatalogYamlFileHostAndDirpathFromEnvVar() (string, error) {
	return utils.GetFromEnvVar(kurtosisPackageCatalogYamlFileHostAndDirpathEnvVarKey, "Kurtosis package catalog YAML file host and dirpath")
}
