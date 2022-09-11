# Purpose

This repo is used to store a program used to benchmark vps performance. The main benchmark script is a combination of https://github.com/teddysun/across/blob/master/bench.sh and https://github.com/masonr/yet-another-bench-script

Sample kubernetes deployment template at [docs/sample_kubernetes_deployment.yaml](docs/sample_kubernetes_deployment.yaml)

## Test Client DB Insert Locally

```
go run collector/main.go collector/sql.go 127.0.0.1:localhost,127.0.0.2:localhost2 "postgres://username:password@localhost:5432/benchmark-japan-vps?sslmode=disable"
```

## Commands to test locally

```
go run server/main.go
curl     --header "Content-Type: application/json"     --data '{}'     http://localhost:8000/service.v1.BenchmarkService/GetResults
```

## Running in a server

Copy over the start script:

```
scp run_server.sh root@ip_address:/root
```

May need to open the port with iptables:

```
iptables -A INPUT -p tcp --dport 8000 -j ACCEPT
```
