package repl

import (
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/kamilsk/grafaman/internal/filter"
	"github.com/kamilsk/grafaman/internal/provider"
	"github.com/kamilsk/grafaman/internal/reporter/coverage"
)

func NewCoverageExecutor(
	prefix string,
	metrics provider.Metrics,
	queries provider.Queries,
	printer interface{ PrintCoverage(*coverage.Report) error },
	logger *logrus.Logger,
) func(string) {
	return func(pattern string) {
		metrics, err := filter.Filter(metrics, pattern, prefix)
		if err != nil {
			logger.WithError(err).WithField("pattern", pattern).Error("repl: filter metrics")
			return
		}
		sort.Sort(metrics)

		reporter := coverage.New(queries)
		report, err := reporter.Report(metrics)
		if err != nil {
			logger.WithError(err).Error("repl: make report")
			return
		}

		if err := printer.PrintCoverage(report); err != nil {
			logger.WithError(err).Error("repl: print coverage report")
			return
		}
	}
}

func NewMetricsExecutor(
	prefix string,
	metrics provider.Metrics,
	printer interface{ PrintMetrics(provider.Metrics) error },
	logger *logrus.Logger,
) func(string) {
	return func(pattern string) {
		metrics, err := filter.Filter(metrics, pattern, prefix)
		if err != nil {
			logger.WithError(err).WithField("pattern", pattern).Error("repl: filter metrics")
			return
		}
		sort.Sort(metrics)

		if err := printer.PrintMetrics(metrics); err != nil {
			logger.WithError(err).Error("repl: print metrics")
			return
		}
	}
}
