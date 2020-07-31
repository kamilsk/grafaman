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

	type payload struct {
		Dashboard dashboard `json:"dashboard,omitempty"`
	}

	type response struct {
		Code int     `json:"code,omitempty"`
		Body payload `json:"body,omitempty"`
	}

	t.Run("success", func(t *testing.T) {
		resp := response{
			Code: http.StatusOK,
			Body: payload{
				Dashboard: dashboard{
					Panels: []panel{
						{
							ID:    1,
							Title: "Panel A",
							Type:  "singlestat",
							Targets: []target{
								{
									Query: "sumSeriesWithWildcards(movingSum(apps.services.*.rpc.*, '1min'), 3, 5)",
								},
							},
						},
						{
							ID:    2,
							Title: "Error rate",
							Type:  "row",
							Panels: []panel{
								{
									ID:    3,
									Title: "Panel B",
									Type:  "graph",
									Targets: []target{
										{
											Query: "aliasByNode(movingSum(apps.services.*.errors.*, '1min'), 3, 6, 5)",
										},
									},
								},
							},
						},
					},
					Templating: templating{
						List: []variable{
							{
								Name:    "env",
								Options: []option{},
								Current: currentOption{
									Text:  "prod",
									Value: "prod",
								},
							},
							{
								Name: "source",
								Options: []option{
									{
										Text:  "All",
										Value: "$__all",
									},
									{
										Text:  "service",
										Value: "service",
									},
								},
								Current: currentOption{
									Text:  "All",
									Value: []interface{}{"$__all"},
								},
							},
						},
					},
				},
			},
		}

		file, err := fs.Create("testdata/success.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(resp))
		require.NoError(t, file.Close())
	})
}
