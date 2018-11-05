# OC contest in Golang

# Synopsis

```sh
go install ./...
server-p0 -h 127.0.0.1 -p 8080 -n 2
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
