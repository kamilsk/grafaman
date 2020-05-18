package presenter

import (
	"fmt"
	"io"
	"strconv"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"

	"github.com/kamilsk/grafaman/internal/reporter/coverage"
)

func PrintCoverage(output io.Writer, report *coverage.Report, style *simpletable.Style) error {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "Metric"},
			{Text: "Hits"},
		},
	}
	for _, metric := range report.Metrics {
		r := []*simpletable.Cell{
			{Text: metric.Name},
			{Align: simpletable.AlignRight, Text: strconv.Itoa(metric.Hits)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: "Total"},
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%.2f%%", report.Total)},
		},
	}
	table.SetStyle(style)

	_, err := fmt.Fprintln(output, table.String())
	return errors.Wrap(err, "presenter: output result")
}
