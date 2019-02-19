package preprocessing

import (
	"io"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestWrappedReader(t *testing.T) {
	g := NewGomegaWithT(t)

	input := `foo:
- bar
- unicode: öäü
- props:
    <<: (( merge replace ))
    b: 3
    c: 4
    d: "abc<<:de"
    list:
    - << : (( merge on key ))
    - key: alice
      age: 25
    - key: bob
      age: 24
`

	expected := `foo:
- bar
- unicode: öäü
- props:
    __: (( merge replace ))
    b: 3
    c: 4
    d: "abc<<:de"
    list:
    - __ : (( merge on key ))
    - key: alice
      age: 25
    - key: bob
      age: 24
`

	for l := 6; l < len(input); l++ {
		for k := -2; k <= 2; k ++ {
			reader := strings.NewReader(input)
			var builder strings.Builder
			wr := NewWrappedReader(reader, l)
			p := make([]byte, l+k)
			eof := false
			for !eof {
				n, err := wr.Read(p)
				if err != nil {
					if err == io.EOF {
						eof = true
					} else {
						t.Errorf("l=%d,k=%d Read failed with error: %v", l, k, err)
						return
					}
				}
				_, err = builder.Write(p[:n])
				if err != nil {
					t.Errorf("l=%d,k=%d Write failed with error: %v", l, k, err)
					return
				}
			}
			actual := builder.String()
			g.Expect(len(actual)).To(Equal(len(expected)), "l=%d,k=%d not as expected. %d != %d", l, k, len(actual), len(expected))
			g.Expect(actual).To(Equal(expected), "l=%d,k=%d not as expected.\nExpected:\n%s\nActual:\n%s", l, k, expected, actual)
		}
	}
}
