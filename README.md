# OC contest in Golang

# Synopsis

```sh
go install ./...
server-p0 -h 127.0.0.1 -p 8080 -n 2
```

Checkup:

```sh
$ curl --request POST --data @post.txt localhost:8080
{"current_time":"2018-11-06 00:58:45 +0300","first_name":"john 527bd5b5d689e2c32ae974c6229ff785","id":"test 1234","last_name":"hopkins 99b1084a7fbde6c975d169eb824d44cd","say":"go is the best"}
```

server-p0 is a baseline version. server-p1 uses fasthttp, server-p2 uses fasthttp + easyjson.

## Benchmark

```sh
ab -p post.txt -n 200000 -c 16 -k http://127.0.0.1:8080/
```

Mac OS X 10.13.6, Intel i7 2.9GHz

| Binary name | RPS   |
|-------------|-------|
| server-p0   | 38800 |
| server-p1   | 66200 |
| server-p2   | 85500 |
