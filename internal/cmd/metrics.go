package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	xtime "go.octolab.org/time"
	"go.octolab.org/unsafe"

	"github.com/kamilsk/grafaman/internal/provider/graphite"
	"github.com/kamilsk/grafaman/internal/validator"
)

// TODO:debt
//  - validate metrics by regexp
//  - support collapse option
//  - replace recursion by worker pool
//  - implement auth, if needed

// NewMetricsCommand returns command to fetch metrics from Graphite.
func NewMetricsCommand(style *simpletable.Style) *cobra.Command {
	var (
		collapse int
		last     time.Duration
		fast     bool
	)
	command := cobra.Command{
		Use:   "metrics",
		Short: "fetch metrics from Graphite",
		Long:  "Fetch metrics from Graphite.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			if err := viper.BindPFlag("graphite_url", flags.Lookup("graphite")); err != nil {
				return err
			}
			if err := viper.BindPFlag("graphite_metrics", flags.Lookup("metrics")); err != nil {
				return err
			}
			if viper.GetString("graphite") == "" {
				return errors.New("please provide Graphite API endpoint")
			}
			metrics, checker := viper.GetString("metrics"), validator.Metric()
			if metrics == "" {
				return errors.New("please provide metric prefix")
			}
			if !checker(metrics) {
				return errors.Errorf("invalid metric prefix: %s; it must be simple, e.g. apps.services.name", metrics)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			provider, err := graphite.New(viper.GetString("graphite"))
			if err != nil {
				return err
			}
			metrics, err := provider.Fetch(cmd.Context(), viper.GetString("metrics"), last, fast)
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
	flags.StringP("graphite", "e", "", "Graphite API endpoint")
	flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")
	{
		flags.IntVarP(&collapse, "collapse", "c", 0, "how many levels from the right to collapse by wildcard")
		flags.DurationVar(&last, "last", xtime.Week, "the last interval to fetch")
		flags.BoolVar(&fast, "fast", false, "use tilde `~` to fetch all metrics by one query if supported")
	}
	return &command
}
