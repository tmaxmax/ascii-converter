package ascii

import (
	"strconv"
	"strings"
)

type Base uint16

const (
	Base2 Base = 1 << (iota + 1)
	Base3
	Base4
	Base5
	Base6
	Base7
	Base8
	Base9
	Base10
	Base11
	Base12
	Base13
	Base14
	Base15
	Base16
)

var AllBases = [...]Base{Base2, Base3, Base4, Base5, Base6, Base7, Base8, Base9, Base10, Base11, Base12, Base13, Base14, Base15, Base16}

const (
	AllValidBases Base = Base2 | Base3 | Base4 | Base5 | Base6 | Base7 | Base8 | Base9 | Base10 | Base11 | Base12 | Base13 | Base14 | Base15 | Base16
	NoValidBases  Base = 0
)

func (b Base) String() string {
	var sb strings.Builder
	var index = 1
	var rep string

	for b >>= 1; b != 0; b >>= 1 {
		index++

		if b&1 == 0 {
			continue
		}
		if sb.Len() > 0 {
			sb.WriteByte(',')
		}

		switch index {
		case 2:
			rep = "binary"
		case 8:
			rep = "octal"
		case 10:
			rep = "decimal"
		case 16:
			rep = "hex"
		default:
			rep = "base" + strconv.Itoa(index)
		}

		sb.WriteString(rep)
	}

	if sb.Len() == 0 {
		return "invalid"
	}

	return sb.String()
}

func (b Base) Test(other Base) bool {
	return b&other == other
}

func IsASCIINumber(b byte) bool {
	return ('0' <= b && b <= '9') || ('A' <= b && b <= 'F') || ('a' <= b && b <= 'f')
}

const (
	digitsMask   = 0b00110000
	upperHexMask = 0b01000000
	lowerHexMask = 0b01100000
)

func GetNumericValue(b byte) (uint8, bool) {
	switch {
	case '0' <= b && b <= '9':
		return b ^ digitsMask, true
	case 'A' <= b && b <= 'F':
		return b ^ upperHexMask + 9, true
	case 'a' <= b && b <= 'f':
		return b ^ lowerHexMask + 9, true
	default:
		return 0, false
	}
}

func GetValidBases(input string) Base {
	base := AllValidBases

	for i := 0; i < len(input); i++ {
		num, ok := GetNumericValue(input[i])
		if !ok {
			return NoValidBases
		}

		for j, maxBaseIndex := 0, int(num)-2; j <= maxBaseIndex; j++ {
			base &= ^AllBases[j]
		}
	}

	return base
}

type Convertor struct {
	multiplier int
	base       Base
}

var basesMultipliers = [...]int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

func (b Base) Convertors() []Convertor {
	var convs []Convertor

	for i, ab := range AllBases {
		if !b.Test(ab) {
			continue
		}

		convs = append(convs, Convertor{
			multiplier: basesMultipliers[i],
			base:       ab,
		})
	}

	return convs
}

func (c Convertor) Base() Base {
	return c.base
}

func (c Convertor) Convert(input string) (byte, bool) {
	var res int

	for i := 0; i < len(input); i++ {
		v, _ := GetNumericValue(input[i])
		res *= c.multiplier
		res += int(v)
	}

	return byte(res), 0 <= res && res <= 127
}
