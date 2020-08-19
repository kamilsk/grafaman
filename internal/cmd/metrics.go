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

	"github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/model"
	"github.com/kamilsk/grafaman/internal/provider/graphite"
	"github.com/kamilsk/grafaman/internal/provider/graphite/cache"
	"github.com/kamilsk/grafaman/internal/repl"
)

// NewMetricsCommand returns command to fetch metrics from Graphite.
func NewMetricsCommand(
	config *cnf.Config,
	logger *logrus.Logger,
	printer MetricPrinter,
) *cobra.Command {
	var (
		last     time.Duration
		noCache  bool
		replMode bool
	)

	command := cobra.Command{
		Use:   "metrics",
		Short: "fetch metrics from Graphite",
		Long:  "Fetch metrics from Graphite.",

		PreRunE: func(cmd *cobra.Command, args []string) error {
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
			var provider cache.Graphite
			provider, err := graphite.New(config.Graphite.URL, &http.Client{Timeout: time.Second}, logger)
			if err != nil {
				return err
			}
			if !noCache {
				provider = cache.Decorate(provider, afero.NewOsFs(), logger)
			}

			metrics, err := provider.Fetch(cmd.Context(), config.Graphite.Prefix, last)
			if err != nil {
				return err
			}

			printer.SetPrefix(config.Graphite.Prefix)
			if !replMode {
				metrics = metrics.Filter(config.FilterQuery().MustCompile()).Sort()
				return printer.PrintMetrics(metrics)
			}
			metrics.Sort()
			prompt.New(
				repl.Prefix(config.Graphite.Prefix, repl.NewMetricExecutor(metrics, printer, logger)),
				repl.NewMetricsCompleter(config.Graphite.Prefix, metrics),
			).Run()
			return nil
		},
	}

	flags := command.Flags()
	flags.DurationVar(&last, "last", xtime.Day, "the last interval to fetch")
	flags.BoolVar(&noCache, "no-cache", false, "disable caching")
	flags.BoolVar(&replMode, "repl", false, "enable repl mode")

	return &command
}
