package metrics

import (
	"context"
	"github.com/kurtosis-tech/stacktrace"
)

type PackagesRunCount map[string]uint32

type Reporter struct {
	snowflake *snowflake
}

func CreateReporter() (*Reporter, error) {
	snowflakeObj, err := createSnowflake()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred creating the Snowflake object")
	}

	newMetricsReporter := &Reporter{snowflake: snowflakeObj}

	return newMetricsReporter, nil
}

func (metricsReporter *Reporter) GetPackageRunMetrics(ctx context.Context) (PackagesRunCount, error) {
	runMetricsRows, err := metricsReporter.snowflake.runQueryAndGetPackageRunMetrics(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred running the query to get the package run metrics")
	}
	return runMetricsRows, nil
}
