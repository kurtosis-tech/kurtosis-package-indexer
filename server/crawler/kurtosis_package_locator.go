package crawler

import "github.com/google/go-github/v54/github"

type KurtosisPackageLocator struct {
	Repository *github.Repository

	RootPath string

	KurtosisYamlFileName string
}

func NewKurtosisPackageLocator(repository *github.Repository, rootPath string, kurtosisYamlFileName string) *KurtosisPackageLocator {
	return &KurtosisPackageLocator{
		Repository:           repository,
		RootPath:             rootPath,
		KurtosisYamlFileName: kurtosisYamlFileName,
	}
}
