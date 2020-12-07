package cmd

import (
	"net/http"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	xtime "go.octolab.org/time"
	"golang.org/x/sync/errgroup"

	"github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/model"
	"github.com/kamilsk/grafaman/internal/presenter"
	"github.com/kamilsk/grafaman/internal/progress"
	"github.com/kamilsk/grafaman/internal/provider/grafana"
	"github.com/kamilsk/grafaman/internal/provider/graphite"
	"github.com/kamilsk/grafaman/internal/provider/graphite/cache"
	"github.com/kamilsk/grafaman/internal/repl"
)

// NewCoverageCommand returns command to calculate metrics coverage by queries.
func NewCoverageCommand(config *cnf.Config, logger *logrus.Logger) *cobra.Command {
	var (
		exclude  []string
		last     time.Duration
		noCache  bool
		replMode bool
	)

	command := cobra.Command{
		Use:   "coverage",
		Short: "calculates metrics coverage by queries",
		Long:  "Calculates metrics coverage by queries.",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if config.Grafana.URL == "" {
				return errors.New("please provide Grafana API endpoint")
			}
			if config.Grafana.Dashboard == "" {
				return errors.New("please provide a dashboard unique identifier")
			}
			if config.Graphite.URL == "" {
				return errors.New("please provide Graphite API endpoint")
			}
			if config.Graphite.Prefix == "" {
				return errors.New("please provide metric prefix")
			}
			if prefix := config.Graphite.Prefix; !model.Metric(prefix).Valid() {
				return errors.Errorf("invalid metric prefix: %s; it must be simple, e.g. apps.services.name", prefix)
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			printer := new(presenter.Printer)
			if err := printer.SetOutput(cmd.OutOrStdout()).SetFormat(config.Output.Format); err != nil {
				return err
			}
			printer.SetPrefix(config.Graphite.Prefix)

			indicator := progress.New()

			var (
				metrics   model.Metrics
				dashboard *model.Dashboard
			)

			g, ctx := errgroup.WithContext(cmd.Context())
			g.Go(func() error {
				var provider cache.Graphite
				provider, err := graphite.New(config.Graphite.URL, &http.Client{Timeout: config.Graphite.Timeout}, logger, indicator)
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
				provider, err := grafana.New(config.Grafana.URL, &http.Client{Timeout: config.Grafana.Timeout}, logger, indicator)
				if err != nil {
					return err
				}
				dashboard, err = provider.Fetch(ctx, config.Grafana.Dashboard)
				return err
			})
			if err := g.Wait(); err != nil {
				return err
			}

			queries, err := dashboard.Queries(model.Config{
				SkipRaw:        false,
				SkipDuplicates: false,
				NeedSorting:    true,
				Unpack:         true,
			})
			if err != nil {
				return err
			}

			reporter := model.NewCoverageReporter(queries)

			if !replMode {
				metrics := metrics.Filter(config.FilterQuery().MustCompile()).Sort()
				return printer.PrintCoverageReport(reporter.CoverageReport(metrics))
			}
			metrics.Sort()
			prompt.New(
				repl.Prefix(config.Graphite.Prefix, repl.NewCoverageReportExecutor(metrics, reporter, printer, logger)),
				repl.NewMetricsCompleter(config.Graphite.Prefix, metrics),
			).Run()
			return nil
		},
	}

	flags := command.Flags()
	flags.StringArrayVar(&exclude, "exclude", nil, "queries to exclude metrics from coverage, e.g. *.median")
	flags.DurationVar(&last, "last", xtime.Day, "the last interval to fetch")
	flags.BoolVar(&noCache, "no-cache", false, "disable caching")
	flags.BoolVar(&replMode, "repl", false, "enable repl mode")

	return &command
}
