package metrics

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/ticker"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	queryFrequency = 1 * time.Hour
)

type PackagesRunCount map[string]uint32

type Reporter struct {
	ctx              context.Context
	snowflake        *snowflake
	packagesRunCount PackagesRunCount // TODO move it into the storage
	ticker           *ticker.Ticker
	store            store.KurtosisIndexerStore
}

func CreateReporter(ctx context.Context, store store.KurtosisIndexerStore) (*Reporter, error) {
	snowflakeObj, err := createSnowflake()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred creating the Snowflake object")
	}

	newMetricsReporter := &Reporter{
		ctx:              ctx,
		snowflake:        snowflakeObj,
		packagesRunCount: nil,
		ticker:           nil,
		store:            store,
	}

	return newMetricsReporter, nil
}

func (reporter *Reporter) Schedule(forceRunNow bool) error {
	if reporter.ticker != nil {
		logrus.Infof("Reporter already scheduled - stopping it first")
		reporter.ticker.Stop()
	}

	reporter.ticker = ticker.NewTicker(0, queryFrequency) // TODO check the initial delay
	go func() {
		for {
			select {
			case <-reporter.ctx.Done():
				logrus.Info("Reporter has been closed. Returning")
				reporter.ticker.Stop()
				return
			case tickerTime := <-reporter.ticker.C:
				logrus.Debugf("Reporter ticker time '%s'", tickerTime.String())
				if err := reporter.upgradePackagesRunMetrics(reporter.ctx); err != nil {
					logrus.Errorf("an error occurred upgrading the package run metrics from the reporter schduller. Error was>\n%s", err.Error())
				}
			}
		}
	}()
	return nil
}

func (reporter *Reporter) upgradePackagesRunMetrics(ctx context.Context) error {
	runMetricsRows, err := reporter.snowflake.runQueryAndGetPackageRunMetrics(ctx)
	if err != nil {
		return stacktrace.Propagate(err, "an error occurred running the query to get the package run metrics")
	}
	reporter.packagesRunCount = runMetricsRows

	return nil
}

func (reporter *Reporter) GetPackageRunMetrics() PackagesRunCount {
	return reporter.packagesRunCount
}
