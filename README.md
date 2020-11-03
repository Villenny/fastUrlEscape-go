[![GitHub issues](https://img.shields.io/github/issues/Villenny/fastUrlEscape-go)](https://github.com/Villenny/fastUrlEscape-go/issues)
[![GitHub forks](https://img.shields.io/github/forks/Villenny/fastUrlEscape-go)](https://github.com/Villenny/fastUrlEscape-go/network)
[![GitHub stars](https://img.shields.io/github/stars/Villenny/fastUrlEscape-go)](https://github.com/Villenny/fastUrlEscape-go/stargazers)
[![GitHub license](https://img.shields.io/github/license/Villenny/fastUrlEscape-go)](https://github.com/Villenny/fastUrlEscape-go/blob/master/LICENSE)
![Go](https://github.com/Villenny/fastUrlEscape-go/workflows/Go/badge.svg?branch=main)
![Codecov branch](https://img.shields.io/codecov/c/github/villenny/fastUrlEscape-go/main)
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
	buf = AppendPathEscape(buf[:0], some_string_to_escape))
```

if you need large buffers, probably you want
```
	// used to amortize the cost of the memclear of the buffer.
	// go doesnt provide any way to allocate a byte slice on the stack that isnt cleared to all 0's
	var pool = sync.Pool{
		New: func() interface{} {
			bb := bytes.Buffer{}
			bb.Grow(8192)
			return &bb
		},
	}
	
	func BytesAsString(bs []byte) string {
		// from strings.Builder
		return *(*string)(unsafe.Pointer(&bs))
	}

	import "github.com/villenny/fastUrlEscape-go"

	bb := pool.Get().(*bytes.Buffer)
	buf := bb.Bytes()[:0]
	
	buf = append(buf, "SOME STRING"...)
	buf = strconv.AppendInt(buf, 666, 10)
	buf = AppendQueryEscape(buf, "some string to escape"))

	if bb.Cap() > 8192 {
		bb = pool.New().(*bytes.Buffer)
	}
	pool.Put(bb)

```


## Benchmark

- from my machine
- at [1024]byte and bigger, you're better off with a sync pool than a stack variable (due to the cost of clearing the stack allocation to all zeros)

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
BenchmarkPathEscape-8                            4432632               544 ns/op             160 B/op          2 allocs/op
BenchmarkQueryEscape
BenchmarkQueryEscape-8                           4384117               573 ns/op             160 B/op          2 allocs/op
BenchmarkAppendPathEscape
BenchmarkAppendPathEscape-8                     15910676               154 ns/op               0 B/op          0 allocs/op
BenchmarkAppendQueryEscape
BenchmarkAppendQueryEscape-8                    15205501               158 ns/op               0 B/op          0 allocs/op
BenchmarkAppendQueryEscape_syncPool
BenchmarkAppendQueryEscape_syncPool-8           15302218               156 ns/op               0 B/op          0 allocs/op
PASS
ok      github.com/villenny/fastUrlEscape-go    13.953s

```

## Contact

Ryan Haksi [ryan.haksi@gmail.com]

## License

Available under the BSD [License](/LICENSE). Or any license really, do what you like

