package filter

import (
	"strings"

	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/kamilsk/grafaman/internal/provider"
)

func Exclude(metrics provider.Metrics, pattern string) (provider.Metrics, error) {
	if len(metrics) == 0 {
		return metrics, nil
	}

	matcher, err := glob.Compile(pattern)
	if err != nil {
		return nil, errors.Wrapf(err, "compile pattern %q", pattern)
	}

	filtered := make(provider.Metrics, 0, len(metrics))
	for _, metric := range metrics {
		if !matcher.Match(string(metric)) {
			filtered = append(filtered, metric)
		}
	}
	return filtered, nil
}

func Filter(metrics provider.Metrics, pattern, prefix string) (provider.Metrics, error) {
	if len(metrics) == 0 || pattern == "" {
		return metrics, nil
	}

	if !strings.HasPrefix(pattern, prefix) {
		pattern = prefix + "." + pattern
	}

	matcher, err := glob.Compile(pattern)
	if err != nil {
		return metrics, errors.Wrapf(err, "filter: compile pattern %q", pattern)
	}

	filtered := make(provider.Metrics, 0, len(metrics))
	for _, metric := range metrics {
		if matcher.Match(string(metric)) {
			filtered = append(filtered, metric)
		}
	}
	return filtered, nil
}
