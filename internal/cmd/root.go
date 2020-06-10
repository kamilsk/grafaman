package cmd

import (
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"

	"github.com/kamilsk/grafaman/internal/presenter"
)

// New returns the new root command.
func New() *cobra.Command {
	var (
		debug   bool
		format  string
		logger  = logrus.New()
		printer = new(presenter.Printer)
	)
	command := cobra.Command{
		Use:   "grafaman",
		Short: "metrics coverage reporter",
		Long:  "Metrics coverage reporter for Graphite and Grafana.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := printer.SetOutput(cmd.OutOrStdout()).SetFormat(format); err != nil {
				return err
			}

			logger.SetOutput(ioutil.Discard)
			if debug {
				logger.SetLevel(logrus.DebugLevel)
				logger.SetOutput(cmd.ErrOrStderr())
				go func() {
					logger.Warning("start listen and serve pprof at http://localhost:8888/debug/pprof/")
					logger.Fatal(http.ListenAndServe(":8888", http.DefaultServeMux))
				}()
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
				if err, sub := config.ReadInConfig(), config.Sub("envs.local.env_vars"); err == nil && sub != nil {
					return viper.MergeConfigMap(sub.AllSettings())
				}
				err = nil
			}
			return err
		},
		SilenceErrors: false,
		SilenceUsage:  true,
	}
	command.AddCommand(
		NewCoverageCommand(logger, printer),
		NewMetricsCommand(logger, printer),
		NewQueriesCommand(logger, printer),
	)
	flags := command.PersistentFlags()
	flags.BoolVar(&debug, "debug", false, "enable debug")
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
