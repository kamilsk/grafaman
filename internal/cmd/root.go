package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/kamilsk/grafaman/internal/cnf"
)

// New returns the new root command.
func New() *cobra.Command {
	var (
		config = new(cnf.Config)
		logger = logrus.New()
	)

	command := cobra.Command{
		Use:   "grafaman",
		Short: "metrics coverage reporter",
		Long:  "Metrics coverage reporter for Graphite and Grafana.",

		SilenceErrors: false,
		SilenceUsage:  true,
	}

	flags := command.PersistentFlags()
	flags.StringVar(&config.File, "env-file", ".env.paas", "file with environment variables; fallback to app.toml")

	command.AddCommand(
		cnf.Apply(
			NewCacheLookupCommand(config, logger), viper.New(),
			cnf.WithConfig(config),
			cnf.WithGraphiteMetrics(),
		),
		cnf.Apply(
			NewCoverageCommand(config, logger), viper.New(),
			cnf.WithConfig(config),
			cnf.WithDebug(config, logger),
			cnf.WithGrafana(),
			cnf.WithGraphite(),
			cnf.WithOutputFormat(),
		),
		cnf.Apply(
			NewMetricsCommand(config, logger), viper.New(),
			cnf.WithConfig(config),
			cnf.WithDebug(config, logger),
			cnf.WithGraphite(),
			cnf.WithOutputFormat(),
		),
		cnf.Apply(
			NewQueriesCommand(config, logger), viper.New(),
			cnf.WithConfig(config),
			cnf.WithGrafana(),
			cnf.WithGraphiteMetrics(),
			cnf.WithOutputFormat(),
		),
	)

	return &command
}
