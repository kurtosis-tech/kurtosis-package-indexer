package crawler

import (
	"github.com/google/go-github/v54/github"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/consts"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"go.starlark.net/syntax"
	"reflect"
	"regexp"
	"strings"
)

const (
	mainFunctionName   = "run"
	starlarkTrueValue  = "True"
	starlarkFalseValue = "False"
)

var (
	argTypeInCommentRegexp         = regexp.MustCompile(`#\s*type\s*:\s*([a-zA-Z]*)\s*`)
	argTypeInCommentRegexpMatchNum = 2
)

func ParseStarlarkMainDotStar(kurtosisYamlContent *github.RepositoryContent) (*KurtosisMainDotStar, error) {
	rawFileContent, err := kurtosisYamlContent.GetContent()
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred getting the content of the '%s' file", consts.DefaultKurtosisYamlFilename)
	}

	return ParseMainDotStarContent(rawFileContent)
}

func ParseMainDotStarContent(rawFileContent string) (*KurtosisMainDotStar, error) {
	parsedStarlarkFile, err := syntax.Parse("", rawFileContent, syntax.RetainComments)
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

	parsedMainFunctionDocstring, err := extractAndParseDocstring(mainFunctionObj)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred extracting the docstring of the run function")
	}

	kurtosisMainDotStar := reconcileRunFunctionArgumentWithDocstring(mainFunctionArguments, parsedMainFunctionDocstring)
	return kurtosisMainDotStar, nil
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
		extractedTypeFromComment, _ := parseTypeFromCommentsIfPossible(rawArg.Comments())
		starlarkFunctionArgument, err := parseExpr(rawArg, extractedTypeFromComment)
		if err != nil {
			return nil, stacktrace.Propagate(err, "Unable to parse Starlark function argument: %v", rawArg)
		}
		allFunctionArguments = append(allFunctionArguments, starlarkFunctionArgument)
	}
	return allFunctionArguments, nil
}

func parseExpr(rawArg syntax.Expr, extractedTypeFromComment *StarlarkArgumentType) (*StarlarkFunctionArgument, error) {
	switch typedArg := rawArg.(type) {
	case *syntax.Ident:
		return &StarlarkFunctionArgument{
			Name:         typedArg.Name,
			Description:  "",
			Type:         extractedTypeFromComment,
			IsRequired:   true,
			DefaultValue: nil,
		}, nil
	case *syntax.BinaryExpr:
		parsedArgument, err := parseExpr(typedArg.X, extractedTypeFromComment) // X, left side of binary expr, is the arg name
		if err != nil {
			return nil, stacktrace.Propagate(err, "Unable to parse Starlark function argument: %v", typedArg)
		}
		parsedDefaultValue := parseDefaultValue(typedArg.Y) // Y, right side of binary expr, is the default value
		return &StarlarkFunctionArgument{
			Name:         parsedArgument.Name,
			Description:  "",
			Type:         extractedTypeFromComment,
			IsRequired:   false,
			DefaultValue: parsedDefaultValue,
		}, nil
	default:
		return nil, stacktrace.NewError("Type of function parameter not handled: %v", reflect.TypeOf(rawArg))
	}
}

// Returns pointer to raw string of the default value if [rawDefaultValue] is a bool, int, or string
// Returns nil otherwise
func parseDefaultValue(rawDefaultValue syntax.Expr) *string {
	var parsedDefaultValue string
	switch typedDefaultValue := rawDefaultValue.(type) {
	case *syntax.Ident: // True and False are Ident
		if typedDefaultValue.Name == starlarkTrueValue || typedDefaultValue.Name == starlarkFalseValue {
			parsedDefaultValue = strings.ReplaceAll(typedDefaultValue.Name, `"`, "") // strip quotes
			return &parsedDefaultValue
		} else {
			return nil
		}
	case *syntax.Literal: // string and int are Literals
		parsedDefaultValue = strings.ReplaceAll(typedDefaultValue.Raw, `"`, "") // strip quotes
		return &parsedDefaultValue
	default: // rawDefaultValue is not a primitive or supported yet
		// TODO: add support for dict and list default values
		logrus.Warnf("Attempted to parse a default value that's not supported. Returning no default value.")
		return nil
	}
}

