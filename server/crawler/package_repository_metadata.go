package crawler

import (
	"fmt"
	"github.com/kurtosis-tech/stacktrace"
	"net/url"
	"time"
)

const gitHuhDownloadBaseURL = "https://raw.githubusercontent.com"

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

	// Stars is the number of stars of the repository hosting the Kurtosis package
	Stars uint64

	// LastCommitTime the time of the last commit on the main branch
	LastCommitTime time.Time

	// DefaultBranch is the repository default branch, it's the branch used during `kurtosis run {package}`
	DefaultBranch string
}

func NewPackageRepositoryMetadata(
	owner string,
	name string,
	rootPath string,
	kurtosisYamlFileName string,
	stars uint64,
	lastCommitTime time.Time,
	defaultBranch string,
) *PackageRepositoryMetadata {
	return &PackageRepositoryMetadata{
		Owner:                owner,
		Name:                 name,
		RootPath:             rootPath,
		KurtosisYamlFileName: kurtosisYamlFileName,
		Stars:                stars,
		LastCommitTime:       lastCommitTime,
		DefaultBranch:        defaultBranch,
	}
}

func (metadata *PackageRepositoryMetadata) GetLocator() string {
	return fmt.Sprintf("%s/%s/%s", metadata.Owner, metadata.Name, metadata.RootPath)
}

func (metadata *PackageRepositoryMetadata) GetDownloadRootURL() (string, error) {
	rootURL, err := url.JoinPath(gitHuhDownloadBaseURL, metadata.Owner, metadata.Name, metadata.DefaultBranch, metadata.RootPath)
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred creating the repository download root path for '%+v'", metadata)
	}
	return rootURL, nil
}
