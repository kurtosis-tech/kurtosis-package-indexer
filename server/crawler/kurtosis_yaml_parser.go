package crawler

import (
	"github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/consts"
	"github.com/kurtosis-tech/stacktrace"
	"gopkg.in/yaml.v3"
)

func ParseKurtosisYaml(kurtosisYamlContent *github.RepositoryContent) (string, string, string, error) {
	rawFileContent, err := kurtosisYamlContent.GetContent()
	if err != nil {
		return "", "", "", stacktrace.Propagate(err, "An error occurred getting the content of the '%s' file", consts.DefaultKurtosisYamlFilename)
	}

	fileContent := new(KurtosisYaml)
	if err = yaml.Unmarshal([]byte(rawFileContent), fileContent); err != nil {
		return "", "", "", stacktrace.Propagate(err, "An error occurred parsing YAML for '%s'", consts.DefaultKurtosisYamlFilename)
	}

	if fileContent.Name == "" {
		return "", "", "", stacktrace.NewError("Kurtosis YAML file had an empty name. This is invalid.")
	}
	return fileContent.Name, fileContent.Description, kurtosisYamlContent.GetSHA(), nil
}
