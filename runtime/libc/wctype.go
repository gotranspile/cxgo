package libc

import (
	"unicode"
)

type WChar = uint16

func IsAlpha(c rune) bool {
	if c < 0 {
		return false
	}
	return unicode.IsLetter(c)
}

func IsAlnum(c rune) bool {
	if c < 0 {
		return false
	}
	return unicode.IsLetter(c) || unicode.IsNumber(c)
}
