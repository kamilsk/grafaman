package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	xtime "go.octolab.org/time"
	"go.octolab.org/unsafe"
	"golang.org/x/sync/errgroup"

	entity "github.com/kamilsk/grafaman/internal/provider"
	"github.com/kamilsk/grafaman/internal/provider/grafana"
	"github.com/kamilsk/grafaman/internal/provider/graphite"
	"github.com/kamilsk/grafaman/internal/reporter/coverage"
)

// NewCoverageCommand returns command to calculate metrics coverage by queries.
// TODO:debt
//  - validate subset by regexp
//  - support last option
//  - support collapse option
//  - support graphite functions (e.g. sum, etc.)
//  - implement auth, if needed
func NewCoverageCommand(style *simpletable.Style) *cobra.Command {
	var (
		exclude []string
		trim    []string
		last    time.Duration
		fast    bool
	)
	command := cobra.Command{
		Use:   "coverage",
		Short: "calculates metrics coverage",
		Long:  "Calculates metrics coverage by queries.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			if err := viper.BindPFlag("grafana", flags.Lookup("grafana")); err != nil {
				return err
			}
			if err := viper.BindPFlag("dashboard", flags.Lookup("dashboard")); err != nil {
				return err
			}
			if err := viper.BindPFlag("graphite", flags.Lookup("graphite")); err != nil {
				return err
			}
			if err := viper.BindPFlag("metrics", flags.Lookup("metrics")); err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			metricsProvider, err := graphite.New(viper.GetString("graphite"))
			if err != nil {
				return err
			}
			dashboardProvider, err := grafana.New(viper.GetString("grafana"))
			if err != nil {
				return err
			}

			var (
				metrics   entity.Metrics
				dashboard *entity.Dashboard
			)
			g, ctx := errgroup.WithContext(cmd.Context())
			g.Go(func() error {
				var err error
				metrics, err = metricsProvider.Fetch(ctx, viper.GetString("metrics"), fast)
				return err
			})
			g.Go(func() error {
				var err error
				dashboard, err = dashboardProvider.Fetch(ctx, viper.GetString("dashboard"))
				return err
			})
			if err := g.Wait(); err != nil {
				return err
			}

			queries, err := dashboard.Queries(entity.Transform{
				SkipRaw:        false,
				SkipDuplicates: false,
				NeedSorting:    true,
				Unpack:         true,
				TrimPrefixes:   trim,
			})
			if err != nil {
				return err
			}

			reporter := coverage.New(exclude)
			report, err := reporter.Report(metrics, queries)
			if err != nil {
				return err
			}

			table := simpletable.New()
			table.Header = &simpletable.Header{
				Cells: []*simpletable.Cell{
					{Text: "Metric"},
					{Text: "Hits"},
				},
			}
			for _, metric := range report.Metrics {
				r := []*simpletable.Cell{
					{Text: metric.Name},
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

			unsafe.DoSilent(fmt.Fprintln(cmd.OutOrStdout(), table.String()))
			return nil
		},
	}
	flags := command.Flags()
	flags.String("grafana", "", "Grafana API endpoint.")
	flags.StringP("dashboard", "d", "", "A dashboard unique identifier.")
	flags.String("graphite", "", "Graphite API endpoint.")
	flags.StringP("metrics", "m", "", "The required subset of metrics. Must be a simple prefix.")
	{
		flags.StringArrayVar(&exclude, "exclude", nil, "Patterns to exclude metrics from coverage, e.g. *.median")
		flags.StringArrayVar(&trim, "trim", nil, "Trim prefixes from queries.")
		flags.DurationVar(&last, "last", xtime.Week, "The last interval to fetch.")
		flags.BoolVar(&fast, "fast", false, "Use tilde `~` to fetch all metrics by one query if supported.")
	}
	return &command
}
