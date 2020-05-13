package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"go.octolab.org/fn"
	xtime "go.octolab.org/time"
	"go.octolab.org/unsafe"

	"github.com/kamilsk/grafaman/internal/provider/graphite"
)

// NewMetricsCommand returns command to fetch metrics from Graphite.
// TODO:debt
//  - validate subset by regexp
//  - support collapse option
//  - support last option
//  - replace recursion by worker pool
//  - implement auth, if needed
func NewMetricsCommand(style *simpletable.Style) *cobra.Command {
	var (
		endpoint string
		subset   string
		collapse int
		last     time.Duration
		fast     bool
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
			metrics, err := provider.Fetch(cmd.Context(), subset, fast)
			if err != nil {
				return err
			}
			sort.Sort(metrics)

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

			unsafe.DoSilent(fmt.Fprintln(cmd.OutOrStdout(), table.String()))
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVarP(&endpoint, "endpoint", "e", "", "Graphite API endpoint.")
	flags.StringVarP(&subset, "subset", "s", "", "The required subset of metrics. Must be a simple prefix.")
	flags.IntVarP(&collapse, "collapse", "c", 0, "How many levels from the right to collapse by wildcard.")
	flags.DurationVar(&last, "last", xtime.Week, "The last interval to fetch.")
	flags.BoolVar(&fast, "fast", false, "Use tilde `~` to fetch all metrics by one query if supported.")
	fn.Must(
		func() error { return command.MarkFlagRequired("endpoint") },
		func() error { return command.MarkFlagRequired("subset") },
	)
	return &command
}
