package repl_test

import (
	"testing"

	"github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/assert"

	"github.com/kamilsk/grafaman/internal/provider"
	. "github.com/kamilsk/grafaman/internal/repl"
)

func TestCompleter(t *testing.T) {
	tests := map[string]struct {
		prefix   string
		metrics  provider.Metrics
		document prompt.Document
		expected []prompt.Suggest
	}{
		"empty char": {
			prefix:  "apps.services.awesome",
			metrics: all,
			document: func() prompt.Document {
				buf := prompt.NewBuffer()
				buf.InsertText("", false, true)
				return *buf.Document()
			}(),
			expected: []prompt.Suggest{
				{Text: "service"},
				{Text: "token_per_url"},
			},
		},
		"first char": {
			prefix:  "apps.services.awesome",
			metrics: all,
			document: func() prompt.Document {
				buf := prompt.NewBuffer()
				buf.InsertText("s", false, true)
				return *buf.Document()
			}(),
			expected: []prompt.Suggest{
				{Text: "service"},
			},
		},
		"first char with wildcard": {
			prefix:  "apps.services.awesome",
			metrics: all,
			document: func() prompt.Document {
				buf := prompt.NewBuffer()
				buf.InsertText("s*", false, true)
				return *buf.Document()
			}(),
			expected: []prompt.Suggest{
				{Text: "service"},
			},
		},
		"first word": {
			prefix:  "apps.services.awesome",
			metrics: all,
			document: func() prompt.Document {
				buf := prompt.NewBuffer()
				buf.InsertText("service", false, true)
				return *buf.Document()
			}(),
			expected: []prompt.Suggest{
				{Text: "service"},
			},
		},
		"first complete word": {
			prefix:  "apps.services.awesome",
			metrics: all,
			document: func() prompt.Document {
				buf := prompt.NewBuffer()
				buf.InsertText("service.", false, true)
				return *buf.Document()
			}(),
			expected: []prompt.Suggest{
				{Text: "api"},
				{Text: "rpc"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			completer := NewMetricsCompleter(test.prefix, test.metrics)
			assert.Equal(t, test.expected, completer(test.document))
		})
	}
}

// helpers

var all = provider.Metrics{
	"apps.services.awesome.service.api.service-x_get_POST.request_time.499.count",
	"apps.services.awesome.service.api.service-x_get_POST.request_time.499.max",
	"apps.services.awesome.service.api.service-x_get_POST.request_time.499.mean",
	"apps.services.awesome.service.api.service-x_get_POST.request_time.499.median",
	"apps.services.awesome.service.api.service-x_get_POST.request_time.499.min",
	"apps.services.awesome.service.api.service-x_get_POST.request_time.499.sum",
	"apps.services.awesome.service.rpc.client.service-y.method.error.count",
	"apps.services.awesome.service.rpc.client.service-y.method.error.max",
	"apps.services.awesome.service.rpc.client.service-y.method.error.mean",
	"apps.services.awesome.service.rpc.client.service-y.method.error.median",
	"apps.services.awesome.service.rpc.client.service-y.method.error.min",
	"apps.services.awesome.service.rpc.client.service-y.method.error.sum",
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.999",
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.99",
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.98",
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.95",
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.75",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.999",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.99",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.98",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.95",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.75",
	"apps.services.awesome.token_per_url.get.service-x",
	"apps.services.awesome.token_per_url.get.service-y",
	"apps.services.awesome.token_per_url.get.service-z",
}
