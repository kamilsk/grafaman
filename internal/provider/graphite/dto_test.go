package graphite

import (
	"encoding/json"
	"flag"
	"net/http"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update golden files")

func TestDumpStubs(t *testing.T) {
	fs := afero.NewMemMapFs()
	if *update {
		fs = afero.NewOsFs()
	}

	type response struct {
		Code int   `json:"code,omitempty"`
		Body []dto `json:"body,omitempty"`
	}

	t.Run("success", func(t *testing.T) {
		file, err := fs.Create("testdata/success.1.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(response{
			Code: http.StatusOK,
			Body: []dto{
				{
					ID:   "apps.services.awesome-service",
					Text: "awesome-service",
					Leaf: 0,
				},
			},
		}))
		require.NoError(t, file.Close())

		file, err = fs.Create("testdata/success.2.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(response{
			Code: http.StatusOK,
			Body: []dto{
				{
					ID:   "apps.services.awesome-service.metric",
					Text: "metric",
					Leaf: 0,
				},
			},
		}))
		require.NoError(t, file.Close())

		file, err = fs.Create("testdata/success.3.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(response{
			Code: http.StatusOK,
			Body: []dto{
				{
					ID:   "apps.services.awesome-service.metric.a",
					Text: "a",
					Leaf: 1,
				},
				{
					ID:   "apps.services.awesome-service.metric.b",
					Text: "b",
					Leaf: 1,
				},
				{
					ID:   "apps.services.awesome-service.metric.c",
					Text: "c",
					Leaf: 1,
				},
			},
		}))
		require.NoError(t, file.Close())
	})

	t.Run("parallel", func(t *testing.T) {
		file, err := fs.Create("testdata/parallel.1.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(response{
			Code: http.StatusOK,
			Body: []dto{
				{
					ID:   "apps.services.awesome-service",
					Text: "awesome-service",
					Leaf: 0,
				},
			},
		}))
		require.NoError(t, file.Close())

		file, err = fs.Create("testdata/parallel.2.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(response{
			Code: http.StatusOK,
			Body: []dto{
				{
					ID:   "apps.services.awesome-service.external",
					Text: "external",
					Leaf: 0,
				},
				{
					ID:   "apps.services.awesome-service.internal",
					Text: "internal",
					Leaf: 0,
				},
			},
		}))
		require.NoError(t, file.Close())

		file, err = fs.Create("testdata/parallel.3-1.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(response{
			Code: http.StatusOK,
			Body: []dto{
				{
					ID:   "apps.services.awesome-service.external.a",
					Text: "a",
					Leaf: 1,
				},
				{
					ID:   "apps.services.awesome-service.external.b",
					Text: "b",
					Leaf: 1,
				},
				{
					ID:   "apps.services.awesome-service.external.c",
					Text: "c",
					Leaf: 1,
				},
			},
		}))
		require.NoError(t, file.Close())

		file, err = fs.Create("testdata/parallel.3-2.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(response{
			Code: http.StatusOK,
			Body: []dto{
				{
					ID:   "apps.services.awesome-service.internal.a",
					Text: "a",
					Leaf: 1,
				},
				{
					ID:   "apps.services.awesome-service.internal.b",
					Text: "b",
					Leaf: 1,
				},
				{
					ID:   "apps.services.awesome-service.internal.c",
					Text: "c",
					Leaf: 1,
				},
			},
		}))
		require.NoError(t, file.Close())
	})
}
