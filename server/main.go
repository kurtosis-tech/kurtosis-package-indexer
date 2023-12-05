package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/kurtosis-tech/kurtosis-package-indexer/api/golang/generated/generatedconnect"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/crawler"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/resource"
	"github.com/kurtosis-tech/kurtosis-package-indexer/server/store"
	connect_server "github.com/kurtosis-tech/kurtosis/connect-server"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	sf "github.com/snowflakedb/gosnowflake"
	"log"
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

// for sharing a single DB instance
var db *sql.DB

func main() {
	conn, err := GetConnection()
	if err != nil {
		log.Fatal("Error while connection ", err)
		return
	}
	defer conn.Close()

	fmt.Println("Successfully get the connection..")
}

func GetConnection() (conn *sql.Conn, err error) {

	dns, err := sf.DSN(&sf.Config{
		Account:   "qtjzlxq-us27029",
		User:      "",
		Password:  "",
		Database:  "SEGMENT_EVENTS",
		Schema:    "",
		Warehouse: "",
		Region:    "",
		Role:      "PRODUCT_ANALYTICS_READER",
	})
	if err != nil {
		log.Fatal("Error while DNS string: ", err)
		return conn, err
	}

	db, err := sql.Open("snowflake", dns)
	if err != nil {
		log.Fatal("Error while open DB: ", err)
		return conn, err
	}
	defer db.Close()

	db.SetMaxIdleConns(100)
	db.SetMaxOpenConns(50)
	db.SetConnMaxLifetime(3 * time.Second)

	log.Println("Database Connection Successful..!")

	ctx := context.Background()
	conn, err = db.Conn(ctx)
	if err != nil {
		log.Fatal("Error while open connection: ", err)
		return conn, err
	}

	log.Println("Connection Successful..!")

	query := "SELECT 1"
	rows, err := conn.QueryContext(ctx, query) // no cancel is allowed
	if err != nil {
		log.Fatalf("failed to run a query. %v, err: %v", query, err)
	}
	defer rows.Close()
	var v int
	for rows.Next() {
		err := rows.Scan(&v)
		if err != nil {
			log.Fatalf("failed to get result. err: %v", err)
		}
		if v != 1 {
			log.Fatalf("failed to get 1. got: %v", v)
		}
	}
	if rows.Err() != nil {
		fmt.Printf("ERROR: %v\n", rows.Err())
		return
	}
	fmt.Printf("Congrats! You have successfully run %v with Snowflake DB!\n", query)

	return conn, err

}

/*
func main() {
	ctx := context.Background()
	configureLogger()

	// Set up the store which will store all the packages. For now all in memory
	indexerStore, err := store.InstantiateStoreFromEnvVar()
	if err != nil {
		exitFailure(err)
	}
	defer indexerStore.Close()

	// Set up the crawler which will populate the store on a periodical basis
	indexerCtx, cancelFunc := context.WithCancel(ctx)
	indexerCrawler := crawler.NewGithubCrawler(indexerCtx, indexerStore)
	if err := indexerCrawler.Schedule(false); err != nil {
		exitFailure(err)
	}
	defer cancelFunc()

	if err := runServer(indexerStore, indexerCrawler); err != nil {
		exitFailure(err)
	}
	logrus.Exit(successExitCode)
}*/

func runServer(indexerStore store.KurtosisIndexerStore, indexerCrawler *crawler.GithubCrawler) error {
	kurtosisPackageIndexerResource := resource.NewKurtosisPackageIndexer(indexerStore, indexerCrawler)
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

func exitFailure(err error) {
	logrus.Error(err.Error())
	logrus.Exit(failureExitCode)
}
