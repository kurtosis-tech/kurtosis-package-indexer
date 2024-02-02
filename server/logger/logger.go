package logger

import (
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/utils"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
)

const (
	forceColors   = true
	fullTimestamp = true

	logMethodAlongWithLogLine = true
	functionPathSeparator     = "."
	emptyFunctionName         = ""

	loggerLogLeverEnvVarKey = "LOGGER_LOG_LEVEL"
)

func ConfigureLogger() error {

	if err := SetLevel(); err != nil {
		return stacktrace.Propagate(err, "an error occurred setting the Logrus log level")
	}

	// This allows the filename & function to be reported
	logrus.SetReportCaller(logMethodAlongWithLogLine)
	// NOTE: we'll want to change the ForceColors to false if we ever want structured logging
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:               forceColors,
		DisableColors:             false,
		ForceQuote:                false,
		DisableQuote:              false,
		EnvironmentOverrideColors: false,
		DisableTimestamp:          false,
		FullTimestamp:             fullTimestamp,
		TimestampFormat:           "",
		DisableSorting:            false,
		SortingFunc:               nil,
		DisableLevelTruncation:    false,
		PadLevelText:              false,
		QuoteEmptyFields:          false,
		FieldMap:                  nil,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			fullFunctionPath := strings.Split(f.Function, functionPathSeparator)
			functionName := fullFunctionPath[len(fullFunctionPath)-1]
			_, filename := path.Split(f.File)
			return emptyFunctionName, formatFilenameFunctionForLogs(filename, functionName)
		},
	})
	return nil
}

func SetLevel() error {
	loggerLogLevelEnVarValue, err := utils.GetFromEnvVar(loggerLogLeverEnvVarKey, "Logger log level")
	if err != nil {
		return stacktrace.Propagate(err, "an error occurred while getting the logger log level value from the env var")
	}

	logrusLogLevel, err := logrus.ParseLevel(loggerLogLevelEnVarValue)
	if err != nil {
		return stacktrace.Propagate(err, "an error occurred parsing lol leve string '%s' to a Logrus log level", logrusLogLevel)
	}

	logrus.SetLevel(logrusLogLevel)

	return nil
}

func formatFilenameFunctionForLogs(filename string, functionName string) string {
	var output strings.Builder
	output.WriteString("[")
	output.WriteString(filename)
	output.WriteString(":")
	output.WriteString(functionName)
	output.WriteString("]")
	return output.String()
}
