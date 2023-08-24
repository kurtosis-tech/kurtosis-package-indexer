package api_constructors

import "github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"

func NewGetPackagesResponse(packages ...*generated.KurtosisPackage) *generated.GetPackagesResponse {
	return &generated.GetPackagesResponse{Packages: packages}
}

func NewKurtosisPackage(name string, args ...*generated.PackageArg) *generated.KurtosisPackage {
	return &generated.KurtosisPackage{
		Name: name,
		Args: args,
	}
}

func NewUntypedPackageArg(name string, isRequired bool) *generated.PackageArg {
	return &generated.PackageArg{
		Name:       name,
		IsRequired: isRequired,
		Type:       nil,
	}
}

func NewPackageArg(name string, isRequired bool, argType generated.PackageArgType) *generated.PackageArg {
	return &generated.PackageArg{
		Name:       name,
		IsRequired: isRequired,
		Type:       &argType,
	}
}
