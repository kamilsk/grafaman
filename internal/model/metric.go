package model

import (
	"reflect"
	"regexp"
	"sort"
	"unsafe"
)

// A Metric represents metric name.
type Metric string

// Valid checks that the metric is valid.
func (metric Metric) Valid() bool {
	return validator.MatchString(string(metric))
}

// Metrics represent a slice of metric names.
type Metrics []Metric

// Convert is an efficient converter from a slice of strings
// to a slice of metric names.
func (metrics *Metrics) Convert(src []string) Metrics {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	if metrics == nil {
		return *(*[]Metric)(unsafe.Pointer(header))
	}
	*metrics = *(*[]Metric)(unsafe.Pointer(header))
	return *metrics
}

// Exclude removes metrics matched by the matchers from the input list.
func (metrics Metrics) Exclude(matchers ...Matcher) Metrics {
	if len(metrics) == 0 || len(matchers) == 0 {
		return metrics
	}

	filtered := make(Metrics, 0, len(metrics))
	for _, metric := range metrics {
		var exclude bool
		for _, matcher := range matchers {
			if matcher.Match(string(metric)) {
				exclude = true
				break
			}
		}
		if !exclude {
			filtered = append(filtered, metric)
		}
	}
	return filtered
}

// Filter removes metrics that do not match by any matcher from the input list.
func (metrics Metrics) Filter(matchers ...Matcher) Metrics {
	if len(metrics) == 0 || len(matchers) == 0 {
		return metrics
	}

	filtered := make(Metrics, 0, len(metrics))
	for _, metric := range metrics {
		for _, matcher := range matchers {
			if matcher.Match(string(metric)) {
				filtered = append(filtered, metric)
				break
			}
		}
	}
	return filtered
}

// Sort orders the metrics in-place by ascending.
func (metrics Metrics) Sort() Metrics {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&metrics))
	sort.Strings(*(*[]string)(unsafe.Pointer(header)))
	return metrics
}

var validator = regexp.MustCompile(`^(?:[0-9a-z-_]+\.?)*$`)
