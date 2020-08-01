package cnf_test

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/model"
)

func TestConfig_FilterQuery(t *testing.T) {
	tests := map[string]struct {
		config   map[string]interface{}
		expected model.Query
	}{
		"empty input": {expected: "*"},
		"with prefix only": {
			config: map[string]interface{}{"metrics": "set"}, expected: "set.*",
		},
		"with filter only": {
			config: map[string]interface{}{"filter": "subset.*"}, expected: "subset.*",
		},
		"full input": {
			config: map[string]interface{}{"metrics": "set", "filter": "subset.*"}, expected: "set.subset.*",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var config Config
			require.NoError(t, mapstructure.Decode(test.config, &config))
			assert.Equal(t, test.expected, config.FilterQuery())
		})
	}
}
