package fastUrlEscape

import (
	"bytes"
	"net/url"
	"sync"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func BytesAsString(bs []byte) string {
	// from strings.Builder
	return *(*string)(unsafe.Pointer(&bs))
}

const DONT_ESCAPE_ME = "there-is-nothing-to-escape-here"
const ESCAPE_ME = "https://host.domain.com/some/url/path?arg1=one&arg2=two"
const ESCAPE_ME_QUERY_ARG = "one & two ? three / four"
const ESCAPE_ME_QUERY_ONLY_SPACE = "one two three four"
const ESCAPE_ME_ALL_THE_SPECIAL_CHARS = " ?&=#+%!<>#\"{}|\\^[]`☺\t:/@$'()*,;"

func TestPathEscape(t *testing.T) {
	t.Run("early outs when nothing to escape", func(t *testing.T) {
		var buf [1024]byte
		s := BytesAsString(AppendPathEscape(buf[:0], DONT_ESCAPE_ME))
		assert.Equal(t, url.PathEscape(DONT_ESCAPE_ME), s)
	})
	t.Run("matches url.PathEscape for query type string", func(t *testing.T) {
		var buf [1024]byte
		s := BytesAsString(AppendPathEscape(buf[:0], ESCAPE_ME))
		assert.Equal(t, url.PathEscape(ESCAPE_ME), s)

		s = BytesAsString(AppendPathEscape(buf[:0], ESCAPE_ME_QUERY_ARG))
		assert.Equal(t, url.PathEscape(ESCAPE_ME_QUERY_ARG), s)
	})
	t.Run("allocates when buf is too small", func(t *testing.T) {
		var buf [16]byte
		s := BytesAsString(AppendPathEscape(buf[:0], ESCAPE_ME))
		assert.Equal(t, url.PathEscape(ESCAPE_ME), s)
	})
	t.Run("allocates when passed a zero size slice?!?", func(t *testing.T) {
		var buf []byte
		s := BytesAsString(AppendPathEscape(buf, ESCAPE_ME))
		assert.Equal(t, url.PathEscape(ESCAPE_ME), s)
	})
	t.Run("concats to existing buffer", func(t *testing.T) {
		var buf [1024]byte
		slice := append(buf[:0], "hello "...)
		s := BytesAsString(AppendPathEscape(slice, ESCAPE_ME))
		assert.Equal(t, "hello "+url.PathEscape(ESCAPE_ME), s)
	})
	t.Run("all the special chars", func(t *testing.T) {
		var buf [1024]byte
		slice := buf[:0]
		s := BytesAsString(AppendPathEscape(slice, ESCAPE_ME_ALL_THE_SPECIAL_CHARS))
		u := url.PathEscape(ESCAPE_ME_ALL_THE_SPECIAL_CHARS)
		assert.Equal(t, u, s)
	})
}

func TestQueryEscape(t *testing.T) {
	t.Run("does a mix of substitution", func(t *testing.T) {
		var buf [1024]byte
		s := BytesAsString(AppendQueryEscape(buf[:0], ESCAPE_ME_QUERY_ARG))
		assert.Equal(t, url.QueryEscape(ESCAPE_ME_QUERY_ARG), s)
	})
	t.Run("handles special case of all space substitution", func(t *testing.T) {
		var buf [1024]byte
		s := BytesAsString(AppendQueryEscape(buf[:0], ESCAPE_ME_QUERY_ONLY_SPACE))
		assert.Equal(t, url.QueryEscape(ESCAPE_ME_QUERY_ONLY_SPACE), s)
	})
	t.Run("all the special chars", func(t *testing.T) {
		var buf [1024]byte
		slice := buf[:0]
		s := BytesAsString(AppendQueryEscape(slice, ESCAPE_ME_ALL_THE_SPECIAL_CHARS))
		assert.Equal(t, url.QueryEscape(ESCAPE_ME_ALL_THE_SPECIAL_CHARS), s)
	})
}

/*
pkg: bitbucket.org/kidozteam/bidder-server/pkg/helper/fastUrlEscape
BenchmarkPathEscape
BenchmarkPathEscape-8   	 1761361	       665 ns/op	     160 B/op	       2 allocs/op
PASS
*/
func BenchmarkPathEscape(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := url.PathEscape(ESCAPE_ME)
		if s == "" {
			panic("WTF")
		}
	}
}

/*
pkg: github.com/villenny/fastUrlEscape-go
BenchmarkQueryEscape
BenchmarkQueryEscape-8
 2100081	       652 ns/op	     160 B/op	       2 allocs/op
PASS
*/
func BenchmarkQueryEscape(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := url.QueryEscape(ESCAPE_ME)
		if s == "" {
			panic("WTF")
		}
	}
}

/*
pkg: github.com/villenny/fastUrlEscape-go
BenchmarkAppendPathEscape
BenchmarkAppendPathEscape-8
 4017529	       291 ns/op	       0 B/op	       0 allocs/op
PASS
*/
func BenchmarkAppendPathEscape(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf [1024]byte
		s := BytesAsString(AppendPathEscape(buf[:0], ESCAPE_ME))
		if s == "" {
			panic("WTF")
		}
	}
}

/*
pkg: github.com/villenny/fastUrlEscape-go
BenchmarkAppendQueryEscape_Bytearray
BenchmarkAppendQueryEscape_Bytearray-8
 4305537	       292 ns/op	       0 B/op	       0 allocs/op
PASS
*/
func BenchmarkAppendQueryEscape_Bytearray(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// for bufs <= 1024 you can beat a sync.pool
		var buf [1024]byte
		s := BytesAsString(AppendQueryEscape(buf[:0], ESCAPE_ME))
		if s == "" {
			panic("WTF")
		}
	}
}

/*
pkg: github.com/villenny/fastUrlEscape-go
BenchmarkAppendQueryEscape_Make
BenchmarkAppendQueryEscape_Make-8
 4274898	       295 ns/op	       0 B/op	       0 allocs/op
PASS
*/
func BenchmarkAppendQueryEscape_Make(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// for bufs <= 1024 you can beat a sync.pool
		buf := make([]byte, 1024)[:0]
		s := BytesAsString(AppendQueryEscape(buf, ESCAPE_ME))
		if s == "" {
			panic("WTF")
		}
	}
}

// used to amortize the cost of the memclear of the buffer.
// go doesnt provide any way to allocate a byte slice on the stack that isnt cleared to all 0's
var pool = sync.Pool{
	New: func() interface{} {
		bb := bytes.Buffer{}
		bb.Grow(8192)
		return &bb
	},
}

/*
pkg: github.com/villenny/fastUrlEscape-go
BenchmarkAppendQueryEscape_syncPool
BenchmarkAppendQueryEscape_syncPool-8
 4229742	       300 ns/op	       0 B/op	       0 allocs/op
PASS
*/
func BenchmarkAppendQueryEscape_syncPool(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bb := pool.Get().(*bytes.Buffer)

		s := BytesAsString(AppendQueryEscape(bb.Bytes()[:0], ESCAPE_ME))
		if s == "" {
			panic("WTF")
		}

		if bb.Cap() > 8192 {
			bb = pool.New().(*bytes.Buffer)
		}
		pool.Put(bb)
	}
}
