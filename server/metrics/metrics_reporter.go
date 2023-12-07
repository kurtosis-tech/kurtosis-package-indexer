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

	queryIntervalBuffer = 15 * time.Second
)

type PackagesRunCount map[string]uint32

type Reporter struct {
	ctx              context.Context
	snowflake        *snowflake
	packagesRunCount PackagesRunCount // TODO move it into the storage
	ticker           *ticker.Ticker
	store            store.KurtosisIndexerStore
}

func CreateAndScheduleReporter(ctx context.Context, store store.KurtosisIndexerStore) (*Reporter, error) {
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

	if err := newMetricsReporter.schedule(false); err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred scheduling the metrics reporter")
	}

	return newMetricsReporter, nil
}

func (reporter *Reporter) GetPackageRunMetrics() PackagesRunCount {
	return reporter.packagesRunCount
}

func (reporter *Reporter) schedule(forceRunNow bool) error {
	if reporter.ticker != nil {
		logrus.Infof("Reporter already scheduled - stopping it first")
		reporter.ticker.Stop()
	}

	lastMetricsQueryDatetime, err := reporter.store.GetLastCrawlDatetime(reporter.ctx)
	if err != nil {
		return stacktrace.Propagate(err, "An unexpected error occurred retrieving last metrics query datetime from the store")
	}

	var initialDelay time.Duration
	if !forceRunNow && lastMetricsQueryDatetime.Add(queryFrequency).After(time.Now()) {
		initialDelay = time.Until(lastMetricsQueryDatetime.Add(queryFrequency))
	}
	logrus.Infof("Reporter starting with an initial delay of '%v' and a period of '%v'", initialDelay, queryFrequency)
	reporter.ticker = ticker.NewTicker(initialDelay, queryFrequency)
	go func() {
		for {
			select {
			case <-reporter.ctx.Done():
				logrus.Info("Reporter has been closed. Returning")
				reporter.ticker.Stop()
				return
			case tickerTime := <-reporter.ticker.C:
				reporter.doUpdateMetricsNoFailure(reporter.ctx, tickerTime)
			}
		}
	}()
	return nil
}

func (reporter *Reporter) doUpdateMetricsNoFailure(ctx context.Context, tickerTime time.Time) {
	updateMetricsSuccessful := false
	lastMetricsQueryDatetime, err := reporter.store.GetLastMetricsQueryDatetime(ctx)
	if err != nil {
		logrus.Errorf("Could not retrieve last metrics query datetime from the store. This is not critical, query"+
			"will just continue even though the last query might be recent. Error was:\n%s", err.Error())
	}

	// Add a small buffer to avoid false positive when checking if the new query is sooner than the frequency interval
	lastMetricsQueryDatetimeWithBuffer := lastMetricsQueryDatetime.Add(-queryIntervalBuffer)
	if time.Since(lastMetricsQueryDatetimeWithBuffer) < queryFrequency {
		logrus.Infof("Last metrics query happened as '%v' ('%v' ago), which is more recent than the query frequency "+
			"set to '%v', so query will be skipped. If the reporter is running with more than one node, it might be "+
			"that another node did the query in between, and this is totally expected.",
			lastMetricsQueryDatetimeWithBuffer, time.Since(lastMetricsQueryDatetimeWithBuffer), queryFrequency)
		return
	}

	// we persist the metrics query datetime before doing the query so that potential other nodes don't query at the same time
	currentQueryDatetime := time.Now()
	if err = reporter.store.UpdateLastMetricsQueryDatetime(reporter.ctx, currentQueryDatetime); err != nil {
		logrus.Errorf("An error occurred persisting metrics query time to database. This is not critical, but in case of "+
			"a service restart (or in a multiple nodes environment), query might happen more frequently than "+
			"expected. Error was was:\n%v", err.Error())
	} else {
		defer func() {
			if updateMetricsSuccessful {
				return
			}
			logrus.Debugf("Reverting the last query datetime to '%s'. Current value is '%s'",
				lastMetricsQueryDatetime, currentQueryDatetime)
			// revert the query datetime to its previous value
			if err = reporter.store.UpdateLastCrawlDatetime(reporter.ctx, lastMetricsQueryDatetime); err != nil {
				logrus.Errorf("An error occurred reverting the last query datetime to '%s'. Its value"+
					"will remain '%s' and no query will happen before '%s'. Error was:\n%v",
					lastMetricsQueryDatetime, currentQueryDatetime, currentQueryDatetime.Add(queryFrequency), err.Error())
			}
		}()
	}

	logrus.Info("Querying metrics storage for getting Kurtosis package metrics...")
	if err := reporter.updateMetrics(ctx); err != nil {
		logrus.Errorf("An error occurred querying for Kurtosis packages metrics. The last query datetime"+
			"will be reverted to its previous value '%v'. This node will try querying again in '%v'. "+
			"Error was:\n%s", lastMetricsQueryDatetime, queryFrequency, err.Error())
		updateMetricsSuccessful = false
	} else {
		updateMetricsSuccessful = true
	}
	logrus.Infof("... query finished in %v. Success: '%v'. Total metrics rows: %d",
		time.Since(tickerTime), updateMetricsSuccessful, len(reporter.packagesRunCount))
}

func (reporter *Reporter) updateMetrics(ctx context.Context) error {
	runMetricsRows, err := reporter.snowflake.runQueryAndGetPackageRunMetrics(ctx)
	if err != nil {
		return stacktrace.Propagate(err, "an error occurred running the query to get the package run metrics")
	}
	reporter.packagesRunCount = runMetricsRows

	return nil
}
