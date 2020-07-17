package repl

import (
	"github.com/sirupsen/logrus"

	"github.com/kamilsk/grafaman/internal/model"
)

func NewCoverageExecutor(
	metrics model.Metrics,
	reporter interface {
		CoverageReport(model.Metrics) model.CoverageReport
	},
	printer interface {
		PrintCoverage(model.CoverageReport) error
	},
	logger *logrus.Logger,
) func(string) {
	return func(pattern string) {
		metrics := metrics.Filter(model.Query(pattern).MustCompile()).Sort()
		if err := printer.PrintCoverage(reporter.CoverageReport(metrics)); err != nil {
			logger.WithError(err).Error("repl: print coverage report")
			return
		}
	}
}

func NewMetricsExecutor(
	metrics model.Metrics,
	printer interface{ PrintMetrics(model.Metrics) error },
	logger *logrus.Logger,
) func(string) {
	return func(pattern string) {
		metrics := metrics.Filter(model.Query(pattern).MustCompile()).Sort()
		if err := printer.PrintMetrics(metrics); err != nil {
			logger.WithError(err).Error("repl: print metrics")
			return
		}
	}
}
