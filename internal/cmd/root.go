package cmd

import (
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"
	"go.octolab.org/toolkit/cli/debugger"

	"github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/presenter"
)

// New returns the new root command.
func New() *cobra.Command {
	var (
		debug   bool
		host    string
		format  string
		verbose int
		logger  = logrus.New()
		config  = new(cnf.Config)
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
				logger.SetOutput(cmd.ErrOrStderr())
				switch {
				case verbose == 1:
					logger.SetLevel(logrus.WarnLevel)
				case verbose == 2:
					logger.SetLevel(logrus.InfoLevel)
				case verbose > 2:
					logger.SetLevel(logrus.DebugLevel)
				default:
					logrus.SetLevel(logrus.ErrorLevel)
				}

				d, err := debugger.New(debugger.WithSpecificHost(host))
				if err != nil {
					return err
				}
				host, _ := d.Debug(func(err error) { logger.WithError(err).Fatal("run debugger") })
				logger.Warningf("start listen and serve pprof at http://%s/debug/pprof/", host)
			}

			cfg := viper.New()
			cfg.SetConfigFile(viper.GetString("config"))
			cfg.SetConfigType("dotenv")
			err := cfg.ReadInConfig()
			if err == nil {
				if !cfg.InConfig("graphite_metrics") && cfg.InConfig("app_name") {
					cfg.Set("graphite_metrics", "apps.services."+cfg.GetString("app_name"))
				}
				return viper.MergeConfigMap(cfg.AllSettings())
			}

			if os.IsNotExist(err) {
				cfg.SetConfigFile("app.toml")
				cfg.SetConfigType("toml")
				if err, sub := cfg.ReadInConfig(), cfg.Sub("envs.local.env_vars"); err == nil && sub != nil {
					if !sub.InConfig("graphite_metrics") && cfg.InConfig("name") {
						sub.Set("graphite_metrics", "apps.services."+cfg.GetString("name"))
					}
					return viper.MergeConfigMap(sub.AllSettings())
				}
				err = nil
			}
			return err
		},

		SilenceErrors: false,
		SilenceUsage:  true,
	}

	flags := command.PersistentFlags()
	{
		flags.String("env-file", ".env.paas", "read in a file of environment variables; fallback to app.toml")
	}
	flags.BoolVar(&debug, "debug", false, "enable debug")
	flags.StringVar(&host, "debug-host", "localhost:", "specific debug host")
	flags.StringVarP(&format, "format", "f", presenter.DefaultFormat, "output format")
	flags.CountVarP(&verbose, "verbose", "v", "increase the verbosity of messages if debug enabled")

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
		func() error {
			viper.RegisterAlias("name", "app_name")
			return viper.BindEnv("name", "APP_NAME")
		},
	)

	command.AddCommand(
		NewCacheLookupCommand(config, logger),
		NewCoverageCommand(config, logger, printer),
		NewMetricsCommand(config, logger, printer),
		NewQueriesCommand(config, logger, printer),
	)

	return &command
}
