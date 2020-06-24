package repl

import (
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/kamilsk/grafaman/internal/filter"
	"github.com/kamilsk/grafaman/internal/provider"
)

func NewExecutor(
	prefix string,
	metrics provider.Metrics,
	printer interface{ PrintMetrics(provider.Metrics) error },
	logger *logrus.Logger,
) func(string) {
	return func(pattern string) {
		metrics, err := filter.Filter(metrics, pattern, prefix)
		if err != nil {
			logger.WithError(err).WithField("pattern", pattern).Error("repl: filter metrics")
		}
		sort.Sort(metrics)

		if err := printer.PrintMetrics(metrics); err != nil {
			logger.WithError(err).Error("repl: print metrics")
		}
	}
}
