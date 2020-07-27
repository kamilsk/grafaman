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

type metricHit struct {
	Metric string `json:"name"`
	Hits   int    `json:"hits"`
}
