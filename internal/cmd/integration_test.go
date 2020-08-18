package cmd_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

var (
	grafana, graphite *httptest.Server
	buffer            *bytes.Buffer
	root              *cobra.Command
)

var (
	_ = BeforeSuite(func() {
		grafana = httptest.NewServer(http.NewServeMux())
		graphite = httptest.NewServer(http.NewServeMux())
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
