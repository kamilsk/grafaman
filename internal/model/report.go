package model

import "encoding/json"

// A CoverageReport contains information about which metrics
// are covered, and which not.
type CoverageReport struct {
	Metrics []metricHit
}

// Add registers the metric and its hit count.
func (report *CoverageReport) Add(name Metric, hits int) {
	report.Metrics = append(report.Metrics, metricHit{string(name), hits})
}

// MarshalJSON implements the Marshaler interface of the json package.
func (report *CoverageReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(report.Metrics)
}

// Total returns coverage value of the report.
func (report *CoverageReport) Total() float64 {
	if len(report.Metrics) == 0 {
		return 0.0
	}
	var hits int
	for _, hit := range report.Metrics {
		if hit.Hits > 0 {
			hits++
		}
	}
	return 100 * float64(hits) / float64(len(report.Metrics))
}

// NewCoverageReporter returns new metric coverage reporter.
func NewCoverageReporter(queries Queries) *reporter {
	return &reporter{queries.MustMatchers()}
}

type reporter struct {
	matchers []Matcher
}

// CoverageReport builds metric coverage report.
func (reporter *reporter) CoverageReport(metrics Metrics) CoverageReport {
	var report CoverageReport

	coverage := make(map[Metric]int, len(metrics))
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
	return report
}

type metricHit struct {
	Metric string `json:"name"`
	Hits   int    `json:"hits"`
}
