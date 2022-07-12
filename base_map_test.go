package ascii_test

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/tmaxmax/ascii-converter"
)

func TestBaseMap(t *testing.T) {
	b := ascii.Base2 | ascii.Base4 | ascii.Base5 | ascii.Base8 | ascii.Base10 | ascii.Base16
	rep := b.String()

	if rep != "binary,base4,base5,octal,decimal,hex" {
		t.Errorf("Invalid base representation %q", rep)
	}
}

func TestGetValidBases(t *testing.T) {
	type test struct {
		input    string
		expected ascii.Base
	}

	tests := [...]test{
		{"0110", ascii.AllValidBases},
		{"3924820", ascii.Base10 | ascii.Base11 | ascii.Base12 | ascii.Base13 | ascii.Base14 | ascii.Base15 | ascii.Base16},
	}

	for _, testCase := range tests {
		b := ascii.GetValidBases(testCase.input)
		if b != testCase.expected {
			t.Errorf("For input %q: got %q, expected %q", testCase.input, b, testCase.expected)
		}
	}
}

type randGen struct {
	source *rand.Rand
	buf    int64
}

const mask = 0b01111111

func (r *randGen) Next() byte {
	if r.buf == 0 {
		r.buf = r.source.Int63()
	}

	for {
		for i := 0; i < 9; i++ {
			b := byte(r.buf & mask)
			r.buf >>= 7

			if ascii.IsASCIINumber(b) {
				return b
			}
		}

		r.buf = r.source.Int63()
	}
}

func (r *randGen) Fill(buf []byte) {
	for i := range buf {
		buf[i] = r.Next()
	}
}

func (r *randGen) FillString(size int) string {
	var sb strings.Builder
	sb.Grow(size)

	for i := 0; i < size; i++ {
		sb.WriteByte(r.Next())
	}

	return sb.String()
}

func newRandGen() *randGen {
	return &randGen{
		source: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func BenchmarkConversion(b *testing.B) {
	b.StopTimer()
	gen := newRandGen()
	nums := make([]byte, b.N)
	gen.Fill(nums)
	b.StartTimer()
	res := uint8(0)
	for i := 0; i < b.N; i++ {
		res, _ = ascii.GetNumericValue(nums[i])
	}
	_ = res
}

func BenchmarkValid(b *testing.B) {
	b.StopTimer()
	const size = 16
	gen := newRandGen()
	strs := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		strs[i] = gen.FillString(size)
	}
	b.StartTimer()
	res := ascii.NoValidBases
	for i := 0; i < b.N; i++ {
		res = ascii.GetValidBases(strs[i])
	}
	_ = res
}