func parseTypeFromCommentsIfPossible(comments *syntax.Comments) (*StarlarkArgumentType, bool) {
	if comments == nil {
		return nil, false
	}
	// For now only parse the comment that is next to the argument. If we want to we can also parse the other ones
	for _, comment := range comments.Suffix {
		if commentType, ok := parseTypeFromCommentIfPossible(comment.Text); ok {
			return commentType, true
		}
	}
	return nil, false
}

func parseTypeFromCommentIfPossible(comment string) (*StarlarkArgumentType, bool) {
	if !argTypeInCommentRegexp.MatchString(comment) {
		logrus.Infof("Comment '%s' does not match the type regexp. Type cannot be inferred for this argument", comment)
		return nil, false
	}
	matches := argTypeInCommentRegexp.FindStringSubmatch(comment)
	if len(matches) != argTypeInCommentRegexpMatchNum {
		logrus.Infof("Comment '%s' cannot be parsed as the match were: %v. Type cannot be inferred for this argument",
			comment, matches)
		return nil, false
	}
	rawTypeFromComment := matches[1]
	parsedTypeFromComment := parseType(rawTypeFromComment)
	return parsedTypeFromComment, true
}

func extractAndParseDocstring(mainFunction *syntax.DefStmt) (*KurtosisMainDotStar, error) {
	docstringContent := extractRawDocstring(mainFunction)

	kurtosisMainDotStar, err := ParseRunFunctionDocstring(docstringContent)
	if err != nil {
		return nil, stacktrace.Propagate(err, "An error occurred parsing the run function docstring comment")
	}
	return kurtosisMainDotStar, nil
}

func extractRawDocstring(mainFunction *syntax.DefStmt) string {
	var firstStmt syntax.Stmt
	if len(mainFunction.Body) == 0 {
		return ""
	}
	firstStmt = mainFunction.Body[0]
	firstStmtExpr, ok := firstStmt.(*syntax.ExprStmt)
	if !ok {
		return ""
	}
	firstStmtLiteral, ok := firstStmtExpr.X.(*syntax.Literal)
	if !ok {
		return ""
	}
	if firstStmtLiteral.Token != syntax.STRING {
		return ""
	}
	firstStmtContent, ok := firstStmtLiteral.Value.(string)
	if !ok {
		return ""
	}
	return firstStmtContent
}

func reconcileRunFunctionArgumentWithDocstring(runFunctionArguments []*StarlarkFunctionArgument, parsedDocstringContent *KurtosisMainDotStar) *KurtosisMainDotStar {
	finalKurtosisMainDotStar := &KurtosisMainDotStar{
		Description:       parsedDocstringContent.Description,
		Arguments:         nil,
		ReturnDescription: parsedDocstringContent.ReturnDescription,
	}

	indexedArgumentsFromDocstring := map[string]*StarlarkFunctionArgument{}
	for _, argument := range parsedDocstringContent.Arguments {
		indexedArgumentsFromDocstring[argument.Name] = argument
	}

	var packageArguments []*StarlarkFunctionArgument
	for _, argument := range runFunctionArguments {
		assembledArgument := &StarlarkFunctionArgument{
			Name:         argument.Name,
			Description:  "",
			Type:         argument.Type,
			IsRequired:   argument.IsRequired,
			DefaultValue: argument.DefaultValue,
		}
		if argumentFromDocstring, ok := indexedArgumentsFromDocstring[argument.Name]; ok {
			assembledArgument.Description = argumentFromDocstring.Description
			assembledArgument.Type = argumentFromDocstring.Type
		}
		packageArguments = append(packageArguments, assembledArgument)
	}
	finalKurtosisMainDotStar.Arguments = packageArguments
	return finalKurtosisMainDotStar
}
