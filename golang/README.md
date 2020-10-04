# Simple HTTP KV server in Golang

### Start server:
```
go run main.go

INFO[2020-10-03T11:07:10+02:00] cache expire time: 10000ms
INFO[2020-10-03T11:07:10+02:00] Starting server on 8081 ...
```

### Test server manually
```
curl -X POST localhost:8081/k1 -d 'v1'
curl -X POST localhost:8081/size
curl -X GET localhost:8081/k1
```

### Run unit tests
```
go test ./...
```

### Simple Performance test

Use this https://github.com/liufuyang/autocannon-go for some simple performance test, getting some metrics from client size

```
# at root directory of autocannon-go
go build && ./autocannon-go --connections=20 --pipelining=10 --duration=300 --uri=http://localhost:8081
running 300s test @ http://localhost:8081
20 connections with 10 pipelining factor.
Generated  100000000  random int, with range [0,  1000000 ).

---------- Read/GET ----------

+---------+------+------+-------+------+---------+---------+--------+
|  STAT   | 2.5% | 50%  | 97.5% | 99%  |   AVG   |  STDEV  |  MAX   |
+---------+------+------+-------+------+---------+---------+--------+
| Latency | 1 ms | 3 ms | 4 ms  | 7 ms | 2.79 ms | 3.19 ms | 295 ms |
+---------+------+------+-------+------+---------+---------+--------+

+-----------+--------+--------+--------+--------+----------+---------+--------+
|   STAT    |   1%   |  2.5%  |  50%   | 97.5%  |   AVG    |  STDEV  |  MIN   |
+-----------+--------+--------+--------+--------+----------+---------+--------+
| Req/Sec   |  36447 |  38619 |  49174 |  55123 | 48620.78 | 4017.16 |  21414 |
| Bytes/Sec | 3.4 MB | 3.6 MB | 4.6 MB | 5.1 MB | 4.5 MB   | 373 kB  | 2.0 MB |
+-----------+--------+--------+--------+--------+----------+---------+--------+

Req/Bytes counts sampled once per second.

14550044 2xx responses, 0 non 2xx responses.
14550k total requests in 300 seconds, 1.4 GB read for GET.


---------- Write/POST ----------

+---------+------+------+-------+------+---------+---------+--------+
|  STAT   | 2.5% | 50%  | 97.5% | 99%  |   AVG   |  STDEV  |  MAX   |
+---------+------+------+-------+------+---------+---------+--------+
| Latency | 1 ms | 3 ms | 4 ms  | 8 ms | 2.82 ms | 3.67 ms | 293 ms |
+---------+------+------+-------+------+---------+---------+--------+

+-----------+--------+--------+--------+--------+----------+---------+--------+
|   STAT    |   1%   |  2.5%  |  50%   | 97.5%  |   AVG    |  STDEV  |  MIN   |
+-----------+--------+--------+--------+--------+----------+---------+--------+
| Req/Sec   |   8982 |   9641 |  12263 |  13887 | 12154.60 | 1011.21 |   5422 |
| Bytes/Sec | 844 kB | 906 kB | 1.2 MB | 1.3 MB | 1.1 MB   | 95 kB   | 510 kB |
+-----------+--------+--------+--------+--------+----------+---------+--------+

Req/Bytes counts sampled once per second.

3637402 2xx responses, 0 non 2xx responses.
3637k total requests in 300 seconds, 342 MB read for POST.
Done!
```

## Prometheus and Grafana

See the [prometheus](..//prometheus) folder for details to see how to start Prometheus and Grafana to monitor the server metrics

A graph during some performance testing
![demo](../demo.png)