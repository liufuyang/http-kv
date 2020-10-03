# Simple HTTP KV server in Golang

Start server:
```
go run main.go
Starting server on 8081 ...
```

Use this https://github.com/liufuyang/autocannon-go for some simple performance test

```
# at root directory of autocannon-go
go build && ./autocannon-go --connections=1 --pipelining=1 --duration=10 --uri=http://localhost:8081/good
```

Run unit tests
```
go test ./...
```