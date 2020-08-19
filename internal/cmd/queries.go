package cmd

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

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
	var cfg model.Config

	command := cobra.Command{
		Use:   "queries",
		Short: "fetch queries from a Grafana dashboard",
		Long:  "Fetch queries from a Grafana dashboard.",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if config.Grafana.URL == "" {
				return errors.New("please provide Grafana API endpoint")
			}
			if config.Grafana.Dashboard == "" {
				return errors.New("please provide a dashboard unique identifier")
			}
			if prefix := config.Graphite.Prefix; prefix != "" && !model.Metric(prefix).Valid() {
				return errors.Errorf("invalid metric prefix: %s; it must be simple, e.g. apps.services.name", prefix)
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
			queries, err := dashboard.Queries(cfg)
			if err != nil {
				return err
			}

			return printer.PrintQueries(queries)
		},
	}

	flags := command.Flags()
	flags.BoolVar(&cfg.SkipDuplicates, "allow-duplicates", false, "allow duplicates of queries")
	flags.BoolVar(&cfg.SkipRaw, "raw", false, "leave the original values of queries")
	flags.BoolVar(&cfg.NeedSorting, "sort", false, "need to sort queries")

	return &command
}
