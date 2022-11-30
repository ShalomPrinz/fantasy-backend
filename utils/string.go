package utils

import (
	"strings"
	"unicode"
)

// Returns given string, first letter upper case and the rest lower case
func Capitalize(s string) string {
	var first rune
	for _, c := range s {
		first = unicode.ToUpper(c)
		break
	}
	rest := strings.ToLower(s[1:])
	return string(first) + rest
}
