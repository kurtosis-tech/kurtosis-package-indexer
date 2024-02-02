package main

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated/generatedconnect"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/crawler"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/logger"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/metrics"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/resource"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	connect_server "github.com/kurtosis-tech/kurtosis/connect-server"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	kurtosisPackageIndexerPort = 9770

	successExitCode = 0
	failureExitCode = 1

	grpcServerStopGracePeriod = 5 * time.Second
)

func main() {
	ctx := context.Background()
	if err := logger.ConfigureLogger(); err != nil {
		exitFailure(err)
	}

	// Set up the store which will store all the packages. For now all in memory
	// I left a comment to remember that there were implementations for etcd and boltdb but, both were removed
	// because they were not being used in production, simply to remember that both can be added again by checking the commit history
	indexerStore := store.NewInMemoryStore()

	indexerCtx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	// Set up the metrics reporter which will query the metrics storage on a periodical basis
	metricsReporter, err := metrics.CreateAndScheduleReporter(indexerCtx, indexerStore)
	if err != nil {
		exitFailure(stacktrace.Propagate(err, "an error occurred creating and scheduling the metrics reporter while bootstrapping the server "+
			"Check if the required metrics storage env vars are set, this is the most probably failure"))
	}

	// Set up the crawler which will populate the store on a periodical basis
	indexerCrawler, err := crawler.NewGitHubCrawler(indexerCtx, indexerStore)
	if err != nil {
		exitFailure(stacktrace.Propagate(err, "an error occurred creating the GitHubCrawler while bootstrapping the server"))
	}
	if err := indexerCrawler.Schedule(false); err != nil {
		exitFailure(stacktrace.Propagate(err, "an error occurred scheduling the GitHubCrawler while bootstrapping the server"))
	}

	if err := runServer(indexerStore, indexerCrawler, metricsReporter); err != nil {
		exitFailure(err)
	}
	logrus.Exit(successExitCode)
}

func runServer(indexerStore store.KurtosisIndexerStore, indexerCrawler *crawler.GithubCrawler, metricsReporter metrics.Reporter) error {
	kurtosisPackageIndexerResource := resource.NewKurtosisPackageIndexer(indexerStore, indexerCrawler, metricsReporter)
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

func exitFailure(err error) {
	logrus.Error(err.Error())
	logrus.Exit(failureExitCode)
}
