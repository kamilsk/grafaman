package cmd

import (
	"sort"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"
	xtime "go.octolab.org/time"

	"github.com/kamilsk/grafaman/internal/cache"
	"github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/filter"
	entity "github.com/kamilsk/grafaman/internal/provider"
	"github.com/kamilsk/grafaman/internal/provider/graphite"
	"github.com/kamilsk/grafaman/internal/validator"
)

// NewMetricsCommand returns command to fetch metrics from Graphite.
func NewMetricsCommand(
	config *cnf.Config,
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
			fn.Must(
				func() error { return viper.BindPFlag("graphite_url", flags.Lookup("graphite")) },
				func() error { return viper.BindPFlag("graphite_metrics", flags.Lookup("metrics")) },
				func() error { return viper.BindPFlag("filter", flags.Lookup("filter")) },
				func() error { return viper.Unmarshal(config) },
			)

			if config.Graphite.URL == "" {
				return errors.New("please provide Graphite API endpoint")
			}
			if config.Graphite.Prefix == "" {
				return errors.New("please provide metric prefix")
			}
			checker := validator.Metric()
			if !checker(config.Graphite.Prefix) {
				return errors.Errorf(
					"invalid metric prefix: %s; it must be simple, e.g. apps.services.name",
					config.Graphite.Prefix,
				)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var provider entity.Graphite
			provider, err := graphite.New(config.Graphite.URL, logger)
			if err != nil {
				return err
			}
			provider = cache.Wrap(provider, afero.NewOsFs(), logger)

			metrics, err := provider.Fetch(cmd.Context(), config.Graphite.Prefix, last)
			if err != nil {
				return err
			}
			metrics, err = filter.Filter(metrics, config.Graphite.Filter, config.Graphite.Prefix)
			if err != nil {
				return err
			}
			sort.Sort(metrics)

			printer.SetPrefix(config.Graphite.Prefix)
			return printer.PrintMetrics(metrics)
		},
	}

	flags := command.Flags()
	{
		flags.StringP("graphite", "e", "", "Graphite API endpoint")
		flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")
		flags.String("filter", "", "exclude metrics by pattern, e.g. some.*.metric")
	}
	flags.IntVarP(&collapse, "collapse", "c", 0, "how many levels from the right to collapse by wildcard")
	flags.DurationVar(&last, "last", xtime.Week, "the last interval to fetch")

	return &command
}
