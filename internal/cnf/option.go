package cnf

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"
	"go.octolab.org/toolkit/cli/debugger"

	"github.com/kamilsk/grafaman/internal/presenter"
)

// After inserts a new function into the pointer, which calls the self function before and the last after.
func After(pointer *func(*cobra.Command, []string), last func(*cobra.Command, []string)) {
	first := *pointer
	if first == nil {
		first = func(*cobra.Command, []string) {}
	}
	*pointer = func(command *cobra.Command, args []string) {
		first(command, args)
		last(command, args)
	}
}

// AfterE inserts a new function into the pointer, which calls the self function before and the last after.
func AfterE(pointer *func(*cobra.Command, []string) error, last func(*cobra.Command, []string) error) {
	first := *pointer
	if first == nil {
		first = func(*cobra.Command, []string) error { return nil }
	}
	*pointer = func(command *cobra.Command, args []string) error {
		if err := first(command, args); err != nil {
			return err
		}
		return last(command, args)
	}
}

// Apply applies options to the Command.
func Apply(command *cobra.Command, container *viper.Viper, options ...Option) *cobra.Command {
	for _, configure := range options {
		configure(command, container)
	}
	return command
}

// Before inserts a new function into the pointer, which calls the first function before and the self after.
func Before(pointer *func(*cobra.Command, []string), first func(*cobra.Command, []string)) {
	last := *pointer
	if last == nil {
		last = func(*cobra.Command, []string) {}
	}
	*pointer = func(command *cobra.Command, args []string) {
		first(command, args)
		last(command, args)
	}
}

// BeforeE inserts a new function into the pointer, which calls the first function before and the self after.
func BeforeE(pointer *func(*cobra.Command, []string) error, first func(*cobra.Command, []string) error) {
	last := *pointer
	if last == nil {
		last = func(*cobra.Command, []string) error { return nil }
	}
	*pointer = func(command *cobra.Command, args []string) error {
		if err := first(command, args); err != nil {
			return err
		}
		return last(command, args)
	}
}

// An Option is a Command configuration function.
type Option func(*cobra.Command, *viper.Viper)

// WithConfig returns an Option to inject configuration from a container and config files into the Config.
func WithConfig(config *Config) Option {
	return func(command *cobra.Command, container *viper.Viper) {
		BeforeE(&command.PreRunE, func(cmd *cobra.Command, args []string) error {
			cfg := viper.New()
			cfg.SetConfigFile(config.File)
			cfg.SetConfigType("dotenv")
			switch err := cfg.ReadInConfig(); {
			case err == nil:
				fn.Must(func() error { return container.MergeConfigMap(cfg.AllSettings()) })
			case os.IsNotExist(err):
				cfg.SetConfigFile("app.toml")
				cfg.SetConfigType("toml")
				if err, sub := cfg.ReadInConfig(), cfg.Sub("envs.local.env_vars"); err == nil && sub != nil {
					fn.Must(
						func() error {
							return container.MergeConfigMap(map[string]interface{}{"app_name": cfg.GetString("name")})
						},
						func() error { return container.MergeConfigMap(sub.AllSettings()) },
					)
				}
			}

			fn.Must(func() error { return container.Unmarshal(config) })

			// ad hoc
			if config.Graphite.Prefix == "" && config.App != "" {
				config.Graphite.Prefix = fmt.Sprintf("apps.services.%s", config.App)
			}

			return nil
		})
	}
}

// WithDebug returns an Option to inject debugger and configure the logger.
func WithDebug(config *Config, logger *logrus.Logger) Option {
	return func(command *cobra.Command, container *viper.Viper) {
		flags := command.Flags()
		flags.Bool("debug", false, "enable debug")
		flags.String("debug-host", "localhost:", "specific debug host")
		flags.CountP("verbose", "v", "increase the verbosity of messages if debug enabled")

		fn.Must(
			func() error { return container.BindPFlag("debug.enabled", flags.Lookup("debug")) },
			func() error { return container.BindPFlag("debug.host", flags.Lookup("debug-host")) },
			func() error { return container.BindPFlag("debug.level", flags.Lookup("verbose")) },
		)

		BeforeE(&command.PreRunE, func(cmd *cobra.Command, args []string) error {
			logger.SetOutput(ioutil.Discard)
			if config.Debug.Enabled {
				logger.SetOutput(cmd.ErrOrStderr())
				switch verbose := config.Debug.Level; {
				case verbose == 1:
					logger.SetLevel(logrus.WarnLevel)
				case verbose == 2:
					logger.SetLevel(logrus.InfoLevel)
				case verbose > 2:
					logger.SetLevel(logrus.DebugLevel)
				default:
					logger.SetLevel(logrus.ErrorLevel)
				}

				d, err := debugger.New(debugger.WithSpecificHost(config.Debug.Host))
				if err != nil {
					return err
				}
				host, _ := d.Debug(func(err error) { logger.WithError(err).Fatal("run debugger") })
				logger.Warningf("start listen and serve pprof at http://%s/debug/pprof/", host)
			}

			return nil
		})
	}
}

// WithGrafana returns an Option to inject flags related to Grafana configuration.
func WithGrafana() Option {
	return func(command *cobra.Command, container *viper.Viper) {
		flags := command.Flags()
		flags.String("grafana", "", "Grafana API endpoint")
		flags.StringP("dashboard", "d", "", "a dashboard unique identifier")

		container.RegisterAlias("grafana", "grafana_url")
		container.RegisterAlias("dashboard", "grafana_dashboard")

		fn.Must(
			func() error { return container.BindPFlag("grafana_url", flags.Lookup("grafana")) },
			func() error { return container.BindPFlag("grafana_dashboard", flags.Lookup("dashboard")) },
			func() error { return container.BindEnv("grafana", "GRAFANA_URL") },
			func() error { return container.BindEnv("dashboard", "GRAFANA_DASHBOARD") },
		)
	}
}

// WithGraphite returns an Option to inject flags related to Graphite configuration.
func WithGraphite() Option {
	return func(command *cobra.Command, container *viper.Viper) {
		flags := command.Flags()
		flags.StringP("graphite", "e", "", "Graphite API endpoint")
		flags.String("filter", "", "query to filter metrics, e.g. some.*.metric")

		container.RegisterAlias("graphite", "graphite_url")

		fn.Must(
			func() error { return container.BindPFlag("graphite_url", flags.Lookup("graphite")) },
			func() error { return container.BindPFlag("filter", flags.Lookup("filter")) },
			func() error { return container.BindEnv("graphite", "GRAPHITE_URL") },
		)

		WithGraphiteMetrics()(command, container)
	}
}

// WithGraphiteMetrics returns an Option to inject flags related to Graphite configuration.
func WithGraphiteMetrics() Option {
	return func(command *cobra.Command, container *viper.Viper) {
		flags := command.Flags()
		flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")

		container.RegisterAlias("app_name", "app")
		container.RegisterAlias("name", "app")
		container.RegisterAlias("metrics", "graphite_metrics")

		fn.Must(
			func() error { return container.BindPFlag("graphite_metrics", flags.Lookup("metrics")) },
			func() error { return container.BindEnv("app", "APP_NAME") },
			func() error { return container.BindEnv("metrics", "GRAPHITE_METRICS") },
		)
	}
}

// WithOutputFormat returns an Option to inject flags related to output format.
func WithOutputFormat() Option {
	return func(command *cobra.Command, container *viper.Viper) {
		flags := command.Flags()
		flags.StringP("format", "f", presenter.DefaultFormat, "output format")

		fn.Must(func() error { return container.BindPFlag("output.format", flags.Lookup("format")) })
	}
}
