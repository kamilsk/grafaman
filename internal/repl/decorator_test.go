package repl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/kamilsk/grafaman/internal/repl"
)

func TestPrefix(t *testing.T) {
	patterns := []string{
		"pattern1",
		"pattern2",
		"pattern3",
		" \n\r\t ",
	}

	tests := map[string]struct {
		prefix   string
		expected []string
	}{
		"empty prefix": {
			expected: patterns,
		},
		"non-empty prefix": {
			prefix: "prefix",
			expected: []string{
				"prefix.pattern1",
				"prefix.pattern2",
				"prefix.pattern3",
				"prefix.",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for i, pattern := range patterns {
				Prefix(test.prefix, func(pattern string) {
					assert.Equal(t, test.expected[i], pattern)
				})(pattern)
			}
		})
	}
}
