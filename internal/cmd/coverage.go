package cmd

import (
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	xtime "go.octolab.org/time"
	"golang.org/x/sync/errgroup"

	entity "github.com/kamilsk/grafaman/internal/provider"
	"github.com/kamilsk/grafaman/internal/provider/grafana"
	"github.com/kamilsk/grafaman/internal/provider/graphite"
	"github.com/kamilsk/grafaman/internal/reporter/coverage"
	"github.com/kamilsk/grafaman/internal/validator"
)

// NewCoverageCommand returns command to calculate metrics coverage by queries.
func NewCoverageCommand(
	logger *logrus.Logger,
	printer interface{ PrintCoverage(*coverage.Report) error },
) *cobra.Command {
	var (
		exclude []string
		trim    []string
		last    time.Duration
	)
	command := cobra.Command{
		Use:   "coverage",
		Short: "calculates metrics coverage",
		Long:  "Calculates metrics coverage by queries.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			if err := viper.BindPFlag("grafana_url", flags.Lookup("grafana")); err != nil {
				return err
			}
			if err := viper.BindPFlag("grafana_dashboard", flags.Lookup("dashboard")); err != nil {
				return err
			}
			if err := viper.BindPFlag("graphite_url", flags.Lookup("graphite")); err != nil {
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
			metricsProvider, err := graphite.New(viper.GetString("graphite"), logger)
			if err != nil {
				return err
			}
			dashboardProvider, err := grafana.New(viper.GetString("grafana"), logger)
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
				metrics, err = metricsProvider.Fetch(ctx, viper.GetString("metrics"), last)
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

			return printer.PrintCoverage(report)
		},
	}
	flags := command.Flags()
	flags.String("grafana", "", "Grafana API endpoint")
	flags.StringP("dashboard", "d", "", "a dashboard unique identifier")
	flags.String("graphite", "", "Graphite API endpoint")
	flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")
	{
		flags.StringArrayVar(&exclude, "exclude", nil, "patterns to exclude metrics from coverage, e.g. *.median")
		flags.StringArrayVar(&trim, "trim", nil, "trim prefixes from queries")
		flags.DurationVar(&last, "last", xtime.Week, "the last interval to fetch")
	}
	return &command
}
