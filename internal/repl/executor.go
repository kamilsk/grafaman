package repl

import (
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/kamilsk/grafaman/internal/filter"
	"github.com/kamilsk/grafaman/internal/model"
	"github.com/kamilsk/grafaman/internal/provider"
)

func NewCoverageExecutor(
	metrics provider.Metrics,
	reporter interface {
		Report(provider.Metrics) model.Report
	},
	printer interface{ PrintCoverage(model.Report) error },
	logger *logrus.Logger,
) func(string) {
	return func(pattern string) {
		metrics, err := filter.Filter(metrics, pattern)
		if err != nil {
			logger.WithError(err).WithField("pattern", pattern).Error("repl: filter metrics")
			return
		}
		sort.Sort(metrics)

		if err := printer.PrintCoverage(reporter.Report(metrics)); err != nil {
			logger.WithError(err).Error("repl: print coverage report")
			return
		}
	}
}

func NewMetricsExecutor(
	metrics provider.Metrics,
	printer interface{ PrintMetrics(provider.Metrics) error },
	logger *logrus.Logger,
) func(string) {
	return func(pattern string) {
		metrics, err := filter.Filter(metrics, pattern)
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
