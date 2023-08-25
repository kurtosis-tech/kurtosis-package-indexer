package crawler

import (
	"github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	"github.com/kurtosis-tech/stacktrace"
	"gopkg.in/yaml.v3"
)

func ParseKurtosisYaml(kurtosisYamlContent *github.RepositoryContent) (store.KurtosisPackageIdentifier, error) {
	rawFileContent, err := kurtosisYamlContent.GetContent()
	if err != nil {
		return "", stacktrace.Propagate(err, "An error occurred getting the content of the '%s' file", kurtosisYamlFileName)
	}

	fileContent := new(KurtosisYaml)
	if err = yaml.Unmarshal([]byte(rawFileContent), fileContent); err != nil {
		return "", stacktrace.Propagate(err, "An error occurred parsing YAML for '%s'", kurtosisYamlFileName)
	}
	return store.KurtosisPackageIdentifier(fileContent.Name), nil
}