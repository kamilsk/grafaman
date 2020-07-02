package model

import (
	"reflect"
	"unsafe"
)

type Metric struct {
	Name string `json:"name"`
	Hits int    `json:"hits"`
}

type MetricName string

type MetricNames []MetricName

func (metrics *MetricNames) Convert(src []string) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	*metrics = *(*[]MetricName)(unsafe.Pointer(header))
}

func (metrics MetricNames) Len() int           { return len(metrics) }
func (metrics MetricNames) Less(i, j int) bool { return metrics[i] < metrics[j] }
func (metrics MetricNames) Swap(i, j int)      { metrics[i], metrics[j] = metrics[j], metrics[i] }
