package grafana

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kamilsk/grafaman/internal/model"
)

func TestConvertTargets(t *testing.T) {
	tests := map[string]struct {
		targets  []target
		expected []model.Query
	}{
		"issue#7": {
			targets: []target{
				{
					Query: "",
				},
			},
			expected: []model.Query{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, convertTargets(test.targets))
		})
	}
}
