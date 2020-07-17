package reporter

import (
	"github.com/gobwas/glob"

	"github.com/kamilsk/grafaman/internal/model"
)

func MustNew(queries model.Queries) *reporter {
	return &reporter{queries.MustMatchers()}
}

type reporter struct {
	matchers []glob.Glob
}

func (reporter *reporter) CoverageReport(metrics model.Metrics) model.CoverageReport {
	report := new(model.CoverageReport)

	coverage := make(map[model.Metric]int, len(metrics))
	for _, matcher := range reporter.matchers {
		for _, metric := range metrics {
			if matcher.Match(string(metric)) {
				coverage[metric]++
			}
		}
	}

	for _, metric := range metrics {
		report.Add(metric, coverage[metric])
	}
	if len(metrics) > 0 {
		report.Total = 100 * float64(len(coverage)) / float64(len(metrics))
	}
	return *report
}
