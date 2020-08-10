package repl_test

import (
	"errors"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/kamilsk/grafaman/internal/model"
	. "github.com/kamilsk/grafaman/internal/repl"
)

func TestCoverageReportExecutor(t *testing.T) {
	metrics := model.Metrics{
		"metric.a.ok",
		"metric.b.ok",
		"metric.c.ok",
	}
	report := model.CoverageReport{}
	report.Add(metrics[0], 1)
	report.Add(metrics[1], 0)
	report.Add(metrics[2], 2)

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reporter := NewMockCoverageReporter(ctrl)
		reporter.EXPECT().
			CoverageReport(metrics).
			Return(report)

		printer := NewMockCoverageReportPrinter(ctrl)
		printer.EXPECT().
			PrintCoverageReport(report).
			Return(nil)

		executor := NewCoverageReportExecutor(metrics, reporter, printer, logger)
		assert.NotPanics(t, func() { executor("metric.*") })
	})

	t.Run("failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reporter := NewMockCoverageReporter(ctrl)
		reporter.EXPECT().
			CoverageReport(metrics).
			Return(report)

		printer := NewMockCoverageReportPrinter(ctrl)
		printer.EXPECT().
			PrintCoverageReport(report).
			Return(errors.New("unhealthy"))

		executor := NewCoverageReportExecutor(metrics, reporter, printer, logger)
		assert.NotPanics(t, func() { executor("metric.*") })
	})
}

func TestMetricExecutor(t *testing.T) {
	metrics := model.Metrics{
		"metric.a.ok",
		"metric.b.ok",
		"metric.c.ok",
	}

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		printer := NewMockMetricPrinter(ctrl)
		printer.EXPECT().
			PrintMetrics(metrics).
			Return(nil)

		executor := NewMetricExecutor(metrics, printer, logger)
		assert.NotPanics(t, func() { executor("metric.*") })
	})

	t.Run("failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		printer := NewMockMetricPrinter(ctrl)
		printer.EXPECT().
			PrintMetrics(metrics).
			Return(errors.New("unhealthy"))

		executor := NewMetricExecutor(metrics, printer, logger)
		assert.NotPanics(t, func() { executor("metric.*") })
	})
}
