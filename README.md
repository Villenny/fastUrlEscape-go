[![GitHub issues](https://img.shields.io/github/issues/Villenny/fastUrlEscape-go)](https://github.com/Villenny/fastUrlEscape-go/issues)
[![GitHub forks](https://img.shields.io/github/forks/Villenny/fastUrlEscape-go)](https://github.com/Villenny/fastUrlEscape-go/network)
[![GitHub stars](https://img.shields.io/github/stars/Villenny/fastUrlEscape-go)](https://github.com/Villenny/fastUrlEscape-go/stargazers)
[![GitHub license](https://img.shields.io/github/license/Villenny/fastUrlEscape-go)](https://github.com/Villenny/fastUrlEscape-go/blob/master/LICENSE)
![Go](https://github.com/Villenny/fastUrlEscape-go/workflows/Go/badge.svg?branch=master)
![Codecov branch](https://img.shields.io/codecov/c/github/villenny/fastUrlEscape-go/master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Villenny/fastUrlEscape-go)](https://goreportcard.com/report/github.com/Villenny/fastUrlEscape-go)
[![Documentation](https://godoc.org/github.com/Villenny/fastUrlEscape-go?status.svg)](http://godoc.org/github.com/Villenny/fastUrlEscape-go)

# fastUrlEscape-go
- zero allocation url escaping
- 4X faster than net/url


## Install

```
go get -u github.com/Villenny/fastUrlEscape-go
```

## Notable members:
`AppendPathEscape`,
`AppendQueryEscape`,

The expected use case:
- assuming you are doing a lot of path escaping
```
	import "github.com/villenny/fastUrlEscape-go"

	var buf [1024]byte

```


## Benchmark

- from my machine

```
$ ./bench.sh
=== RUN   TestPathEscape
=== RUN   TestPathEscape/early_outs_when_nothing_to_escape
=== RUN   TestPathEscape/matches_url.PathEscape_for_query_type_string
=== RUN   TestPathEscape/allocates_when_buf_is_too_small
=== RUN   TestPathEscape/concats_to_existing_buffer
--- PASS: TestPathEscape (0.00s)
    --- PASS: TestPathEscape/early_outs_when_nothing_to_escape (0.00s)
    --- PASS: TestPathEscape/matches_url.PathEscape_for_query_type_string (0.00s)
    --- PASS: TestPathEscape/allocates_when_buf_is_too_small (0.00s)
    --- PASS: TestPathEscape/concats_to_existing_buffer (0.00s)
=== RUN   TestQueryEscape
--- PASS: TestQueryEscape (0.00s)
goos: windows
goarch: amd64
pkg: github.com/villenny/fastUrlEscape-go
BenchmarkPathEscape
BenchmarkPathEscape-8            4392124               543 ns/op             160 B/op          2 allocs/op
BenchmarkQueryEscape
BenchmarkQueryEscape-8           4376107               563 ns/op             160 B/op          2 allocs/op
BenchmarkAppendPathEscape
BenchmarkAppendPathEscape-8     15805736               155 ns/op               0 B/op          0 allocs/op
BenchmarkAppendQueryEscape
BenchmarkAppendQueryEscape-8    13650548               157 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/villenny/fastUrlEscape-go    11.091s
```

## Contact

Ryan Haksi [ryan.haksi@gmail.com]

## License

Available under the MIT [License](/LICENSE).
