package model_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/kamilsk/grafaman/internal/model"
)

func TestReport(t *testing.T) {
	report := new(CoverageReport)
	assert.Equal(t, 0.0, report.Total())

	report.Add("a", 1)
	report.Add("b", 0)
	report.Add("c", 2)
	assert.Equal(t, 100*float64(2)/float64(3), report.Total())

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	require.NoError(t, json.NewEncoder(buf).Encode(report))

	expected, err := ioutil.ReadFile("testdata/report.json")
	require.NoError(t, err)

	assert.Equal(t, expected, buf.Bytes())
}

func TestCoverageReporter(t *testing.T) {
	t.Run("full covered", func(t *testing.T) {
		reporter := NewCoverageReporter(Queries{"metric.*"})
		report := reporter.CoverageReport(Metrics{
			"metric.a",
			"metric.b",
			"metric.c",
		})
		assert.Len(t, report.Metrics, 3)
		assert.Equal(t, 100.0, report.Total())
	})

	t.Run("partial covered", func(t *testing.T) {
		reporter := NewCoverageReporter(Queries{"metric.b"})
		report := reporter.CoverageReport(Metrics{
			"metric.a",
			"metric.b",
			"metric.c",
		})
		assert.Len(t, report.Metrics, 3)
		assert.Equal(t, 100*float64(1)/float64(3), report.Total())
	})

	t.Run("without matchers", func(t *testing.T) {
		reporter := NewCoverageReporter(nil)
		report := reporter.CoverageReport(Metrics{
			"metric.a",
			"metric.b",
			"metric.c",
		})
		assert.Len(t, report.Metrics, 3)
		assert.Equal(t, 0.0, report.Total())
	})
}
