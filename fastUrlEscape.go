package fastUrlEscape

func makePathLUT() []bool {
	var lut [256]bool
	for i := 0; i < 256; i += 1 {
		lut[i] = shouldPathEscape(byte(i))
	}
	return lut[0:]
}

var pathLUT []bool = makePathLUT()

func makeQueryLUT() []bool {
	var lut [256]bool
	for i := 0; i < 256; i += 1 {
		lut[i] = shouldQueryEscape(byte(i))
	}
	return lut[0:]
}

var queryLUT []bool = makeQueryLUT()

func shouldPathEscape(c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '_', '.', '~', ':', '@', '&', '=', '+', '$':
		return false
	}
	return true
}

func shouldQueryEscape(c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '_', '.', '~', ' ':
		return false
	}
	return true
}

func AppendQueryEscape(buf []byte, s string) []byte {
	return appendEscape(buf, s, queryLUT)
}
func AppendPathEscape(buf []byte, s string) []byte {
	return appendEscape(buf, s, pathLUT)
}

func appendEscape(buf []byte, s string, lut []bool) []byte {
	hexCount := 0
	hasSpaces := false
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			hasSpaces = true
		}
		if lut[s[i]] {
			hexCount++
		}
	}

	if !hasSpaces && hexCount == 0 {
		return append(buf, s...)
	}

	cb := cap(buf)
	lb := len(buf)
	ls := len(s)

	/*
		Honestly this is probably the best solution to being passed a zero size slice, I have no idea why you would do this
		It would guarantee an allocation

			if cb == 0 {
				panic("passed a zero size buffer, use like var buf [512]byte, appendEscape(buf[:0], someQueryArg)")
			}
	*/

	spaceRequired := ls + 2*hexCount
	bufRequired := spaceRequired + lb
	if cb-lb < spaceRequired {
		for cb <<= 1; cb < bufRequired; cb <<= 1 {
			if cb == 0 {
				cb = 64
			}
		}
		buf2 := make([]byte, cb)
		copy(buf2, buf)
		buf = buf2
	}

	newBufLen := lb + spaceRequired
	buf = buf[0:newBufLen]
	t := buf[lb:newBufLen]
	j := 0
	for i := 0; i < ls; i++ {
		c := s[i]
		if lut[c] {
			t[j+2] = "0123456789ABCDEF"[c&15]
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j] = '%'
			j += 3

		} else if c == ' ' {
			// if space isnt converted to hex, convert it to plus symbol
			t[j] = '+'
			j++

		} else {
			t[j] = c
			j++
		}
	}
	return buf
}
