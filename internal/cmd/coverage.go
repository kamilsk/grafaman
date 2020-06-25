package cmd

import (
	"sort"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"
	xtime "go.octolab.org/time"
	"golang.org/x/sync/errgroup"

	"github.com/kamilsk/grafaman/internal/cache"
	"github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/filter"
	"github.com/kamilsk/grafaman/internal/model"
	entity "github.com/kamilsk/grafaman/internal/provider"
	"github.com/kamilsk/grafaman/internal/provider/grafana"
	"github.com/kamilsk/grafaman/internal/provider/graphite"
	"github.com/kamilsk/grafaman/internal/repl"
	"github.com/kamilsk/grafaman/internal/reporter/coverage"
	"github.com/kamilsk/grafaman/internal/validator"
)

// NewCoverageCommand returns command to calculate metrics coverage by queries.
func NewCoverageCommand(
	config *cnf.Config,
	logger *logrus.Logger,
	printer interface {
		SetPrefix(string)
		PrintCoverage(model.Report) error
	},
) *cobra.Command {
	var (
		exclude  []string
		trim     []string
		last     time.Duration
		noCache  bool
		replMode bool
	)

	command := cobra.Command{
		Use:   "coverage",
		Short: "calculates metrics coverage by queries",
		Long:  "Calculates metrics coverage by queries.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			fn.Must(
				func() error { return viper.BindPFlag("grafana_url", flags.Lookup("grafana")) },
				func() error { return viper.BindPFlag("grafana_dashboard", flags.Lookup("dashboard")) },
				func() error { return viper.BindPFlag("graphite_url", flags.Lookup("graphite")) },
				func() error { return viper.BindPFlag("graphite_metrics", flags.Lookup("metrics")) },
				func() error { return viper.BindPFlag("filter", flags.Lookup("filter")) },
				func() error { return viper.Unmarshal(config) },
			)

			if config.Grafana.URL == "" {
				return errors.New("please provide Grafana API endpoint")
			}
			if config.Grafana.Dashboard == "" {
				return errors.New("please provide a dashboard unique identifier")
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
			var (
				metrics   entity.Metrics
				dashboard *entity.Dashboard
			)

			g, ctx := errgroup.WithContext(cmd.Context())
			g.Go(func() error {
				var provider entity.Graphite
				provider, err := graphite.New(config.Graphite.URL, logger)
				if err != nil {
					return err
				}
				if !noCache {
					provider = cache.WrapGraphiteProvider(provider, afero.NewOsFs(), logger)
				}

				metrics, err = provider.Fetch(ctx, config.Graphite.Prefix, last)
				if err != nil {
					return err
				}
				for _, pattern := range exclude {
					metrics, err = filter.Exclude(metrics, pattern)
					if err != nil {
						break
					}
				}
				return err
			})
			g.Go(func() error {
				provider, err := grafana.New(config.Grafana.URL, logger)
				if err != nil {
					return err
				}
				dashboard, err = provider.Fetch(ctx, config.Grafana.Dashboard)
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

			reporter, err := coverage.New(queries)
			if err != nil {
				return err
			}

			printer.SetPrefix(config.Graphite.Prefix)
			if !replMode {
				metrics, err := filter.Filter(metrics, config.Graphite.Filter, config.Graphite.Prefix)
				if err != nil {
					return err
				}
				sort.Sort(metrics)

				return printer.PrintCoverage(reporter.Report(metrics))
			}
			prompt.New(
				repl.NewCoverageExecutor(config.Graphite.Prefix, metrics, queries, printer, logger),
				repl.NewMetricsCompleter(metrics),
			).Run()
			return nil
		},
	}

	flags := command.Flags()
	{
		flags.String("grafana", "", "Grafana API endpoint")
		flags.StringP("dashboard", "d", "", "a dashboard unique identifier")
		flags.String("graphite", "", "Graphite API endpoint")
		flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")
		flags.String("filter", "", "exclude metrics by pattern, e.g. some.*.metric")
	}
	flags.StringArrayVar(&exclude, "exclude", nil, "patterns to exclude metrics from coverage, e.g. *.median")
	flags.StringArrayVar(&trim, "trim", nil, "trim prefixes from queries")
	flags.DurationVar(&last, "last", xtime.Week, "the last interval to fetch")
	flags.BoolVar(&noCache, "no-cache", false, "disable caching")
	flags.BoolVar(&replMode, "repl", false, "enable repl mode")

	return &command
}
