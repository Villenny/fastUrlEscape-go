package fastUrlEscape

import (
	"net/url"
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
	})
	t.Run("allocates when buf is too small", func(t *testing.T) {
		var buf [16]byte
		s := BytesAsString(AppendPathEscape(buf[:0], ESCAPE_ME))
		assert.Equal(t, url.PathEscape(ESCAPE_ME), s)
	})
	t.Run("concats to existing buffer", func(t *testing.T) {
		var buf [1024]byte
		slice := append(buf[:0], "hello "...)
		s := BytesAsString(AppendPathEscape(slice, ESCAPE_ME))
		assert.Equal(t, "hello "+url.PathEscape(ESCAPE_ME), s)
	})
}

func TestQueryEscape(t *testing.T) {
	var buf [1024]byte
	s := BytesAsString(AppendQueryEscape(buf[:0], ESCAPE_ME))
	assert.Equal(t, url.QueryEscape(ESCAPE_ME), s)
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
pkg: bitbucket.org/kidozteam/bidder-server/pkg/helper/fastUrlEscape
BenchmarkQueryEscape
BenchmarkQueryEscape-8   	 1909773	       685 ns/op	     160 B/op	       2 allocs/op
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
pkg: bitbucket.org/kidozteam/bidder-server/pkg/helper/fastUrlEscape
BenchmarkAppendPathEscape
BenchmarkAppendPathEscape-8   	 7673206	       152 ns/op	       0 B/op	       0 allocs/op
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
pkg: bitbucket.org/kidozteam/bidder-server/pkg/helper/fastUrlEscape
BenchmarkAppendQueryEscape
BenchmarkAppendQueryEscape-8   	 6756363	       148 ns/op	       0 B/op	       0 allocs/op
PASS
*/
func BenchmarkAppendQueryEscape(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf [1024]byte
		s := BytesAsString(AppendQueryEscape(buf[:0], ESCAPE_ME))
		if s == "" {
			panic("WTF")
		}
	}
}
