package cmd

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"

	"github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/model"
	"github.com/kamilsk/grafaman/internal/provider/grafana"
)

// NewQueriesCommand returns command to fetch queries from a Grafana dashboard.
func NewQueriesCommand(
	config *cnf.Config,
	logger *logrus.Logger,
	printer QueryPrinter,
) *cobra.Command {
	var (
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
			fn.Must(
				func() error { return viper.BindPFlag("grafana_url", flags.Lookup("grafana")) },
				func() error { return viper.BindPFlag("grafana_dashboard", flags.Lookup("dashboard")) },
				func() error { return viper.BindPFlag("graphite_metrics", flags.Lookup("metrics")) },
				func() error { return viper.Unmarshal(config) },
			)

			if config.Grafana.URL == "" {
				return errors.New("please provide Grafana API endpoint")
			}
			if config.Grafana.Dashboard == "" {
				return errors.New("please provide a dashboard unique identifier")
			}
			if config.Graphite.Prefix != "" && !model.Metric(config.Graphite.Prefix).Valid() {
				return errors.Errorf(
					"invalid metric prefix: %s; it must be simple, e.g. apps.services.name",
					config.Graphite.Prefix,
				)
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			provider, err := grafana.New(config.Grafana.URL, &http.Client{Timeout: time.Second}, logger)
			if err != nil {
				return err
			}
			dashboard, err := provider.Fetch(cmd.Context(), config.Grafana.Dashboard)
			if err != nil {
				return err
			}

			dashboard.Prefix = config.Graphite.Prefix
			queries, err := dashboard.Queries(model.Config{
				SkipRaw:        raw,
				SkipDuplicates: duplicates,
				NeedSorting:    sort,
			})
			if err != nil {
				return err
			}

			return printer.PrintQueries(queries)
		},
	}

	flags := command.Flags()
	{
		flags.StringP("grafana", "e", "", "Grafana API endpoint")
		flags.StringP("dashboard", "d", "", "a dashboard unique identifier")
		flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")
	}
	flags.BoolVar(&duplicates, "allow-duplicates", false, "allow duplicates of queries")
	flags.BoolVar(&raw, "raw", false, "leave the original values of queries")
	flags.BoolVar(&sort, "sort", false, "need to sort queries")

	return &command
}
