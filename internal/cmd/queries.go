package cmd

import (
	"fmt"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/unsafe"

	entity "github.com/kamilsk/grafaman/internal/provider"
	"github.com/kamilsk/grafaman/internal/provider/grafana"
	"github.com/kamilsk/grafaman/internal/validator"
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
			if err := viper.BindPFlag("grafana_url", flags.Lookup("grafana")); err != nil {
				return err
			}
			if err := viper.BindPFlag("grafana_dashboard", flags.Lookup("dashboard")); err != nil {
				return err
			}
			if err := viper.BindPFlag("graphite_metrics", flags.Lookup("metrics")); err != nil {
				return err
			}
			if viper.GetString("grafana") == "" {
				return errors.New("please provide Grafana API endpoint")
			}
			if viper.GetString("dashboard") == "" {
				return errors.New("please provide a dashboard unique identifier")
			}
			if metrics, checker := viper.GetString("metrics"), validator.Metric(); metrics != "" && !checker(metrics) {
				return errors.Errorf("invalid metric prefix: %s; it must be simple, see examples", metrics)
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
	flags.StringP("grafana", "e", "", "Grafana API endpoint")
	flags.StringP("dashboard", "d", "", "a dashboard unique identifier")
	flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")
	{
		flags.StringArrayVar(&trim, "trim", nil, "trim prefixes from queries")
		flags.BoolVar(&duplicates, "allow-duplicates", false, "allow duplicates of queries")
		flags.BoolVar(&raw, "raw", false, "leave the original values of queries")
		flags.BoolVar(&sort, "sort", false, "need to sort queries")
	}
	return &command
}
