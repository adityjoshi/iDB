# iDB

A toy database built in Golang to understand how redis work internally.



## Features (Current)
	•	RESP protocol parsing (partial)
	•	TCP server (synchronous)
	•	Basic command handling:
	•	PING
	•	Minimal encoding/decoding layer



## Benchmark

Tested using redis-benchmark:

```redis-benchmark -n 10000 -t ping_mbulk -c 1 -h localhost -p 1926```

```
Results:
	redis-benchmark -n 10000 -t ping_mbulk -c 1 -h localhost -p 1926                                                                                                                            36s  04:10:22 PM
WARNING: Could not fetch server CONFIG
====== PING_MBULK ======
  10000 requests completed in 0.18 seconds
  1 parallel clients
  3 bytes payload
  keep alive: 1
  multi-thread: no

Latency by percentile distribution:
0.000% <= 0.015 milliseconds (cumulative count 8175)
87.500% <= 0.023 milliseconds (cumulative count 9277)
93.750% <= 0.031 milliseconds (cumulative count 9490)
96.875% <= 0.039 milliseconds (cumulative count 9777)
98.438% <= 0.047 milliseconds (cumulative count 9928)
99.609% <= 0.055 milliseconds (cumulative count 9989)
99.902% <= 0.063 milliseconds (cumulative count 9995)
99.951% <= 0.071 milliseconds (cumulative count 9997)
99.976% <= 0.087 milliseconds (cumulative count 9998)
99.988% <= 0.103 milliseconds (cumulative count 9999)
99.994% <= 0.119 milliseconds (cumulative count 10000)
100.000% <= 0.119 milliseconds (cumulative count 10000)

Cumulative distribution of latencies:
99.990% <= 0.103 milliseconds (cumulative count 9999)
100.000% <= 0.207 milliseconds (cumulative count 10000)

Summary:
  throughput summary: 55248.62 requests per second
  latency summary (msec):
          avg       min       p50       p95       p99       max
        0.015     0.008     0.015     0.039     0.047     0.119
```
```
Environment:
	•	Single client
	•	No pipelining
	•	No concurrency

```

## Getting Started
```
1. Clone the repo

git clone https://github.com/adityjoshi/iDB.git
cd iDB

2. Run the server

go run main.go

Server starts on:

0.0.0.0:1926

```

## Testing
```
Using Redis CLI:

redis-cli -p 1926                                                                                                                                                                                 04:04:38 PM
127.0.0.1:1926> PING
PONG
127.0.0.1:1926> PING HELLO
"HELLO"
127.0.0.1:1926> PING HELLO WO
(error) Err wrong number of argumnents for ping command

```



