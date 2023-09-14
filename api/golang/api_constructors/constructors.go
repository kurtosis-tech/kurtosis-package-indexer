package api_constructors

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
)

func NewGetPackagesResponse(packages ...*generated.KurtosisPackage) *generated.GetPackagesResponse {
	return &generated.GetPackagesResponse{Packages: packages}
}

func NewKurtosisPackage(name string, description string, url string, stars uint64, entrypointDescription string, returnsDescription string, args ...*generated.PackageArg) *generated.KurtosisPackage {
	var optionalEntrypointDescription *string
	if entrypointDescription != "" {
		optionalEntrypointDescription = &entrypointDescription
	}
	var optionalReturnsDescription *string
	if returnsDescription != "" {
		optionalReturnsDescription = &returnsDescription
	}
	return &generated.KurtosisPackage{
		Name:                  name,
		Description:           description,
		Url:                   url,
		Args:                  args,
		Stars:                 stars,
		EntrypointDescription: optionalEntrypointDescription,
		ReturnsDescription:    optionalReturnsDescription,
	}
}

func NewPackageArg(name string, description string, isRequired bool, argType *generated.ArgumentValueType, argTypeV2 *generated.PackageArgumentType) *generated.PackageArg {
	var optionalDescription *string
	if description != "" {
		optionalDescription = &description
	}
	return &generated.PackageArg{
		Name:        name,
		Description: optionalDescription,
		IsRequired:  isRequired,
		Type:        argType,
		TypeV2:      argTypeV2,
	}
}
