package model_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/kamilsk/grafaman/internal/model"
)

func TestMetric_Valid(t *testing.T) {
	tests := map[string]struct {
		metric   Metric
		expected bool
	}{
		"empty":      {"", true},
		"root":       {"apps", true},
		"pod name":   {"apps.services.awesome-service.go.pod-5dbdcd5dbb-6z58f.threadsv", true},
		"jaeger":     {"apps.services.awesome-service.jaeger.finished_spans_sampled_n", true},
		"percentile": {"apps.services.awesome-service.rpc.client.success.ok.percentile.999", true},
		"invalid":    {"$env.apps.$space", false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.metric.Valid())
		})
	}
}

func TestMetrics(t *testing.T) {
	raw := []string{"b", "c", "a"}

	var metrics Metrics
	require.NotPanics(t, func() { assert.Len(t, (*Metrics)(nil).Convert(raw), len(raw)) })
	require.NotPanics(t, func() { assert.Len(t, metrics.Convert(raw), len(raw)) })

	assert.False(t, sort.StringsAreSorted(raw))
	assert.Len(t, metrics.Sort(), len(raw))
	assert.True(t, sort.StringsAreSorted(raw))
}

func TestMetrics_Exclude(t *testing.T) {
	tests := map[string]struct {
		metrics  Metrics
		exclude  []Matcher
		expected Metrics
	}{
		"nil metrics": {
			metrics:  nil,
			exclude:  []Matcher{Query("a.*").MustCompile()},
			expected: nil,
		},
		"empty metrics": {
			metrics:  Metrics{},
			exclude:  []Matcher{Query("a.*").MustCompile()},
			expected: Metrics{},
		},
		"nothing to exclude": {
			metrics: Metrics{
				"metric.a",
				"metric.b",
				"metric.c",
			},
			exclude: []Matcher{Query("a.*").MustCompile()},
			expected: Metrics{
				"metric.a",
				"metric.b",
				"metric.c",
			},
		},
		"with exclusion": {
			metrics: Metrics{
				"metric.a",
				"metric.b",
				"metric.c",
				"metric.d",
			},
			exclude: []Matcher{
				Query("*.b").MustCompile(),
				Query("*.d").MustCompile(),
			},
			expected: Metrics{
				"metric.a",
				"metric.c",
			},
		},
		"all excluded": {
			metrics: Metrics{
				"metric.a",
				"metric.b",
				"metric.c",
				"metric.d",
			},
			exclude:  []Matcher{Query("metric.*").MustCompile()},
			expected: Metrics{},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.metrics.Exclude(test.exclude...))
		})
	}
}

func TestMetrics_Filter(t *testing.T) {
	tests := map[string]struct {
		metrics  Metrics
		filter   []Matcher
		expected Metrics
	}{
		"nil metrics": {
			metrics:  nil,
			filter:   []Matcher{Query("a.*").MustCompile()},
			expected: nil,
		},
		"empty metrics": {
			metrics:  Metrics{},
			filter:   []Matcher{Query("a.*").MustCompile()},
			expected: Metrics{},
		},
		"nothing to match": {
			metrics: Metrics{
				"metric.a",
				"metric.b",
				"metric.c",
			},
			filter:   []Matcher{Query("metric.*.*").MustCompile()},
			expected: Metrics{},
		},
		"partial match": {
			metrics: Metrics{
				"metric.a",
				"metric.b",
				"metric.c",
				"metric.d",
			},
			filter: []Matcher{
				Query("*.b").MustCompile(),
				Query("*.d").MustCompile(),
			},
			expected: Metrics{
				"metric.b",
				"metric.d",
			},
		},
		"all match": {
			metrics: Metrics{
				"metric.a",
				"metric.b",
				"metric.c",
			},
			filter: []Matcher{Query("metric.*").MustCompile()},
			expected: Metrics{
				"metric.a",
				"metric.b",
				"metric.c",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.metrics.Filter(test.filter...))
		})
	}
}
