package cmd_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"go.octolab.org/unsafe"
)

var (
	grafana, graphite *httptest.Server
	buffer            *bytes.Buffer
	root              *cobra.Command
)

var (
	_ = BeforeSuite(func() {
		grafanaAPI := http.NewServeMux()
		grafana = httptest.NewUnstartedServer(grafanaAPI)
		grafana.Start()

		graphiteAPI := http.NewServeMux()
		graphiteAPI.HandleFunc("/metrics/find", func(rw http.ResponseWriter, req *http.Request) {
			if err := req.ParseForm(); err != nil {
				http.Error(rw, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			name := filepath.Join("testdata", "metrics", req.FormValue("query")) + ".json"
			f, err := os.Open(name)
			if errors.Is(err, os.ErrNotExist) {
				http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			if err != nil {
				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			unsafe.DoSilent(io.Copy(rw, f))
			unsafe.Ignore(f.Close())
		})
		graphite = httptest.NewUnstartedServer(graphiteAPI)
		graphite.Start()

		buffer = bytes.NewBuffer(make([]byte, 0, 1024))
	})

	_ = AfterSuite(func() {
		grafana.Close()
		graphite.Close()
	})
)

func TestComponents(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Integration Suite")
}
