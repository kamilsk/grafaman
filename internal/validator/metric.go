package validator

import "regexp"

var metric = regexp.MustCompile(`^(?:[0-9a-z-_]+\.?)+$`)

func Metric() func(string) bool {
	return metric.MatchString
}
