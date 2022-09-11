# Test Client DB Insert Locally

```
go run client/main.go client/sql.go 127.0.0.1:localhost "postgres://username:password@localhost:5432/benchmark-japan-vps?sslmode=disable"
```

# Commands to test locally

```
go run server/main.go
curl     --header "Content-Type: application/json"     --data '{}'     http://localhost:8000/service.v1.BenchmarkService/GetResults
```
