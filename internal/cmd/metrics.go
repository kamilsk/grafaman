package cmd

import (
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	xtime "go.octolab.org/time"

	"github.com/kamilsk/grafaman/internal/cache"
	"github.com/kamilsk/grafaman/internal/filter"
	entity "github.com/kamilsk/grafaman/internal/provider"
	"github.com/kamilsk/grafaman/internal/provider/graphite"
	"github.com/kamilsk/grafaman/internal/validator"
)

// NewMetricsCommand returns command to fetch metrics from Graphite.
func NewMetricsCommand(
	logger *logrus.Logger,
	printer interface {
		SetPrefix(string)
		PrintMetrics(entity.Metrics) error
	},
) *cobra.Command {
	var (
		collapse int
		last     time.Duration
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
			if err := viper.BindPFlag("filter", flags.Lookup("filter")); err != nil {
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
			var provider cache.Graphite
			provider, err := graphite.New(viper.GetString("graphite"), logger)
			if err != nil {
				return err
			}
			provider = cache.Wrap(provider, afero.NewOsFs())
			metrics, err := provider.Fetch(cmd.Context(), viper.GetString("metrics"), last)
			if err != nil {
				return err
			}
			metrics, err = filter.Filter(metrics, viper.GetString("filter"), viper.GetString("metrics"))
			if err != nil {
				return err
			}
			sort.Sort(metrics)

			printer.SetPrefix(viper.GetString("metrics"))
			return printer.PrintMetrics(metrics)
		},
	}
	flags := command.Flags()
	flags.StringP("graphite", "e", "", "Graphite API endpoint")
	flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")
	flags.String("filter", "", "exclude metrics by pattern, e.g. some.*.metric")
	{
		flags.IntVarP(&collapse, "collapse", "c", 0, "how many levels from the right to collapse by wildcard")
		flags.DurationVar(&last, "last", xtime.Week, "the last interval to fetch")
	}
	return &command
}
