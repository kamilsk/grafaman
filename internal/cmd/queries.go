package cmd

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/unsafe"

	entity "github.com/kamilsk/grafaman/internal/provider"
	"github.com/kamilsk/grafaman/internal/provider/grafana"
)

// TODO:debt
//  - validate metrics by regexp
//  - implement auth, if needed

// NewQueriesCommand returns command to fetch queries from a Grafana dashboard.
func NewQueriesCommand(style *simpletable.Style) *cobra.Command {
	var (
		trim       []string
		duplicates bool
		raw        bool
		sort       bool
	)
	command := cobra.Command{
		Use:   "queries",
		Short: "fetch queries from a Grafana dashboard",
		Long:  "Fetch queries from a Grafana dashboard.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			if err := viper.BindPFlag("grafana", flags.Lookup("grafana")); err != nil {
				return err
			}
			if err := viper.BindPFlag("dashboard", flags.Lookup("dashboard")); err != nil {
				return err
			}
			if err := viper.BindPFlag("metrics", flags.Lookup("metrics")); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			provider, err := grafana.New(viper.GetString("grafana"))
			if err != nil {
				return err
			}
			dashboard, err := provider.Fetch(cmd.Context(), viper.GetString("dashboard"))
			if err != nil {
				return err
			}
			dashboard.Prefix = viper.GetString("metrics")

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
	flags.StringP("grafana", "e", "", "Grafana API endpoint.")
	flags.StringP("dashboard", "d", "", "A dashboard unique identifier.")
	flags.StringP("metrics", "m", "", "The required subset of metrics. Must be a simple prefix.")
	{
		flags.StringArrayVar(&trim, "trim", nil, "Trim prefixes from queries.")
		flags.BoolVar(&duplicates, "allow-duplicates", false, "Allow duplicates of queries.")
		flags.BoolVar(&raw, "raw", false, "Leave the original values of queries.")
		flags.BoolVar(&sort, "sort", false, "Need to sort queries.")
	}
	return &command
}
