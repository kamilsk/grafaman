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
	report.Add("a", 1)
	report.Add("b", 0)
	report.Add("c", 2)

	buf := bytes.NewBuffer(make([]byte, 0, 512))
	require.NoError(t, json.NewEncoder(buf).Encode(report))

	expected, err := ioutil.ReadFile("testdata/report.json")
	require.NoError(t, err)

	assert.Equal(t, expected, buf.Bytes())
}
