package coverage

import (
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

func New(raw entity.Queries) (*reporter, error) {
	matchers := make([]glob.Glob, 0, len(raw))
	for _, query := range raw {
		matcher, err := glob.Compile(string(query))
		if err != nil {
			return nil, errors.Wrapf(err, "coverage: compile pattern %q", query)
		}
		matchers = append(matchers, matcher)
	}
	return &reporter{matchers}, nil
}

type reporter struct {
	matchers []glob.Glob
}

type Report struct {
	Metrics []Metric
	Total   float64
}

type Metric struct {
	Name string `json:"name"`
	Hits int    `json:"hits"`
}

func (reporter *reporter) Report(metrics entity.Metrics) *Report {
	report := Report{Metrics: make([]Metric, 0, len(metrics))}
	coverage := make(map[entity.Metric]int, len(metrics))
	for _, matcher := range reporter.matchers {
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
	return &report
}
