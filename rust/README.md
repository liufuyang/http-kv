## Rust implementation

Rust result:
```
./autocannon-go --connections=20 --pipelining=10 --duration=30 --uri=http://localhost:8082
running 30s test @ http://localhost:8082
20 connections with 10 pipelining factor.
Generated  100000000  random int, with range [0,  1000000 ).

---------- Read/GET ----------

+---------+------+------+-------+------+---------+---------+-------+
|  STAT   | 2.5% | 50%  | 97.5% | 99%  |   AVG   |  STDEV  |  MAX  |
+---------+------+------+-------+------+---------+---------+-------+
| Latency | 1 ms | 1 ms | 2 ms  | 2 ms | 1.22 ms | 0.49 ms | 16 ms |
+---------+------+------+-------+------+---------+---------+-------+

+-----------+--------+--------+--------+--------+----------+---------+--------+
|   STAT    |   1%   |  2.5%  |  50%   | 97.5%  |   AVG    |  STDEV  |  MIN   |
+-----------+--------+--------+--------+--------+----------+---------+--------+
| Req/Sec   |  86732 |  86732 |  93570 | 104207 | 94284.10 | 4553.46 |  86732 |
| Bytes/Sec | 8.1 MB | 8.1 MB | 8.7 MB | 9.7 MB | 8.8 MB   | 418 kB  | 8.1 MB |
+-----------+--------+--------+--------+--------+----------+---------+--------+

Req/Bytes counts sampled once per second.

2828523 2xx responses, 0 non 2xx responses.
2829k total requests in 30 seconds, 264 MB read for GET.


---------- Write/POST ----------

+---------+------+------+-------+------+---------+---------+-------+
|  STAT   | 2.5% | 50%  | 97.5% | 99%  |   AVG   |  STDEV  |  MAX  |
+---------+------+------+-------+------+---------+---------+-------+
| Latency | 1 ms | 1 ms | 2 ms  | 2 ms | 1.24 ms | 0.49 ms | 11 ms |
+---------+------+------+-------+------+---------+---------+-------+

+-----------+--------+--------+--------+--------+----------+---------+--------+
|   STAT    |   1%   |  2.5%  |  50%   | 97.5%  |   AVG    |  STDEV  |  MIN   |
+-----------+--------+--------+--------+--------+----------+---------+--------+
| Req/Sec   |  21642 |  21642 |  23313 |  25979 | 23546.57 | 1108.94 |  21642 |
| Bytes/Sec | 2.0 MB | 2.0 MB | 2.2 MB | 2.4 MB | 2.2 MB   | 103 kB  | 2.0 MB |
+-----------+--------+--------+--------+--------+----------+---------+--------+

Req/Bytes counts sampled once per second.

706397 2xx responses, 0 non 2xx responses.
706k total requests in 30 seconds, 66 MB read for POST.
Done!

```

Golang result:
```
./autocannon-go --connections=20 --pipelining=10 --duration=30 --uri=http://localhost:8081
running 30s test @ http://localhost:8081
20 connections with 10 pipelining factor.
Generated  100000000  random int, with range [0,  1000000 ).

---------- Read/GET ----------

+---------+------+------+-------+------+---------+---------+--------+
|  STAT   | 2.5% | 50%  | 97.5% | 99%  |   AVG   |  STDEV  |  MAX   |
+---------+------+------+-------+------+---------+---------+--------+
| Latency | 1 ms | 3 ms | 4 ms  | 4 ms | 2.59 ms | 2.30 ms | 112 ms |
+---------+------+------+-------+------+---------+---------+--------+

+-----------+--------+--------+--------+--------+----------+---------+--------+
|   STAT    |   1%   |  2.5%  |  50%   | 97.5%  |   AVG    |  STDEV  |  MIN   |
+-----------+--------+--------+--------+--------+----------+---------+--------+
| Req/Sec   |  42475 |  42475 |  52071 |  55684 | 51632.70 | 2954.81 |  42475 |
| Bytes/Sec | 4.0 MB | 4.0 MB | 4.8 MB | 5.2 MB | 4.8 MB   | 273 kB  | 4.0 MB |
+-----------+--------+--------+--------+--------+----------+---------+--------+

Req/Bytes counts sampled once per second.

1548981 2xx responses, 0 non 2xx responses.
1549k total requests in 30 seconds, 144 MB read for GET.


---------- Write/POST ----------

+---------+------+------+-------+------+---------+---------+--------+
|  STAT   | 2.5% | 50%  | 97.5% | 99%  |   AVG   |  STDEV  |  MAX   |
+---------+------+------+-------+------+---------+---------+--------+
| Latency | 1 ms | 3 ms | 4 ms  | 4 ms | 2.65 ms | 2.82 ms | 112 ms |
+---------+------+------+-------+------+---------+---------+--------+

+-----------+--------+--------+--------+--------+----------+--------+--------+
|   STAT    |   1%   |  2.5%  |  50%   | 97.5%  |   AVG    | STDEV  |  MIN   |
+-----------+--------+--------+--------+--------+----------+--------+--------+
| Req/Sec   |  10465 |  10465 |  12992 |  13834 | 12907.97 | 755.47 |  10465 |
| Bytes/Sec | 984 kB | 984 kB | 1.2 MB | 1.3 MB | 1.2 MB   | 71 kB  | 984 kB |
+-----------+--------+--------+--------+--------+----------+--------+--------+

Req/Bytes counts sampled once per second.

387239 2xx responses, 0 non 2xx responses.
387k total requests in 30 seconds, 36 MB read for POST.
Done!
```