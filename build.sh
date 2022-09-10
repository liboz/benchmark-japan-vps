#!/usr/bin/env bash
set -e
# reads from tools.go and installs all tools
cat tools.go | grep _ | awk -F'"' '{print $2}' | xargs -tI % go install %

buf generate
CGO_ENABLED=0 GOOS=linux go build -o ./benchmark-japan-vps -a -ldflags '-w -s' ./server/main.go
tar -czf benchmark-japan-vps.tar.gz benchmark-japan-vps bench.sh parse_benchmark_result.sh run_and_parse_bench.sh server/benchmark-server.service