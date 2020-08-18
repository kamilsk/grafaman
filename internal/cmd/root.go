package cmd

import (
	"fmt"
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
		logger  = logrus.New()
		config  = new(cnf.Config)
		printer = new(presenter.Printer)
	)

	provider := viper.GetViper() // TODO:refactor issue#41

	command := cobra.Command{
		Use:   "grafaman",
		Short: "metrics coverage reporter",
		Long:  "Metrics coverage reporter for Graphite and Grafana.",

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			/* TODO:refactor issue#41
			if err := printer.SetOutput(cmd.OutOrStdout()).SetFormat(provider.GetString("output.format")); err != nil {
				return err
			}
			*/

			logger.SetOutput(ioutil.Discard)
			if provider.GetBool("debug.enabled") {
				logger.SetOutput(cmd.ErrOrStderr())
				switch verbose := provider.GetInt("debug.level"); true {
				case verbose == 1:
					logger.SetLevel(logrus.WarnLevel)
				case verbose == 2:
					logger.SetLevel(logrus.InfoLevel)
				case verbose > 2:
					logger.SetLevel(logrus.DebugLevel)
				default:
					logrus.SetLevel(logrus.ErrorLevel)
				}

				d, err := debugger.New(debugger.WithSpecificHost(provider.GetString("debug.host")))
				if err != nil {
					return err
				}
				host, _ := d.Debug(func(err error) { logger.WithError(err).Fatal("run debugger") })
				logger.Warningf("start listen and serve pprof at http://%s/debug/pprof/", host)
			}

			cfg := viper.New()
			cfg.SetConfigFile(provider.GetString("config"))
			cfg.SetConfigType("dotenv")
			err := cfg.ReadInConfig()
			if err == nil {
				if !cfg.InConfig("graphite_metrics") && cfg.InConfig("app_name") {
					cfg.Set("graphite_metrics", fmt.Sprintf("apps.services.%s", cfg.GetString("app_name")))
				}
				return provider.MergeConfigMap(cfg.AllSettings())
			}

			if os.IsNotExist(err) {
				cfg.SetConfigFile("app.toml")
				cfg.SetConfigType("toml")
				if err, sub := cfg.ReadInConfig(), cfg.Sub("envs.local.env_vars"); err == nil && sub != nil {
					if !sub.InConfig("graphite_metrics") && cfg.InConfig("name") {
						sub.Set("graphite_metrics", fmt.Sprintf("apps.services.%s", cfg.GetString("name")))
					}
					return provider.MergeConfigMap(sub.AllSettings())
				}
				err = nil
			}
			return err
		},

		SilenceErrors: false,
		SilenceUsage:  true,
	}

	flags := command.PersistentFlags()
	flags.String("env-file", ".env.paas", "read in a file of environment variables; fallback to app.toml")
	fn.Must(
		func() error { return provider.BindPFlag("config", flags.Lookup("env-file")) },
	)

	command.AddCommand(
		cnf.Apply(
			NewCacheLookupCommand(config, logger), provider,
			//
		),
		cnf.Apply(
			NewCoverageCommand(config, logger, printer), provider,
			cnf.WithDebug(), cnf.WithGrafana(), cnf.WithGraphite(), cnf.WithOutputFormat(),
		),
		cnf.Apply(
			NewMetricsCommand(config, logger, printer), provider,
			cnf.WithDebug(), cnf.WithGraphite(), cnf.WithOutputFormat(),
		),
		cnf.Apply(
			NewQueriesCommand(config, logger, printer), provider,
			cnf.WithGrafana(), cnf.WithOutputFormat(),
		),
	)

	return &command
}
