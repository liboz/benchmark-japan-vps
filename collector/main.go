package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"

	benchmarkv1 "liboz/benchmark-japan-vps/gen/service/v1"        // generated by protoc-gen-go
	"liboz/benchmark-japan-vps/gen/service/v1/benchmarkv1connect" // generated by protoc-gen-connect-go

	"github.com/bufbuild/connect-go"
)

var DEBUG = false
var TARGETS = []Target{}

type Target struct {
	IpAddress string
	Name      string
}

func parseTargetFromString(targetString string) Target {
	split := strings.Split(targetString, ":")
	return Target{IpAddress: split[0], Name: split[1]}
}

func makeClientFromTarget(target Target) (benchmarkv1connect.BenchmarkServiceClient, context.Context, context.CancelFunc) {
	client := benchmarkv1connect.NewBenchmarkServiceClient(
		http.DefaultClient,
		fmt.Sprintf("http://%s:8000", target.IpAddress),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3) // wait up to 3 minutes
	return client, ctx, cancel
}

func getResult(ch chan []*benchmarkv1.BenchmarkResult, target Target) {
	client, ctx, cancel := makeClientFromTarget(target)
	defer cancel()

	res, err := client.GetResults(
		ctx,
		connect.NewRequest(&benchmarkv1.GetResultsRequest{}),
	)
	if err != nil {
		log.Printf("Error getting results for target %s: %s", target, err)
		ch <- nil
		return
	}

	for _, result := range res.Msg.Results {
		result.IpAddress = target.IpAddress
		result.Name = target.Name
	}

	ch <- res.Msg.Results
}

func deleteOldResults(target Target, endTime int64) {
	client, ctx, cancel := makeClientFromTarget(target)
	defer cancel()

	_, err := client.DeleteOldResults(
		ctx,
		connect.NewRequest(&benchmarkv1.DeleteOldResultsRequest{LatestEndTime: endTime}),
	)
	if err != nil {
		log.Printf("Error deleting old results from target %s with endTime %d: %s", target, endTime, err)
		return
	}
	log.Printf("Deleted old results from target %s with endTime %d", target, endTime)
}

func getResultsFromAllTargets() []*benchmarkv1.BenchmarkResult {
	result := []*benchmarkv1.BenchmarkResult{}

	start := time.Now()
	ch := make(chan []*benchmarkv1.BenchmarkResult)
	returnedTargets := 0

	for _, target := range TARGETS {
		go getResult(ch, target)
	}

	for res := range ch {
		returnedTargets += 1
		if res != nil {
			result = append(result, res...)
		}
		if returnedTargets == len(TARGETS) {
			break
		}
	}

	elapsedTime := time.Since(start)
	log.Printf("Took %s to get information from all targets %s", elapsedTime, TARGETS)

	return result
}

func insertSpeedTestResults(txn *sql.Tx, benchmarkId int64, speedTestResults []*benchmarkv1.SpeedTestResult) error {
	var (
		placeholders []string
		vals         []interface{}
	)
	for index, speedTestResult := range speedTestResults {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d)",
			index*5+1,
			index*5+2,
			index*5+3,
			index*5+4,
			index*5+5,
		))

		vals = append(vals, benchmarkId, speedTestResult.City, speedTestResult.Country, speedTestResult.UploadSpeed, speedTestResult.DownloadSpeed)
	}

	insertStatement := speedTestResultsSql(placeholders)
	_, err := txn.Exec(insertStatement, vals...)
	if err != nil {
		log.Printf("Error inserting speed_test_results: %s, Attempted to insert with statement [%s] and vals [%s]", err, insertStatement, vals)
	}

	return err
}

