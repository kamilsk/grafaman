package filter

import (
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/kamilsk/grafaman/internal/provider"
)

func Filter(metrics provider.Metrics, pattern string) (provider.Metrics, error) {
	if pattern == "" {
		return metrics, nil
	}
	matcher, err := glob.Compile(pattern)
	if err != nil {
		return metrics, errors.Wrapf(err, "filter: compile pattern %q", pattern)
	}
	filtered := metrics[:0]
	for _, metric := range metrics {
		if matcher.Match(string(metric)) {
			filtered = append(filtered, metric)
		}
	}
	return filtered, nil
}
