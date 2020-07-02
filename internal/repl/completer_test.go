package repl_test

import (
	"testing"

	"github.com/c-bata/go-prompt"
	"github.com/stretchr/testify/assert"

	"github.com/kamilsk/grafaman/internal/model"
	. "github.com/kamilsk/grafaman/internal/repl"
)

func TestCompleter(t *testing.T) {
	tests := map[string]struct {
		prefix   string
		metrics  model.MetricNames
		document prompt.Document
		expected []prompt.Suggest
	}{
		"empty char": {
			prefix:   "apps.services.awesome",
			metrics:  metrics,
			document: document(""),
			expected: []prompt.Suggest{
				{Text: "service."},
				{Text: "token_per_url."},
			},
		},
		"first char": {
			prefix:   "apps.services.awesome",
			metrics:  metrics,
			document: document("s"),
			expected: []prompt.Suggest{
				{Text: "service."},
			},
		},
		"first char with wildcard": {
			prefix:   "apps.services.awesome",
			metrics:  metrics,
			document: document("s*"),
			expected: []prompt.Suggest{
				{Text: "service."},
			},
		},
		"first word": {
			prefix:   "apps.services.awesome",
			metrics:  metrics,
			document: document("service"),
			expected: []prompt.Suggest{
				{Text: "service."},
			},
		},
		"first segment": {
			prefix:   "apps.services.awesome",
			metrics:  metrics,
			document: document("service."),
			expected: []prompt.Suggest{
				{Text: "service.api."},
				{Text: "service.rpc."},
			},
		},
		"every time ordered": {
			prefix:   "apps.services.awesome",
			metrics:  metrics,
			document: document("service.rpc.client.service-y.method.ok.percentile."),
			expected: []prompt.Suggest{
				{Text: "service.rpc.client.service-y.method.ok.percentile.75"},
				{Text: "service.rpc.client.service-y.method.ok.percentile.95"},
				{Text: "service.rpc.client.service-y.method.ok.percentile.98"},
				{Text: "service.rpc.client.service-y.method.ok.percentile.99"},
				{Text: "service.rpc.client.service-y.method.ok.percentile.999"},
			},
		},
		"nothing found": {
			prefix:   "apps.services.cool",
			metrics:  metrics,
			document: document("service."),
			expected: nil,
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

var metrics = model.MetricNames{
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
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.75",
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.95",
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.98",
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.99",
	"apps.services.awesome.service.rpc.client.service-y.method.error.percentile.999",
	"apps.services.awesome.service.rpc.client.service-y.method.error.sum",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.95",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.75",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.99",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.98",
	"apps.services.awesome.service.rpc.client.service-y.method.ok.percentile.999",
	"apps.services.awesome.token_per_url.get.service-x",
	"apps.services.awesome.token_per_url.get.service-y",
	"apps.services.awesome.token_per_url.get.service-z",
}

func document(text string) prompt.Document {
	buf := prompt.NewBuffer()
	buf.InsertText(text, false, true)
	return *buf.Document()
}
