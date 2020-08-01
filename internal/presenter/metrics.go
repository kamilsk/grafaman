package presenter

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"

	"github.com/kamilsk/grafaman/internal/model"
)

// PrintMetrics prints metrics in a specific format.
func (printer *Printer) PrintMetrics(metrics model.Metrics) error {
	switch printer.format {
	case formatJSON:
		return printMetricsAsJSON(printer.output, metrics)
	case formatTSV:
		return printMetricsAsTSV(printer.output, metrics)
	default:
		return printMetricsAsTable(printer.output, metrics, styles[printer.format], printer.prefix)
	}
}

func printMetricsAsJSON(output io.Writer, metrics model.Metrics) error {
	return errors.Wrap(json.NewEncoder(output).Encode(metrics), "presenter: output result as json")
}

func printMetricsAsTable(output io.Writer, metrics model.Metrics, style *simpletable.Style, prefix string) error {
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
			{Text: fmt.Sprintf("Total: %d", len(metrics))},
		},
	}
	table.SetStyle(style)

	_, err := fmt.Fprintln(output, table.String())
	return errors.Wrap(err, "presenter: output result as table")
}

func printMetricsAsTSV(output io.Writer, metrics model.Metrics) error {
	for _, metric := range metrics {
		if _, err := fmt.Fprintln(output, metric); err != nil {
			return errors.Wrap(err, "presenter: output result as TSV")
		}
	}
	return nil
}
