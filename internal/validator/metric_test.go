package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/kamilsk/grafaman/internal/validator"
)

func TestMetric(t *testing.T) {
	matcher := Metric()

	tests := map[string]struct {
		metric   string
		expected bool
	}{
		"root":       {"apps", true},
		"pod name":   {"apps.services.awesome-service.go.pod-5dbdcd5dbb-6z58f.threadsv", true},
		"jaeger":     {"apps.services.awesome-service.jaeger.finished_spans_sampled_n", true},
		"percentile": {"apps.services.awesome-service.rpc.client.success.ok.percentile.999", true},
		"invalid":    {"$env.apps.$space", false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, matcher(test.metric))
		})
	}
}
