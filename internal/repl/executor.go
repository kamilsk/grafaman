package repl

import (
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/kamilsk/grafaman/internal/filter"
	"github.com/kamilsk/grafaman/internal/model"
)

func NewCoverageExecutor(
	metrics model.MetricNames,
	reporter interface {
		Report(model.MetricNames) model.Report
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
	metrics model.MetricNames,
	printer interface{ PrintMetrics(model.MetricNames) error },
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
