package cmd

import "github.com/kamilsk/grafaman/internal/model"

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test

// A CoverageReportPrinter defines behavior of a coverage report printer.
type CoverageReportPrinter interface {
	SetPrefix(string)
	PrintCoverageReport(model.CoverageReport) error
}

// A MetricPrinter defines behavior of metrics printer.
type MetricPrinter interface {
	SetPrefix(string)
	PrintMetrics(model.Metrics) error
}

// A QueryPrinter defines behavior of queries printer.
type QueryPrinter interface {
	PrintQueries(model.Queries) error
}
