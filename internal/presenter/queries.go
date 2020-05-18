package presenter

import (
	"fmt"
	"io"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

func PrintQueries(output io.Writer, queries entity.Queries, style *simpletable.Style) error {
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
			{Text: fmt.Sprintf("Total: %d", queries.Len())},
		},
	}
	table.SetStyle(style)

	_, err := fmt.Fprintln(output, table.String())
	return errors.Wrap(err, "presenter: output result")
}
