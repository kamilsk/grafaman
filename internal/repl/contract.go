package repl

import "github.com/kamilsk/grafaman/internal/model"

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test

// A CoverageReporter defines behavior of a coverage reporter.
type CoverageReporter interface {
	CoverageReport(model.Metrics) model.CoverageReport
}

// A CoverageReportPrinter defines behavior of a coverage report printer.
type CoverageReportPrinter interface {
	PrintCoverageReport(model.CoverageReport) error
}

// A MetricPrinter defines behavior of metrics printer.
type MetricPrinter interface {
	PrintMetrics(model.Metrics) error
}
