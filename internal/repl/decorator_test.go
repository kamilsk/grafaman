package repl_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/kamilsk/grafaman/internal/repl"
)

func TestPrefix(t *testing.T) {
	inputs := []string{
		"input1",
		"input2",
		"input3",
		" \n\r\t ",
	}

	tests := map[string]struct {
		prefix   string
		expected []string
	}{
		"empty prefix": {
			expected: inputs,
		},
		"non-empty prefix": {
			prefix: "prefix",
			expected: []string{
				"prefix.input1",
				"prefix.input2",
				"prefix.input3",
				"prefix.",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for i, input := range inputs {
				Prefix(test.prefix, func(input string) {
					assert.Equal(t, test.expected[i], input)
				})(input)
			}
		})
	}
}
