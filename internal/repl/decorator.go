package repl

import "strings"

func Prefix(prefix string, fn func(string)) func(string) {
	if prefix == "" {
		return fn
	}
	return func(pattern string) {
		fn(prefix + "." + strings.TrimSpace(pattern))
	}
}
