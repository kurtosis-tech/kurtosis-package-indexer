package api_constructors

import (
	"fmt"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func NewGetPackagesResponse(packages ...*generated.KurtosisPackage) *generated.GetPackagesResponse {
	return &generated.GetPackagesResponse{Packages: packages}
}

func NewKurtosisPackage(
	name string,
	description string,
	repository *generated.PackageRepository,
	stars uint64,
	entrypointDescription string,
	returnsDescription string,
	parsingResult string,
	parsingTime *timestamppb.Timestamp,
	version string,
	iconURL string,
	args ...*generated.PackageArg,
) *generated.KurtosisPackage {
	// construct the URL from the repository object for now. Remove it if it's not needed by the FE
	url := fmt.Sprintf("%s/%s/%s/blob/main/%s", repository.BaseUrl, repository.Owner, repository.Name, repository.RootPath)

	return &generated.KurtosisPackage{
		Name:                  name,
		Description:           description,
		RepositoryMetadata:    repository,
		Url:                   &url,
		Args:                  args,
		Stars:                 stars,
		EntrypointDescription: entrypointDescription,
		ReturnsDescription:    returnsDescription,
		ParsingResult:         parsingResult,
		ParsingTime:           parsingTime,
		Version:               version,
		IconUrl:               iconURL,
	}
}

func NewPackageArg(name string, description string, isRequired bool, argType *generated.PackageArgumentType, defaultValue *string) *generated.PackageArg {
	return &generated.PackageArg{
		Name:         name,
		Description:  description,
		IsRequired:   isRequired,
		TypeV2:       argType,
		DefaultValue: defaultValue,
	}
}

func NewPackageRepository(
	baseUrl string,
	owner string,
	name string,
	rootPath string,
	lastCommitTime time.Time,
	defaultBranch string,
) *generated.PackageRepository {
	return &generated.PackageRepository{
		BaseUrl:        baseUrl,
		Owner:          owner,
		Name:           name,
		RootPath:       rootPath,
		LastCommitTime: timestamppb.New(lastCommitTime),
		DefaultBranch:  defaultBranch,
	}
}

func NewReadPackageResponse(
	pack *generated.KurtosisPackage,
) *generated.ReadPackageResponse {
	return &generated.ReadPackageResponse{Package: pack}
}
