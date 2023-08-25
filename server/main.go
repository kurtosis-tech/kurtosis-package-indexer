package main

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated/generatedconnect"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/resource"
	connect_server "github.com/kurtosis-tech/kurtosis/connect-server"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strings"
	"time"
)

const (
	kurtosisPackageIndexerPort = 9770

	successExitCode = 0
	failureExitCode = 1

	grpcServerStopGracePeriod = 5 * time.Second

	forceColors   = true
	fullTimestamp = true

	logMethodAlongWithLogLine = true
	functionPathSeparator     = "."
	emptyFunctionName         = ""
)

func main() {
	ctx := context.Background()
	configureLogger()

	if err := runServer(ctx); err != nil {
		logrus.Errorf("An error occurred when running the Kurtosis package indexer")
		logrus.Exit(failureExitCode)
	}
	logrus.Exit(successExitCode)
}

func runServer(ctx context.Context) error {
	kurtosisPackageIndexerResource := resource.NewKurtosisPackageIndexer()
	connectGoHandler := resource.NewKurtosisPackageIndexerHandlerImpl(kurtosisPackageIndexerResource)

	apiPath, handler := generatedconnect.NewKurtosisPackageIndexerHandler(connectGoHandler)

	apiServer := connect_server.NewConnectServer(
		kurtosisPackageIndexerPort,
		grpcServerStopGracePeriod,
		handler,
		apiPath,
	)

	logrus.Infof("Kurtosis Package Indexer running and listening on port %d", kurtosisPackageIndexerPort)
	if err := apiServer.RunServerUntilInterruptedWithCors(cors.AllowAll()); err != nil {
		logrus.Error("An error occurred running the server", err)
	}
	return nil
}

func configureLogger() {
	logrus.SetLevel(logrus.DebugLevel)
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
