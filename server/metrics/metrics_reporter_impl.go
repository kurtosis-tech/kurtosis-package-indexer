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
	queryFrequency = 30 * time.Minute

	queryIntervalBuffer = 15 * time.Second
)

type reporterImpl struct {
	ctx       context.Context
	snowflake *snowflake
	ticker    *ticker.Ticker
	store     store.KurtosisIndexerStore
}

func (reporter *reporterImpl) Schedule(forceRunNow bool) error {
	if reporter.ticker != nil {
		logrus.Infof("reporterImpl already scheduled - stopping it first")
		reporter.ticker.Stop()
	}

	lastMetricsQueryDatetime, err := reporter.store.GetLastMainCrawlDatetime(reporter.ctx)
	if err != nil {
		return stacktrace.Propagate(err, "An unexpected error occurred retrieving last metrics query datetime from the store")
	}

	var initialDelay time.Duration
	if !forceRunNow && lastMetricsQueryDatetime.Add(queryFrequency).After(time.Now()) {
		initialDelay = time.Until(lastMetricsQueryDatetime.Add(queryFrequency))
	}
	logrus.Infof("reporterImpl starting with an initial delay of '%v' and a period of '%v'", initialDelay, queryFrequency)
	reporter.ticker = ticker.NewTicker(initialDelay, queryFrequency)
	go func() {
		for {
			select {
			case <-reporter.ctx.Done():
				logrus.Info("reporterImpl has been closed. Returning")
				reporter.ticker.Stop()
				return
			case tickerTime := <-reporter.ticker.C:
				reporter.doUpdateMetricsNoFailure(reporter.ctx, tickerTime)
			}
		}
	}()
	return nil
}

func (reporter *reporterImpl) doUpdateMetricsNoFailure(ctx context.Context, tickerTime time.Time) {
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
			if err = reporter.store.UpdateLastMainCrawlDatetime(reporter.ctx, lastMetricsQueryDatetime); err != nil {
				logrus.Errorf("An error occurred reverting the last query datetime to '%s'. Its value"+
					"will remain '%s' and no query will happen before '%s'. Error was:\n%v",
					lastMetricsQueryDatetime, currentQueryDatetime, currentQueryDatetime.Add(queryFrequency), err.Error())
			}
		}()
	}

	logrus.Infof("Querying metrics storage for getting Kurtosis package metrics from '%s' to '%s'...", lastMetricsQueryDatetime, currentQueryDatetime)
	updatedPackagesRunCount, err := reporter.updateMetrics(ctx, lastMetricsQueryDatetime, currentQueryDatetime)
	if err != nil {
		logrus.Errorf("An error occurred querying for Kurtosis packages metrics. The last query datetime"+
			"will be reverted to its previous value '%v'. This node will try querying again in '%v'. "+
			"Error was:\n%s", lastMetricsQueryDatetime, queryFrequency, err.Error())
		updateMetricsSuccessful = false
	} else {
		updateMetricsSuccessful = true
	}
	logrus.Infof("... query finished in %v. Success: '%v'. Total packages updated: %d",
		time.Since(tickerTime), updateMetricsSuccessful, updatedPackagesRunCount)
}

func (reporter *reporterImpl) updateMetrics(ctx context.Context, fromTime time.Time, toTime time.Time) (int, error) {

	newPackagesRunCount, err := reporter.snowflake.getPackageRunMetricsInDateRange(ctx, fromTime, toTime)
	if err != nil {
		return 0, stacktrace.Propagate(err, "an error occurred running the query to get the package run metrics from '%s' to '%s'", fromTime, toTime)
	}
	packagesRunCount, err := reporter.store.GetPackagesRunCount(ctx)
	if err != nil {
		return 0, stacktrace.Propagate(err, "an error occurred getting packages run count from the store")
	}

	for packageName, newCount := range newPackagesRunCount {
		finalCount := newCount
		if currentCount, found := packagesRunCount[packageName]; found {
			finalCount = currentCount + newCount
		}
		packagesRunCount[packageName] = finalCount
	}

	if err := reporter.store.UpdatePackagesRunCount(ctx, packagesRunCount); err != nil {
		return 0, stacktrace.Propagate(err, "an error occurred updating packages run count")
	}

	return len(packagesRunCount), nil
}
