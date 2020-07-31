package cmd

import (
	"net/http"
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

	"github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/model"
	"github.com/kamilsk/grafaman/internal/provider/grafana"
	"github.com/kamilsk/grafaman/internal/provider/graphite"
	"github.com/kamilsk/grafaman/internal/provider/graphite/cache"
	"github.com/kamilsk/grafaman/internal/repl"
)

// NewCoverageCommand returns command to calculate metrics coverage by queries.
func NewCoverageCommand(
	config *cnf.Config,
	logger *logrus.Logger,
	printer interface {
		SetPrefix(string)
		PrintCoverage(model.CoverageReport) error
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
			if !model.Metric(config.Graphite.Prefix).Valid() {
				return errors.Errorf(
					"invalid metric prefix: %s; it must be simple, e.g. apps.services.name",
					config.Graphite.Prefix,
				)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				metrics   model.Metrics
				dashboard *model.Dashboard
			)

			g, ctx := errgroup.WithContext(cmd.Context())
			g.Go(func() error {
				var provider cache.Graphite
				provider, err := graphite.New(config.Graphite.URL, &http.Client{Timeout: time.Second}, logger)
				if err != nil {
					return err
				}
				if !noCache {
					provider = cache.Decorate(provider, afero.NewOsFs(), logger)
				}

				metrics, err = provider.Fetch(ctx, config.Graphite.Prefix, last)
				if err != nil {
					return err
				}

				metrics = metrics.Exclude(new(model.Queries).Convert(exclude).MustMatchers()...)
				return nil
			})
			g.Go(func() error {
				provider, err := grafana.New(config.Grafana.URL, &http.Client{Timeout: time.Second}, logger)
				if err != nil {
					return err
				}
				dashboard, err = provider.Fetch(ctx, config.Grafana.Dashboard)
				return err
			})
			if err := g.Wait(); err != nil {
				return err
			}

			queries, err := dashboard.Queries(model.Transform{
				SkipRaw:        false,
				SkipDuplicates: false,
				NeedSorting:    true,
				Unpack:         true,
				TrimPrefixes:   trim,
			})
			if err != nil {
				return err
			}

			coverage := model.NewCoverageReporter(queries)

			printer.SetPrefix(config.Graphite.Prefix)
			if !replMode {
				metrics := metrics.Filter(config.FilterQuery().MustCompile()).Sort()
				return printer.PrintCoverage(coverage.CoverageReport(metrics))
			}
			metrics.Sort()
			prompt.New(
				repl.Prefix(config.Graphite.Prefix, repl.NewCoverageExecutor(metrics, coverage, printer, logger)),
				repl.NewMetricsCompleter(config.Graphite.Prefix, metrics),
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
		flags.String("filter", "", "query to filter metrics, e.g. some.*.metric")
	}
	flags.StringArrayVar(&exclude, "exclude", nil, "queries to exclude metrics from coverage, e.g. *.median")
	flags.StringArrayVar(&trim, "trim", nil, "trim prefixes from queries")
	flags.DurationVar(&last, "last", xtime.Day, "the last interval to fetch")
	flags.BoolVar(&noCache, "no-cache", false, "disable caching")
	flags.BoolVar(&replMode, "repl", false, "enable repl mode")

	return &command
}
