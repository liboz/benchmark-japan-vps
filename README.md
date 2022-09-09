# Commands to test locally

```
go run server/main.go
curl     --header "Content-Type: application/json"     --data '{}'     http://localhost:8000/service.v1.BenchmarkService/GetResult
```
