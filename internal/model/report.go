package model

type Report struct {
	Metrics []metric
	Total   float64
}

type metric struct {
	Name string `json:"name"`
	Hits int    `json:"hits"`
}

func (report *Report) Add(name Metric, hits int) {
	report.Metrics = append(report.Metrics, metric{string(name), hits})
}
