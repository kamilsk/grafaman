package model

import (
	"reflect"
	"unsafe"
)

type Metric string

type Metrics []Metric

func (metrics *Metrics) Convert(src []string) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	*metrics = *(*[]Metric)(unsafe.Pointer(header))
}

func (metrics Metrics) Len() int           { return len(metrics) }
func (metrics Metrics) Less(i, j int) bool { return metrics[i] < metrics[j] }
func (metrics Metrics) Swap(i, j int)      { metrics[i], metrics[j] = metrics[j], metrics[i] }