func insertPingTestResults(txn *sql.Tx, benchmarkId int64, pingTestResults []*benchmarkv1.PingTestResult) error {
	var (
		placeholders []string
		vals         []interface{}
	)
	for index, pingTestResult := range pingTestResults {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			index*7+1,
			index*7+2,
			index*7+3,
			index*7+4,
			index*7+5,
			index*7+6,
			index*7+7,
		))

		vals = append(vals, benchmarkId, pingTestResult.Url, pingTestResult.DroppedPackets, pingTestResult.MinimumPing, pingTestResult.AveragePing, pingTestResult.MaximumPing, pingTestResult.StandardDeviation)
	}

	insertStatement := pingTestResultsSql(placeholders)
	_, err := txn.Exec(insertStatement, vals...)
	if err != nil {
		log.Printf("Error inserting ping_test_result: %s, Attempted to insert with statement [%s] and vals [%s]", err, insertStatement, vals)
	}

	return err
}

func insertIntoPostgres(db *sql.DB, newResults []*benchmarkv1.BenchmarkResult) {
	start := time.Now()
	for _, result := range newResults {
		placeholder := fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			1,
			2,
			3,
			4,
			5,
			6,
			7,
		)

		txn, err := db.Begin()
		if err != nil {
			continue
		}

		insertStatement := benchmarkResultsSql(placeholder)
		vals := []interface{}{result.IpAddress, result.Name, result.StartTime, result.EndTime,
			result.IoSpeed, result.SingleCoreGeekbench, result.MultiCoreGeekbench}
		row := txn.QueryRow(insertStatement, vals...)

		var lastInsertId int64 = 0
		err = row.Scan(&lastInsertId)
		if err != nil {
			log.Printf("Error inserting benchmark_results: %s, Attempted to insert with statement [%s] and vals [%s]", err, insertStatement, vals)
			txn.Rollback()
			continue
		}

		if len(result.SpeedTestResults) > 0 {
			err = insertSpeedTestResults(txn, lastInsertId, result.SpeedTestResults)

			if err != nil {
				txn.Rollback()
				continue
			}
		}

		err = insertPingTestResults(txn, lastInsertId, result.PingTestResults)

		if err != nil {
			txn.Rollback()
			continue
		}

		if err := txn.Commit(); err != nil {
			log.Print(err)
		}

		log.Printf("Inserted %s as %d", result, lastInsertId)
		go func(result *benchmarkv1.BenchmarkResult) {
			deleteOldResults(Target{IpAddress: result.IpAddress, Name: result.Name}, result.EndTime)
		}(result)
	}
	elapsedTime := time.Since(start)
	log.Printf("Took %s to insert gathered data into postgres", elapsedTime)
}

func createTablesIfNotExisting(db *sql.DB) error {
	txn, err := db.Begin()

	if err != nil {
		log.Print("Could not begin transaction to create initial tables", err)
		return err
	}

	_, err = txn.Exec(INITIAL_TABLE_SQL)
	if err != nil {
		log.Print("Could not create initial tables", err)
		return err
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

func main() {
	argsWithoutProg := os.Args[1:]
	targetStrings := strings.Split(argsWithoutProg[0], ",")
	for _, targetString := range targetStrings {
		TARGETS = append(TARGETS, parseTargetFromString(targetString))
	}

	if len(argsWithoutProg) == 3 {
		DEBUG, _ = strconv.ParseBool(argsWithoutProg[2])
	}

	log.Printf("Starting with DEBUG: [%v], TARGETS: %s\n", DEBUG, TARGETS)

	// debug mode to manually run benchmark
	if DEBUG {
		client := benchmarkv1connect.NewBenchmarkServiceClient(
			http.DefaultClient,
			fmt.Sprintf("http://%s:8000", TARGETS[0]),
		)
		res, err := client.StartBenchmark(
			context.Background(),
			connect.NewRequest(&benchmarkv1.StartBenchmarkRequest{}),
		)

		if err != nil {
			log.Println(err)
			return
		}
		log.Println(res)
		return
	}

	connStr := argsWithoutProg[1]
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = createTablesIfNotExisting(db)
	if err != nil {
		log.Fatal(err)
	}

	// check results regularly
	for {
		currTime := time.Now().Format("2006-01-02 15:04:05")
		newResults := getResultsFromAllTargets()
		if len(newResults) == 0 {
			fmt.Printf("%s: found no new results\n", currTime)
		} else {
			insertIntoPostgres(db, newResults)
		}
		time.Sleep(10 * time.Minute)
	}
}
