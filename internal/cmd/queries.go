package cmd

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"go.octolab.org/fn"
	"go.octolab.org/unsafe"

	entity "github.com/kamilsk/grafaman/internal/provider"
	"github.com/kamilsk/grafaman/internal/provider/grafana"
)

// NewQueriesCommand returns command to fetch queries from a Grafana dashboard.
// TODO
//  - validate subset by regexp
//  - implement auth, if needed
func NewQueriesCommand(style *simpletable.Style) *cobra.Command {
	var (
		endpoint   string
		uid        string
		subset     string
		trim       []string
		raw        bool
		duplicates bool
		sort       bool
	)
	command := cobra.Command{
		Use:   "queries",
		Short: "fetch queries from a Grafana dashboard",
		Long:  "Fetch queries from a Grafana dashboard.",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider, err := grafana.New(endpoint)
			if err != nil {
				return err
			}
			dashboard, err := provider.Fetch(cmd.Context(), uid)
			if err != nil {
				return err
			}
			dashboard.Subset = subset

			queries, err := dashboard.Queries(entity.Transform{
				SkipRaw:        raw,
				SkipDuplicates: duplicates,
				TrimPrefixes:   trim,
				NeedSorting:    sort,
			})
			if err != nil {
				return err
			}

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

			unsafe.DoSilent(fmt.Fprintln(cmd.OutOrStdout(), table.String()))
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVarP(&endpoint, "endpoint", "e", "", "Grafana API endpoint.")
	flags.StringVarP(&uid, "dashboard", "d", "", "A dashboard unique identifier.")
	flags.StringVarP(&subset, "subset", "s", "", "The required subset of metrics. Must be a simple prefix.")
	flags.StringArrayVarP(&trim, "trim", "t", nil, "Trim prefixes from queries.")
	flags.BoolVar(&raw, "raw", false, "Leave the original values of queries.")
	flags.BoolVar(&duplicates, "allow-duplicates", false, "Allow duplicates of queries.")
	flags.BoolVar(&sort, "sort", false, "Need to sort queries.")
	fn.Must(
		func() error { return command.MarkFlagRequired("endpoint") },
		func() error { return command.MarkFlagRequired("dashboard") },
	)
	return &command
}
