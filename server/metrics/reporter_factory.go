package metrics

import (
	"context"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/utils"
	"github.com/kurtosis-tech/stacktrace"
)

const (
	ciEnvVarKey = "CI"
	runningInCI = "true"
)

func CreateAndScheduleReporter(ctx context.Context, store store.KurtosisIndexerStore) (Reporter, error) {

	var newMetricsReporter Reporter

	isRunningInCI, err := utils.GetFromEnvVar(ciEnvVarKey, "Running in CI")
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting env var key value with key '%s'", ciEnvVarKey)
	}

	if isRunningInCI == runningInCI {
		// to avoid calling Snowflake from the CI builds
		newMetricsReporter = &doNothingReporter{}
	} else {
		snowflakeObj, err := createSnowflake()
		if err != nil {
			return nil, stacktrace.Propagate(err, "an error occurred creating the Snowflake object")
		}

		newMetricsReporter = &reporterImpl{
			ctx:       ctx,
			snowflake: snowflakeObj,
			ticker:    nil,
			store:     store,
		}
	}

	if err := newMetricsReporter.Schedule(false); err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred scheduling the metrics reporter")
	}

	return newMetricsReporter, nil
}
