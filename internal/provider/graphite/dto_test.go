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
		resp := response{
			Code: http.StatusOK,
			Body: []dto{},
		}

		file, err := fs.Create("testdata/success.json")
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(resp))
		require.NoError(t, file.Close())
	})
}
