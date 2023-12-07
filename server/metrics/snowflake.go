package metrics

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	sf "github.com/snowflakedb/gosnowflake"
	"os"
	"time"
)

const (
	// KurtosisSnowflakeAccountIdentifierEnvVarKey using this format: https://docs.snowflake.com/en/user-guide/admin-account-identifier#format-1-preferred-account-name-in-your-organization
	KurtosisSnowflakeAccountIdentifierEnvVarKey = "KURTOSIS_SNOWFLAKE_ACCOUNT_IDENTIFIER"
	// KurtosisSnowflakeUserEnvVarKey should be a user with only read access to public metrics
	KurtosisSnowflakeUserEnvVarKey      = "KURTOSIS_SNOWFLAKE_USER"
	KurtosisSnowflakePasswordEnvVarKey  = "KURTOSIS_SNOWFLAKE_PASSWORD"
	KurtosisSnowflakeDatabaseEnvVarKey  = "KURTOSIS_SNOWFLAKE_DB"
	KurtosisSnowflakeWarehouseEnvVarKey = "KURTOSIS_SNOWFLAKE_WAREHOUSE"
	KurtosisSnowflakeRoleEnvVarKey      = "KURTOSIS_SNOWFLAKE_ROLE"

	SnowflakeDriverName = "snowflake"

	//SnowflakeDBIdleConnections this be small because they will be used only for the job task, and it won't be run over each server request
	SnowflakeDBIdleConnections = 10
	//SnowflakeDBMaxOpenConnections this be small because they will be used only for the job task, and it won't be run over each server request
	SnowflakeDBMaxOpenConnections = 5
	SnowflakeDBConnMaxLifeTime    = 5 * time.Second
	SnowflakeQueryTimeout         = 60 * time.Second

	kurtosianUserSQLQuery = "SELECT USER_ID FROM SEGMENT_EVENTS.KURTOSIS_METRICS_LIBRARY.KNOWN_USERS WHERE IS_KURTOSIAN=TRUE"
	isCISQLQueryCondition = "FALSE"

	// This query is pretty much the same we have in the "Top Usage Packages" table in the Snowflake "usage metrics dashboard"
	selectPackageRunMetricSQLQueryFormat = `SELECT IFNULL(kp.name, k.package_id) as package_name, COUNT(k.PACKAGE_ID) AS COUNT 
FROM SEGMENT_EVENTS.KURTOSIS_METRICS_LIBRARY.KURTOSIS_RUN k 
LEFT JOIN SEGMENT_EVENTS.KURTOSIS_METRICS_LIBRARY.KNOWN_PACKAGES kp ON k.PACKAGE_ID = kp.PACKAGE_ID 
                                                                    WHERE k.USER_ID NOT IN ( %s ) AND k.IS_CI = %s 
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

func (snowflake *snowflake) runQueryAndGetPackageRunMetrics(ctx context.Context) (PackagesRunCount, error) {
	conn, err := snowflake.getConnection(ctx)
	if err != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the Snowflake dB connection")
	}

	ctxWithTimeout, _ := context.WithTimeout(ctx, SnowflakeQueryTimeout)

	query := fmt.Sprintf(selectPackageRunMetricSQLQueryFormat, kurtosianUserSQLQuery, isCISQLQueryCondition)
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
	result := PackagesRunCount{}
	var rowCounter int //TODO remove
	for rows.Next() {
		rowCounter++
		if err := rows.Scan(&packageName, &count); err != nil {
			return nil, stacktrace.Propagate(err, "an error occurred scanning the query result rows")
		}
		result[packageName] = count
	}
	if rows.Err() != nil {
		return nil, stacktrace.Propagate(err, "an error occurred getting the query result rows")
	}
	logrus.Debugf("run metrics query successfully executed, '%v' rows received", len(result))

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

	db, err := sql.Open(SnowflakeDriverName, snowflake.dsn)
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
	db.SetMaxIdleConns(SnowflakeDBIdleConnections)
	db.SetMaxOpenConns(SnowflakeDBMaxOpenConnections)
	db.SetConnMaxLifetime(SnowflakeDBConnMaxLifeTime)

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
	return getFromEnvVar(KurtosisSnowflakeAccountIdentifierEnvVarKey, "Snowflake account identifier")
}

func getSnowflakeUserFromEnvVar() (string, error) {
	return getFromEnvVar(KurtosisSnowflakeUserEnvVarKey, "Snowflake user")
}

func getSnowflakePasswordFromEnvVar() (string, error) {
	return getFromEnvVar(KurtosisSnowflakePasswordEnvVarKey, "Snowflake password")
}

func getSnowflakeDatabaseFromEnvVar() (string, error) {
	return getFromEnvVar(KurtosisSnowflakeDatabaseEnvVarKey, "Snowflake database")
}

func getSnowflakeWarehouseFromEnvVar() (string, error) {
	return getFromEnvVar(KurtosisSnowflakeWarehouseEnvVarKey, "Snowflake warehouse")
}

func getSnowflakeRoleFromEnvVar() (string, error) {
	return getFromEnvVar(KurtosisSnowflakeRoleEnvVarKey, "Snowflake role")
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
