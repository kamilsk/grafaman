package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"

	"github.com/kamilsk/grafaman/internal/presenter"
)

// New returns the new root command.
func New() *cobra.Command {
	var (
		format  string
		printer = new(presenter.Printer)
	)
	command := cobra.Command{
		Use:   "grafaman",
		Short: "Metrics coverage reporter",
		Long:  "Metrics coverage reporter for Graphite and Grafana.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := printer.SetOutput(cmd.OutOrStdout()).SetFormat(format); err != nil {
				return err
			}

			config := viper.New()
			config.SetConfigFile(viper.GetString("config"))
			config.SetConfigType("dotenv")
			err := config.ReadInConfig()
			if err == nil {
				return viper.MergeConfigMap(config.AllSettings())
			}

			if os.IsNotExist(err) {
				config.SetConfigFile("app.toml")
				config.SetConfigType("toml")
				if err := config.ReadInConfig(); err == nil {
					return viper.MergeConfigMap(config.Sub("envs.local.env_vars").AllSettings())
				}
				err = nil
			}
			return err
		},
		SilenceErrors: false,
		SilenceUsage:  true,
	}
	command.AddCommand(
		NewCoverageCommand(printer),
		NewMetricsCommand(printer),
		NewQueriesCommand(printer),
	)
	flags := command.PersistentFlags()
	flags.StringVarP(&format, "format", "f", printer.DefaultFormat(), "output format")
	flags.String("env-file", ".env.paas", "read in a file of environment variables; fallback to app.toml")
	fn.Must(
		func() error { return viper.BindPFlag("config", flags.Lookup("env-file")) },
		func() error {
			viper.RegisterAlias("grafana", "grafana_url")
			return viper.BindEnv("grafana", "GRAFANA_URL")
		},
		func() error {
			viper.RegisterAlias("dashboard", "grafana_dashboard")
			return viper.BindEnv("dashboard", "GRAFANA_DASHBOARD")
		},
		func() error {
			viper.RegisterAlias("graphite", "graphite_url")
			return viper.BindEnv("graphite", "GRAPHITE_URL")
		},
		func() error {
			viper.RegisterAlias("metrics", "graphite_metrics")
			return viper.BindEnv("metrics", "GRAPHITE_METRICS")
		},
	)
	return &command
}
