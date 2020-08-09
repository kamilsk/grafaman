package repl

import (
	"github.com/sirupsen/logrus"

	"github.com/kamilsk/grafaman/internal/model"
)

func NewCoverageExecutor(
	metrics model.Metrics,
	reporter CoverageReporter,
	printer CoverageReportPrinter,
	logger *logrus.Logger,
) func(string) {
	return func(q string) {
		metrics := metrics.Filter(model.Query(q).MustCompile()).Sort()
		if err := printer.PrintCoverageReport(reporter.CoverageReport(metrics)); err != nil {
			logger.WithError(err).Error("repl: print coverage report")
			return
		}
	}
}

func NewMetricsExecutor(
	metrics model.Metrics,
	printer MetricPrinter,
	logger *logrus.Logger,
) func(string) {
	return func(q string) {
		metrics := metrics.Filter(model.Query(q).MustCompile()).Sort()
		if err := printer.PrintMetrics(metrics); err != nil {
			logger.WithError(err).Error("repl: print metrics")
			return
		}
	}
}
