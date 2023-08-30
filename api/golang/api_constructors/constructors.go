package api_constructors

import "github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"

func NewGetPackagesResponse(packages ...*generated.KurtosisPackage) *generated.GetPackagesResponse {
	return &generated.GetPackagesResponse{Packages: packages}
}

func NewKurtosisPackage(name string, description string, stars uint64, args ...*generated.PackageArg) *generated.KurtosisPackage {
	return &generated.KurtosisPackage{
		Name:        name,
		Description: description,
		Args:        args,
		Stars:       stars,
	}
}

// NewUntypedPackageArg creates an argument with no type.
// Right now, supported types are only STRING, INT, FLOAT and BOOL, so this should be used for more evolved arguments
// like dictionaries or lists. Support for those can be added later if needed
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
