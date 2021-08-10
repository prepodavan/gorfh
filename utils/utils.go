package utils

import (
	"unicode"
	"unicode/utf8"
)

func IsName(p []byte) bool {
	i := 0
	var r rune
	var s int
	for len(p) > 0 {
		r, s = utf8.DecodeRune(p)
		if r == utf8.RuneError {
			return false
		}
		cond := unicode.IsLetter(r)
		if i != 0 {
			cond = cond || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.'
		}
		if !cond {
			return false
		}
		i++
		p = p[s:]
	}
	return r != '.'
}
