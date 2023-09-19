package api_constructors

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
)

func NewGetPackagesResponse(packages ...*generated.KurtosisPackage) *generated.GetPackagesResponse {
	return &generated.GetPackagesResponse{Packages: packages}
}

func NewKurtosisPackage(name string, description string, url string, stars uint64, entrypointDescription string, returnsDescription string, args ...*generated.PackageArg) *generated.KurtosisPackage {
	return &generated.KurtosisPackage{
		Name:                  name,
		Description:           description,
		Url:                   url,
		Args:                  args,
		Stars:                 stars,
		EntrypointDescription: entrypointDescription,
		ReturnsDescription:    returnsDescription,
	}
}

func NewPackageArg(name string, description string, isRequired bool, argType *generated.PackageArgumentType) *generated.PackageArg {
	return &generated.PackageArg{
		Name:        name,
		Description: description,
		IsRequired:  isRequired,
		TypeV2:      argType,
	}
}
