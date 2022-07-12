package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tmaxmax/ascii-converter"
)

var bases = [...]ascii.Base{
	ascii.Base10,
	ascii.Base2,
	ascii.Base16,
	ascii.Base8,
	ascii.Base3,
	ascii.Base4,
	ascii.Base5,
	ascii.Base6,
	ascii.Base7,
	ascii.Base9,
	ascii.Base11,
	ascii.Base12,
	ascii.Base13,
	ascii.Base14,
	ascii.Base15,
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	s.Split(bufio.ScanWords)

	validBases := ascii.AllValidBases
	tokens := []string(nil)

	for s.Scan() {
		token := s.Text()
		tokens = append(tokens, token)
		validBases &= ascii.GetValidBases(token)
	}

	fmt.Println("first pass:", validBases)
	if validBases == ascii.NoValidBases {
		return
	}

	convertors := validBases.Convertors()
	converted := map[ascii.Base]*strings.Builder{}
	invalidBases := map[ascii.Base]struct{}{}

	for _, token := range tokens {
		if len(invalidBases) == len(convertors) {
			break
		}

		for _, convertor := range convertors {
			base := convertor.Base()
			if _, ok := invalidBases[base]; ok {
				continue
			}

			v, ok := convertor.Convert(token)
			if ok {
				sb := converted[base]
				if sb == nil {
					sb = new(strings.Builder)
					converted[base] = sb
				}
				sb.WriteByte(v)
				continue
			}

			delete(converted, base)
			invalidBases[base] = struct{}{}
			validBases &= ^base
		}
	}

	fmt.Println("second pass:", validBases)

	for _, base := range bases {
		if text := converted[base]; text != nil {
			fmt.Printf("%s: %s\n", base, text)
		}
	}
}
