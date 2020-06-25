package repl

import "strings"

func Prefix(prefix string, fn func(string)) func(string) {
	return func(pattern string) {
		pattern = strings.TrimSpace(pattern)
		if pattern != "" {
			pattern = prefix + "." + pattern
		}
		fn(pattern)
	}
}
