package presenter_test

import (
	"errors"
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/kamilsk/grafaman/internal/presenter"
)

var update = flag.Bool("update", false, "update golden files")

func TestPrinter_SetFormat(t *testing.T) {
	printer := new(Printer)
	assert.Error(t, printer.SetFormat("unknown"))
	assert.NoError(t, printer.SetFormat("json"))
}

// helpers

type unhealthy bool

func (fs unhealthy) Write([]byte) (int, error) {
	return 0, errors.New("unhealthy")
}
