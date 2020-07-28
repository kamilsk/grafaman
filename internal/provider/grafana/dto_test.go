package grafana

import (
	"encoding/json"
	"flag"
	"net/http"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kamilsk/grafaman/internal/model"
)

var update = flag.Bool("update", false, "update golden files")

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

func TestDumpStubs(t *testing.T) {
	fs := afero.NewMemMapFs()
	if *update {
		fs = afero.NewOsFs()
	}

	type response struct {
		Code int       `json:"code,omitempty"`
		Body dashboard `json:"body,omitempty"`
	}

	t.Run("success", func(t *testing.T) {
		resp := response{
			Code: http.StatusOK,
			Body: dashboard{},
		}

		file, err := fs.Create("testdata/success.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(resp))
		require.NoError(t, file.Close())
	})
}
