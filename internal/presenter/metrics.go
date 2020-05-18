package presenter

import (
	"fmt"
	"io"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

func PrintMetrics(output io.Writer, metrics entity.Metrics, style *simpletable.Style) error {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "Metric"},
		},
	}
	for _, metric := range metrics {
		r := []*simpletable.Cell{
			{Text: string(metric)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Text: fmt.Sprintf("Total: %d", metrics.Len())},
		},
	}
	table.SetStyle(style)

	_, err := fmt.Fprintln(output, table.String())
	return errors.Wrap(err, "presenter: output result")
}
