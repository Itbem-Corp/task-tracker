// Package stringutil provides utility functions for string manipulation.
package stringutil

import (
	"strings"
	"unicode"
)

// Reverse returns the input string reversed, correctly handling
// multi-byte UTF-8 characters by operating on runes.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// IsPalindrome reports whether the input string reads the same forwards
// and backwards. The comparison is case-insensitive and ignores
// non-alphanumeric characters (spaces, punctuation, etc.).
func IsPalindrome(s string) bool {
	s = strings.ToLower(s)
	runes := make([]rune, 0, len(s))
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			runes = append(runes, r)
		}
	}

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		if runes[i] != runes[j] {
			return false
		}
	}
	return true
}
