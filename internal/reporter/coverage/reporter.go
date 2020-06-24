package coverage

import (
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

func New() *reporter {
	return &reporter{}
}

type reporter struct{}

type Report struct {
	Metrics []Metric
	Total   float64
}

type Metric struct {
	Name string `json:"name"`
	Hits int    `json:"hits"`
}

func (reporter *reporter) Report(metrics entity.Metrics, queries entity.Queries) (*Report, error) {
	report := Report{Metrics: make([]Metric, 0, len(metrics))}
	coverage := make(map[entity.Metric]int, len(metrics))
	for _, query := range queries {
		matcher, err := glob.Compile(string(query))
		if err != nil {
			return nil, errors.Wrapf(err, "coverage: compile pattern %q", query)
		}
		for _, metric := range metrics {
			if matcher.Match(string(metric)) {
				coverage[metric]++
			}
		}
	}
	for _, metric := range metrics {
		report.Metrics = append(report.Metrics, Metric{Name: string(metric), Hits: coverage[metric]})
	}
	if len(metrics) > 0 {
		report.Total = 100 * float64(len(coverage)) / float64(len(metrics))
	}
	return &report, nil
}
