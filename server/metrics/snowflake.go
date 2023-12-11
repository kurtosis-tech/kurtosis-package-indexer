package metrics

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/types"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	sf "github.com/snowflakedb/gosnowflake"
	"os"
	"time"
)

const (
	// kurtosisSnowflakeAccountIdentifierEnvVarKey using this format: https://docs.snowflake.com/en/user-guide/admin-account-identifier#format-1-preferred-account-name-in-your-organization
	kurtosisSnowflakeAccountIdentifierEnvVarKey = "KURTOSIS_SNOWFLAKE_ACCOUNT_IDENTIFIER"
	// kurtosisSnowflakeUserEnvVarKey should be a user with only read access to public metrics
	kurtosisSnowflakeUserEnvVarKey      = "KURTOSIS_SNOWFLAKE_USER"
	kurtosisSnowflakePasswordEnvVarKey  = "KURTOSIS_SNOWFLAKE_PASSWORD"
	kurtosisSnowflakeDatabaseEnvVarKey  = "KURTOSIS_SNOWFLAKE_DB"
	kurtosisSnowflakeWarehouseEnvVarKey = "KURTOSIS_SNOWFLAKE_WAREHOUSE"
	kurtosisSnowflakeRoleEnvVarKey      = "KURTOSIS_SNOWFLAKE_ROLE"

	snowflakeDriverName = "snowflake"

	//snowflakeDBIdleConnections this be small because they will be used only for the job task, and it won't be run over each server request
	snowflakeDBIdleConnections = 10
	//snowflakeDBMaxOpenConnections this be small because they will be used only for the job task, and it won't be run over each server request
	snowflakeDBMaxOpenConnections = 5
	snowflakeDBConnMaxLifeTime    = 5 * time.Second
	snowflakeQueryTimeout         = 60 * time.Second

	snowflakeTimestampFormat = "2006-01-02 15:04:05"

	kurtosianUserSQLQuery = "SELECT USER_ID FROM SEGMENT_EVENTS.KURTOSIS_METRICS_LIBRARY.KNOWN_USERS WHERE IS_KURTOSIAN=TRUE"
	isCISQLQueryCondition = "FALSE"

	// This query is pretty much the same we have in the "Top Usage Packages" table in the Snowflake "usage metrics dashboard"
	selectPackageRunMetricSQLQueryFormat = `SELECT IFNULL(kp.name, k.package_id) as package_name, COUNT(k.PACKAGE_ID) AS COUNT 
FROM SEGMENT_EVENTS.KURTOSIS_METRICS_LIBRARY.KURTOSIS_RUN k 
LEFT JOIN SEGMENT_EVENTS.KURTOSIS_METRICS_LIBRARY.KNOWN_PACKAGES kp ON k.PACKAGE_ID = kp.PACKAGE_ID 
WHERE k.USER_ID NOT IN ( %s ) AND k.IS_CI = %s
AND (ORIGINAL_TIMESTAMP >= ('%s')::timestamp AND ORIGINAL_TIMESTAMP < ('%s')::timestamp)
GROUP BY package_name;`
)

type snowflake struct {
	dsn string
}

func createSnowflake() (*snowflake, error) {
	dsn, err := getSnowflakeDSN()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the Snowflake DSN")
	}

	newSnowflake := &snowflake{
		dsn: dsn,
	}

	return newSnowflake, nil
}

func (snowflake *snowflake) getPackageRunMetricsInDateRange(ctx context.Context, fromTime time.Time, toTime time.Time) (types.PackagesRunCount, error) {
	conn, err := snowflake.getConnection(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the Snowflake dB connection")
	}

	ctxWithTimeout, _ := context.WithTimeout(ctx, snowflakeQueryTimeout)

	query := fmt.Sprintf(
		selectPackageRunMetricSQLQueryFormat,
		kurtosianUserSQLQuery,
		isCISQLQueryCondition,
		fromTime.Format(snowflakeTimestampFormat),
		toTime.Format(snowflakeTimestampFormat),
	)
	rows, err := conn.QueryContext(ctxWithTimeout, query)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred running the query for getting the run metrics")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Warningf("an error occurred closing the query rows")
		}
	}()
	var packageName string
	var count uint32
	result := types.PackagesRunCount{}
	for rows.Next() {
		if err := rows.Scan(&packageName, &count); err != nil {
			return nil, stacktrace.Propagate(err, "an error occurred scanning the query result rows")
		}
		result[packageName] = count
	}
	if rows.Err() != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the query result rows")
	}
	logrus.Debugf("run metrics query successfully executed, '%v' packages received", len(result))

	return result, nil
}

