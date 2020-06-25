package model

type Report struct {
	Metrics []Metric
	Total   float64
}

type Metric struct {
	Name string `json:"name"`
	Hits int    `json:"hits"`
}
