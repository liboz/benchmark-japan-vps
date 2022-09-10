# Test Client DB Insert Locally

```
go run client/main.go client/sql.go target "postgres://username:password@localhost:5432/benchmark-japan-vps?sslmode=disable"
```

# Commands to test locally

```
go run server/main.go
curl     --header "Content-Type: application/json"     --data '{}'     http://localhost:8000/service.v1.BenchmarkService/GetResults
```