func (snowflake *snowflake) getConnection(ctx context.Context) (conn *sql.Conn, err error) {

	db, err := snowflake.getSnowflakeDB()
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the Snowflake dB")
	}

	logrus.Debugf("Connecting with Snowflake database...")
	conn, err = db.Conn(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred connecting to the Snowflake dB")
	}
	logrus.Debugf("...successful connection.")

	return conn, err
}

func (snowflake *snowflake) getSnowflakeDB() (*sql.DB, error) {

	db, err := sql.Open(snowflakeDriverName, snowflake.dsn)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred opening the Snowflake dB")
	}
	defer func() {
		/*if err := db.Close(); err != nil {
			logrus.Warningf("an error occurred closing the Snowflake dB")
		}*/ // TODO handle this
	}()

	// these values should be small because they will be used only for the job task,
	// and it won't be run over each server request
	db.SetMaxIdleConns(snowflakeDBIdleConnections)
	db.SetMaxOpenConns(snowflakeDBMaxOpenConnections)
	db.SetConnMaxLifetime(snowflakeDBConnMaxLifeTime)

	return db, nil
}

func getSnowflakeDSN() (string, error) {
	snowflakeAccount, err := getSnowflakeAccountFromEnvVar()
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred getting the Snowflake account identifier")
	}

	snowflakeUser, err := getSnowflakeUserFromEnvVar()
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred getting the Snowflake user")
	}

	snowflakePassword, err := getSnowflakePasswordFromEnvVar()
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred getting the Snowflake password")
	}

	snowflakeDatabase, err := getSnowflakeDatabaseFromEnvVar()
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred getting the Snowflake database")
	}

	snowflakeWarehouse, err := getSnowflakeWarehouseFromEnvVar()
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred getting the Snowflake warehouse")
	}

	snowflakeRole, err := getSnowflakeRoleFromEnvVar()
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred getting the Snowflake role")
	}

	dsn, err := sf.DSN(&sf.Config{
		Account:   snowflakeAccount,
		User:      snowflakeUser,
		Password:  snowflakePassword,
		Database:  snowflakeDatabase,
		Warehouse: snowflakeWarehouse,
		Role:      snowflakeRole,
	})
	if err != nil {
		return "", stacktrace.Propagate(err, "an error occurred configuring DSN for the Snowflake dB connection")
	}
	return dsn, nil
}

func getSnowflakeAccountFromEnvVar() (string, error) {
	return getFromEnvVar(kurtosisSnowflakeAccountIdentifierEnvVarKey, "Snowflake account identifier")
}

func getSnowflakeUserFromEnvVar() (string, error) {
	return getFromEnvVar(kurtosisSnowflakeUserEnvVarKey, "Snowflake user")
}

func getSnowflakePasswordFromEnvVar() (string, error) {
	return getFromEnvVar(kurtosisSnowflakePasswordEnvVarKey, "Snowflake password")
}

func getSnowflakeDatabaseFromEnvVar() (string, error) {
	return getFromEnvVar(kurtosisSnowflakeDatabaseEnvVarKey, "Snowflake database")
}

func getSnowflakeWarehouseFromEnvVar() (string, error) {
	return getFromEnvVar(kurtosisSnowflakeWarehouseEnvVarKey, "Snowflake warehouse")
}

func getSnowflakeRoleFromEnvVar() (string, error) {
	return getFromEnvVar(kurtosisSnowflakeRoleEnvVarKey, "Snowflake role")
}

func getFromEnvVar(
	key string,
	subject string,
) (string, error) {
	value := os.Getenv(key)
	if len(value) < 1 {
		return "", stacktrace.NewError("No '%s' env var was found. Must be provided as env var %s", subject, key)
	}
	logrus.Debugf("Successfully loaded env var '%s'", subject)
	return value, nil
}
