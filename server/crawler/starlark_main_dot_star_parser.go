package crawler

import (
	"github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"go.starlark.net/syntax"
	"reflect"
	"regexp"
)

const (
	mainFunctionName = "run"
)

func ParseStarlarkMainDoStar(kurtosisYamlContent *github.RepositoryContent) (*KurtosisMainDotStar, error) {
	rawFileContent, err := kurtosisYamlContent.GetContent()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the content of the '%s' file", kurtosisYamlFileName)
	}

	parsedStarlarkFile, err := syntax.LegacyFileOptions().Parse("", rawFileContent, syntax.RetainComments)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred parsing the Starlark file")
	}

	mainFunctionObj, err := extractMainFunction(parsedStarlarkFile)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred extracting function '%s' from Starlark file", mainFunctionName)
	}

	mainFunctionArguments, err := extractFunctionArguments(mainFunctionObj)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred analysing the arguments of the main function from the Starlark file")
	}

	return &KurtosisMainDotStar{
		Arguments: mainFunctionArguments,
	}, nil
}

func extractMainFunction(parsedFile *syntax.File) (*syntax.DefStmt, error) {
	for _, rawStmt := range parsedFile.Stmts {
		defStmt, ok := rawStmt.(*syntax.DefStmt)
		if !ok {
			continue
		}
		if defStmt.Name.Name == mainFunctionName {
			return defStmt, nil
		}
	}
	return nil, stacktrace.NewError("No main statement found in the Starlark file")
}

func extractFunctionArguments(mainFunction *syntax.DefStmt) ([]*StarlarkFunctionArgument, error) {
	var allFunctionArguments []*StarlarkFunctionArgument
	for _, rawArg := range mainFunction.Params {
		extractedInfoFromComment, _ := parseTypeFromCommentsIfPossible(rawArg.Comments())
		starlarkFunctionArgument, err := parseExpr(rawArg, extractedInfoFromComment)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Unable to parse Starlark function argument: %v", rawArg)
		}
		allFunctionArguments = append(allFunctionArguments, starlarkFunctionArgument)
	}
	return allFunctionArguments, nil
}

func parseExpr(rawArg syntax.Expr, extractedInfoFromComment string) (*StarlarkFunctionArgument, error) {
	switch typedArg := rawArg.(type) {
	case *syntax.Ident:
		return &StarlarkFunctionArgument{
			Name:       typedArg.Name,
			Type:       extractedInfoFromComment,
			IsRequired: true,
		}, nil
	case *syntax.BinaryExpr:
		parsedArgument, err := parseExpr(typedArg.X, extractedInfoFromComment)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Unable to parse Starlark function argument: %v", typedArg)
		}
		return &StarlarkFunctionArgument{
			Name:       parsedArgument.Name,
			Type:       extractedInfoFromComment,
			IsRequired: false,
		}, nil
	default:
		return nil, stacktrace.NewError("Type of function parameter no handled: %v", reflect.TypeOf(rawArg))
	}
}

func parseTypeFromCommentsIfPossible(comments *syntax.Comments) (string, bool) {
	if comments == nil {
		return "", false
	}
	// For now only parse the comment that is next to the argument. If we want to we can also parse the other ones
	for _, comment := range comments.Suffix {
		if commentType, ok := parseTypeFromCommentIfPossible(comment.Text); ok {
			return commentType, true
		}
	}
	return "", false
}

func parseTypeFromCommentIfPossible(comment string) (string, bool) {
	rp := regexp.MustCompile("#\\s*type\\s*:\\s*([a-zA-Z]*)\\s*")
	if !rp.MatchString(comment) {
		logrus.Infof("Comment '%s' does not match the type regexp. Type cannot be inferred for this argument", comment)
		return "", false
	}
	match := rp.FindStringSubmatch(comment)
	if len(match) != 2 {
		logrus.Infof("Comment '%s' cannot be parsed as the match were: %v. Type cannot be inferred for this argument", comment, match)
		return "", false
	}
	return match[1], true
}
