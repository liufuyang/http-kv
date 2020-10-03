# Simple HTTP KV server in Golang

Start server:
```
go run main.go

INFO[2020-10-03T11:07:10+02:00] cache expire time: 10000ms
INFO[2020-10-03T11:07:10+02:00] Starting server on 8081 ...
```

Test server manually
```
curl -X POST localhost:8081/k1 -d 'v1'
curl -X POST localhost:8081/size
curl -X GET localhost:8081/k1
```

Use this https://github.com/liufuyang/autocannon-go for some simple performance test

```
# at root directory of autocannon-go
go build && ./autocannon-go --connections=1 --pipelining=1 --duration=10 --uri=http://localhost:8081
```

Run unit tests
```
go test ./...
```