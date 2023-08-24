package main

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/resource"
	minimal_grpc_server "github.com/kurtosis-tech/minimal-grpc-server/golang/server"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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

	if err := runServer(ctx); err != nil {
		logrus.Errorf("An error occurred when running the Kurtosis package indexer")
		logrus.Exit(failureExitCode)
	}
	logrus.Exit(successExitCode)
}

func runServer(ctx context.Context) error {
	kurtosisPackageIndexerResource := resource.NewKurtosisPackageIndexer()

	kurtosisPackageIndexerRegistrationFunc := func(grpcServer *grpc.Server) {
		generated.RegisterKurtosisPackageIndexerServer(grpcServer, kurtosisPackageIndexerResource)
	}

	kurtosisPackageIndexerServer := minimal_grpc_server.NewMinimalGRPCServer(
		kurtosisPackageIndexerPort,
		grpcServerStopGracePeriod,
		[]func(*grpc.Server){
			kurtosisPackageIndexerRegistrationFunc,
		},
	)

	logrus.Infof("Kurtosis Package Indexer running and listening on port %d", kurtosisPackageIndexerPort)
	if err := kurtosisPackageIndexerServer.RunUntilStopped(ctx.Done()); err != nil {
		return stacktrace.Propagate(err, "An error occurred running the Kurtosis Package Indexer")
	}
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
