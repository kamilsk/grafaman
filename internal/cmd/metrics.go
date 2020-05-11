package cmd

import (
	"fmt"
	"sort"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"go.octolab.org/fn"
	"go.octolab.org/unsafe"

	"github.com/kamilsk/grafaman/internal/provider/graphite"
)

// NewMetricsCommand returns command to fetch metrics from Graphite.
// TODO
//  - validate subset by regexp
//  - support collapse option
//  - replace recursion by worker pool
//  - implement auth, if needed
func NewMetricsCommand(style *simpletable.Style) *cobra.Command {
	var (
		endpoint string
		subset   string
		collapse int
	)
	command := cobra.Command{
		Use:   "metrics",
		Short: "fetch metrics from Graphite",
		Long:  "Fetch metrics from Graphite.",
		RunE: func(cmd *cobra.Command, args []string) error {
			provider, err := graphite.New(endpoint)
			if err != nil {
				return err
			}
			nodes, err := provider.Fetch(cmd.Context(), subset)
			if err != nil {
				return err
			}
			sort.Sort(nodes)

			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Text: "Metric name"},
				},
			}
			for _, metric := range nodes {
				r := []*simpletable.Cell{
					{Text: metric.ID},
				}
				table.Body.Cells = append(table.Body.Cells, r)
			}
			table.Footer = &simpletable.Footer{
				Cells: []*simpletable.Cell{
					{Text: fmt.Sprintf("Total: %d", nodes.Len())},
				},
			}
			table.SetStyle(style)

			unsafe.DoSilent(fmt.Fprintln(cmd.OutOrStdout(), table.String()))
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVarP(&endpoint, "endpoint", "e", "", "Graphite API endpoint.")
	flags.StringVarP(&subset, "subset", "s", "", "The required subset of metrics. Must be a simple prefix.")
	flags.IntVarP(&collapse, "collapse", "c", 0, "How many levels from the right to collapse by wildcard.")
	fn.Must(
		func() error { return command.MarkFlagRequired("endpoint") },
		func() error { return command.MarkFlagRequired("subset") },
	)
	return &command
}
