package presenter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

func (printer *Printer) PrintMetrics(metrics entity.Metrics) error {
	switch printer.format {
	case formatJSON:
		return PrintMetricsAsJSON(printer.output, metrics)
	case formatTSV:
		return PrintMetricsAsTSV(printer.output, metrics)
	default:
		return PrintMetricsAsTable(printer.output, metrics, styles[printer.format], printer.prefix)
	}
}

func PrintMetricsAsJSON(output io.Writer, metrics entity.Metrics) error {
	return errors.Wrap(json.NewEncoder(output).Encode(metrics), "presenter: output result as json")
}

func PrintMetricsAsTable(output io.Writer, metrics entity.Metrics, style *simpletable.Style, prefix string) error {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: fmt.Sprintf("Metric of %s", prefix)},
		},
	}
	for _, metric := range metrics {
		r := []*simpletable.Cell{
			{Text: strings.TrimPrefix(strings.TrimPrefix(string(metric), prefix), ".")},
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
	return errors.Wrap(err, "presenter: output result as table")
}

func PrintMetricsAsTSV(output io.Writer, metrics entity.Metrics) error {
	for _, metric := range metrics {
		if _, err := fmt.Fprintln(output, metric); err != nil {
			return errors.Wrap(err, "presenter: output result as TSV")
		}
	}
	return nil
}
