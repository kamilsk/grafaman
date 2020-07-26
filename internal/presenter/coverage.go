package presenter

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"

	"github.com/kamilsk/grafaman/internal/model"
)

func (printer *Printer) PrintCoverage(report model.CoverageReport) error {
	switch printer.format {
	case formatJSON:
		return PrintCoverageAsJSON(printer.output, report)
	case formatTSV:
		return PrintCoverageAsTSV(printer.output, report)
	default:
		return PrintCoverageAsTable(printer.output, report, styles[printer.format], printer.prefix)
	}
}

func PrintCoverageAsJSON(output io.Writer, report model.CoverageReport) error {
	return errors.Wrap(json.NewEncoder(output).Encode(report), "presenter: output result as json")
}

func PrintCoverageAsTable(output io.Writer, report model.CoverageReport, style *simpletable.Style, prefix string) error {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: fmt.Sprintf("Metric of %s", prefix)},
			{Text: "Hits"},
		},
	}
	for _, metric := range report.Metrics {
		r := []*simpletable.Cell{
			{Text: strings.TrimPrefix(strings.TrimPrefix(metric.Metric, prefix), ".")},
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
	return errors.Wrap(err, "presenter: output result as table")
}

func PrintCoverageAsTSV(output io.Writer, report model.CoverageReport) error {
	for _, metric := range report.Metrics {
		if _, err := fmt.Fprintln(output, metric.Metric, "\t", strconv.Itoa(metric.Hits)); err != nil {
			return errors.Wrap(err, "presenter: output result as TSV")
		}
	}
	return nil
}
