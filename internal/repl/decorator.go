package repl

import "strings"

// Prefix wraps an input string by the prefix.
func Prefix(prefix string, fn func(string)) func(string) {
	if prefix == "" {
		return fn
	}
	return func(input string) {
		fn(prefix + "." + strings.TrimSpace(input))
	}
}
