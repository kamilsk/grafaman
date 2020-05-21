package grafana

import (
	"testing"

	"github.com/stretchr/testify/assert"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

func TestConvertTargets(t *testing.T) {
	tests := map[string]struct {
		targets  []target
		expected []entity.Query
	}{
		"issue#7": {
			targets: []target{
				{
					Query: "",
				},
			},
			expected: []entity.Query{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, convertTargets(test.targets))
		})
	}
}
