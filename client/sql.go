package main

import (
	"fmt"
	"strings"
)

var INITIAL_TABLE_SQL = `CREATE TABLE IF NOT EXISTS benchmark_results (
    benchmark_id bigserial PRIMARY KEY,
    target inet,
    name VARCHAR(255),
    start_time bigint,
    end_time bigint,
    io_speed bigint,
    single_core_geekbench bigint,
    multi_core_geekbench bigint
);

CREATE TABLE IF NOT EXISTS speed_test_results (
    speed_test_result_id bigserial PRIMARY KEY,
    benchmark_id bigint,
    city VARCHAR(255),
    country VARCHAR(255),
    upload_speed bigint,
    download_speed bigint,
    CONSTRAINT fk_benchmark
      FOREIGN KEY(benchmark_id)
	  REFERENCES benchmark_results(benchmark_id)
);

CREATE TABLE IF NOT EXISTS ping_test_results (
    ping_test_result_id bigserial PRIMARY KEY,
    benchmark_id bigint,
    url VARCHAR(255),
    dropped_packets bigint,
    minimum_ping double precision,
    average_ping double precision,
    maximum_ping double precision,
    standard_deviation double precision,
    CONSTRAINT fk_benchmark
      FOREIGN KEY(benchmark_id)
	  REFERENCES benchmark_results(benchmark_id)
);`

var BENCHMARK_RESULT_SQL = `INSERT INTO 
benchmark_results(target,start_time,end_time,io_speed,single_core_geekbench,multi_core_geekbench) 
VALUES %s RETURNING benchmark_id`

var PING_TEST_RESULTS_SQL = `INSERT INTO 
ping_test_results(benchmark_id,url,dropped_packets,minimum_ping,average_ping,maximum_ping,standard_deviation) 
VALUES %s`

var SPEED_TEST_RESULTS_SQL = `INSERT INTO 
speed_test_results(benchmark_id,city,country,upload_speed,download_speed) 
VALUES %s`

func benchmarkResultsSql(placeholder string) string {
	return fmt.Sprintf(BENCHMARK_RESULT_SQL, placeholder)
}

func pingTestResultsSql(placeholders []string) string {
	return fmt.Sprintf(PING_TEST_RESULTS_SQL, strings.Join(placeholders, ","))
}

func speedTestResultsSql(placeholders []string) string {
	return fmt.Sprintf(SPEED_TEST_RESULTS_SQL, strings.Join(placeholders, ","))
}
