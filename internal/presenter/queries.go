package presenter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"

	"github.com/kamilsk/grafaman/internal/model"
)

func (printer *Printer) PrintQueries(queries model.Queries) error {
	switch printer.format {
	case formatJSON:
		return PrintQueriesAsJSON(printer.output, queries)
	case formatTSV:
		return PrintQueriesAsTSV(printer.output, queries)
	default:
		return PrintQueriesAsTable(printer.output, queries, styles[printer.format])
	}
}

func PrintQueriesAsJSON(output io.Writer, queries model.Queries) error {
	return errors.Wrap(json.NewEncoder(output).Encode(queries), "presenter: output result as json")
}

func PrintQueriesAsTable(output io.Writer, queries model.Queries, style *simpletable.Style) error {
	table := simpletable.New()
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Text: "Query"},
		},
	}
	for _, query := range queries {
		r := []*simpletable.Cell{
			{Text: string(query)},
		}
		table.Body.Cells = append(table.Body.Cells, r)
	}
	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Text: fmt.Sprintf("Total: %d", len(queries))},
		},
	}
	table.SetStyle(style)

	_, err := fmt.Fprintln(output, table.String())
	return errors.Wrap(err, "presenter: output result as table")
}

func PrintQueriesAsTSV(output io.Writer, queries model.Queries) error {
	for _, query := range queries {
		if _, err := fmt.Fprintln(output, query); err != nil {
			return errors.Wrap(err, "presenter: output result as TSV")
		}
	}
	return nil
}
