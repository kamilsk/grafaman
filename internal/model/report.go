package model

import "encoding/json"

// A CoverageReport contains information about which metrics
// are covered, and which not.
type CoverageReport struct {
	Metrics []metricHit
	Total   float64
}

// Add registers the metric and its hit count.
func (report *CoverageReport) Add(name Metric, hits int) {
	report.Metrics = append(report.Metrics, metricHit{string(name), hits})
}

// MarshalJSON implements the Marshaler interface of the json package.
func (report *CoverageReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(report.Metrics)
}

type metricHit struct {
	Name string `json:"name"`
	Hits int    `json:"hits"`
}
