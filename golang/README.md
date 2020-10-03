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

Run unit tests
```
go test ./...
```

Use this https://github.com/liufuyang/autocannon-go for some simple performance test

```
# at root directory of autocannon-go
go build && ./autocannon-go --connections=10 --pipelining=10 --duration=30 --uri=http://localhost:8081
running 30s test @ http://localhost:8081
10 connections with 10 pipelining factor.
Generated  100000000  random int, with range [0,  1000000 ).

---------- Read/GET ----------

+---------+------+------+-------+------+---------+---------+--------+
|  STAT   | 2.5% | 50%  | 97.5% | 99%  |   AVG   |  STDEV  |  MAX   |
+---------+------+------+-------+------+---------+---------+--------+
| Latency | 0 ms | 1 ms | 2 ms  | 3 ms | 1.06 ms | 2.63 ms | 363 ms |
+---------+------+------+-------+------+---------+---------+--------+


+-----------+--------+--------+--------+--------+----------+---------+--------+
|   STAT    |   1%   |  2.5%  |  50%   | 97.5%  |   AVG    |  STDEV  |  MIN   |
+-----------+--------+--------+--------+--------+----------+---------+--------+
| Req/Sec   |  31890 |  31890 |  51789 |  58167 | 51299.63 | 5046.04 |  31890 |
| Bytes/Sec | 3.0 MB | 3.0 MB | 4.9 MB | 5.5 MB | 4.8 MB   | 473 kB  | 3.0 MB |
+-----------+--------+--------+--------+--------+----------+---------+--------+

Req/Bytes counts sampled once per second.


1538990 2xx responses, 0 non 2xx responses.
1539k total requests in 30 seconds, 144 MB read for GET.

---------- Write/POST ----------

+---------+------+------+-------+------+---------+---------+--------+
|  STAT   | 2.5% | 50%  | 97.5% | 99%  |   AVG   |  STDEV  |  MAX   |
+---------+------+------+-------+------+---------+---------+--------+
| Latency | 0 ms | 1 ms | 2 ms  | 3 ms | 1.07 ms | 3.41 ms | 363 ms |
+---------+------+------+-------+------+---------+---------+--------+


+-----------+--------+--------+--------+--------+----------+---------+--------+
|   STAT    |   1%   |  2.5%  |  50%   | 97.5%  |   AVG    |  STDEV  |  MIN   |
+-----------+--------+--------+--------+--------+----------+---------+--------+
| Req/Sec   |   8000 |   8000 |  12957 |  14427 | 12835.43 | 1267.15 |   8000 |
| Bytes/Sec | 752 kB | 752 kB | 1.2 MB | 1.4 MB | 1.2 MB   | 119 kB  | 752 kB |
+-----------+--------+--------+--------+--------+----------+---------+--------+

Req/Bytes counts sampled once per second.


385063 2xx responses, 0 non 2xx responses.
385k total requests in 30 seconds, 36 MB read for POST.
Done!
```
