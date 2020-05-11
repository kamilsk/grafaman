package provider

type Metric string

type Metrics []Metric

func (metrics Metrics) Len() int           { return len(metrics) }
func (metrics Metrics) Less(i, j int) bool { return metrics[i] < metrics[j] }
func (metrics Metrics) Swap(i, j int)      { metrics[i], metrics[j] = metrics[j], metrics[i] }
