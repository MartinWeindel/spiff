package preprocessing

import (
	"fmt"
	"io"
)

type State int

const (
	StartingWhitespaceOrDash State = 0
	MergeArrow1              State = 1
	MergeArrow2              State = 2
	Other                    State = 3
)

type wrappedreader struct {
	inner io.Reader
	state State
	buf   []byte
	pos   int
	len   int
}

func NewWrappedReader(inner io.Reader, bufferSize int) io.Reader {
	buf := make([]byte, bufferSize)
	return &wrappedreader{inner, StartingWhitespaceOrDash, buf, 0, 0}
}

func (wr *wrappedreader) Read(p []byte) (n int, err error) {
	if len(p) < 4 {
		panic(fmt.Sprintf("wrapperreader does not work with buffersize %d, (min size 4)", len(p)))
	}
	for n < len(p) {
		if wr.pos == wr.len {
			if err == io.EOF {
				return n, err
			}
			wr.len, err = wr.inner.Read(wr.buf)
			wr.pos = 0
			if err != nil && err != io.EOF {
				return n, err
			}
			if wr.len == 0 {
				continue
			}
		}
		b := wr.buf[wr.pos]
		wr.pos++
		p[n] = b
		n++
		if b&0x80 != 0 {
			// pass any unicode rune bytes
			wr.state = Other
		} else if b == '\n' || b == '\r' {
			wr.state = StartingWhitespaceOrDash
		} else if wr.state == Other {
			// no change
		} else if wr.state == StartingWhitespaceOrDash {
			switch b {
			case ' ', '\t', '-':
				// no change
			case '<':
				if n+2 >= len(p) {
					wr.pos--
					n--
					return n, err
				}
				wr.state = MergeArrow1
			default:
				wr.state = Other
			}
		} else if wr.state == MergeArrow1 {
			switch b {
			case '<':
				wr.state = MergeArrow2
			default:
				wr.state = Other
			}
		} else if wr.state == MergeArrow2 {
			switch b {
			case ':', ' ', '\t':
				// replace `<<:` with `__:`
				p[n-3] = '_'
				p[n-2] = '_'
				wr.state = Other
			default:
				wr.state = Other
			}
		}
	}
	return n, nil
}
