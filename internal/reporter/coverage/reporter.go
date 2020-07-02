package coverage

import (
	"github.com/gobwas/glob"
	"github.com/pkg/errors"

	"github.com/kamilsk/grafaman/internal/model"
)

func New(raw model.Queries) (*reporter, error) {
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

func (reporter *reporter) Report(metrics model.Metrics) model.Report {
	report := new(model.Report)

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
