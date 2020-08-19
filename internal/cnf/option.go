package cnf

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"

	"github.com/kamilsk/grafaman/internal/presenter"
)

const ConfigKey = "config"

func Apply(command *cobra.Command, provider *viper.Viper, options ...Option) *cobra.Command {
	for _, configure := range options {
		configure(command, provider)
	}
	return command
}

type Option func(command *cobra.Command, provider *viper.Viper)

func WithConfig(config *Config) Option {
	return func(command *cobra.Command, provider *viper.Viper) {
		next := command.PreRunE
		if next == nil {
			next = func(cmd *cobra.Command, args []string) error { return nil }
		}
		command.PreRunE = func(cmd *cobra.Command, args []string) error {
			cfg := viper.New()
			cfg.SetConfigFile(config.File)
			cfg.SetConfigType("dotenv")
			switch err := cfg.ReadInConfig(); true {
			case err == nil:
				fn.Must(func() error { return provider.MergeConfigMap(cfg.AllSettings()) })
			case os.IsNotExist(err):
				cfg.SetConfigFile("app.toml")
				cfg.SetConfigType("toml")
				if err, sub := cfg.ReadInConfig(), cfg.Sub("envs.local.env_vars"); err == nil && sub != nil {
					fn.Must(func() error { return provider.MergeConfigMap(sub.AllSettings()) })
				}
				err = nil // ignore, if implicit fallback doesn't work
			}

			fn.Must(func() error { return provider.Unmarshal(config) })
			return next(cmd, args)
		}
	}
}

func WithDebug() Option {
	return func(command *cobra.Command, provider *viper.Viper) {
		flags := command.Flags()
		flags.Bool("debug", false, "enable debug")
		flags.String("debug-host", "localhost:", "specific debug host")
		flags.CountP("verbose", "v", "increase the verbosity of messages if debug enabled")

		fn.Must(
			func() error { return provider.BindPFlag("debug.enabled", flags.Lookup("debug")) },
			func() error { return provider.BindPFlag("debug.host", flags.Lookup("debug-host")) },
			func() error { return provider.BindPFlag("debug.level", flags.Lookup("verbose")) },
		)
	}
}

func WithGrafana() Option {
	return func(command *cobra.Command, provider *viper.Viper) {
		flags := command.Flags()
		flags.String("grafana", "", "Grafana API endpoint")
		flags.StringP("dashboard", "d", "", "a dashboard unique identifier")

		provider.RegisterAlias("grafana", "grafana_url")
		provider.RegisterAlias("dashboard", "grafana_dashboard")

		fn.Must(
			func() error { return provider.BindPFlag("grafana_url", flags.Lookup("grafana")) },
			func() error { return provider.BindPFlag("grafana_dashboard", flags.Lookup("dashboard")) },
			func() error { return provider.BindEnv("grafana", "GRAFANA_URL") },
			func() error { return provider.BindEnv("dashboard", "GRAFANA_DASHBOARD") },
		)
	}
}

func WithGraphite() Option {
	return func(command *cobra.Command, provider *viper.Viper) {
		flags := command.Flags()
		flags.StringP("graphite", "e", "", "Graphite API endpoint")
		flags.String("filter", "", "query to filter metrics, e.g. some.*.metric")

		provider.RegisterAlias("graphite", "graphite_url")

		fn.Must(
			func() error { return provider.BindPFlag("graphite_url", flags.Lookup("graphite")) },
			func() error { return provider.BindPFlag("filter", flags.Lookup("filter")) },
			func() error { return provider.BindEnv("graphite", "GRAPHITE_URL") },
		)

		WithGraphiteMetrics()(command, provider)
	}
}

func WithGraphiteMetrics() Option {
	return func(command *cobra.Command, provider *viper.Viper) {
		flags := command.Flags()
		flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")

		provider.RegisterAlias("app_name", "app")
		provider.RegisterAlias("name", "app")
		provider.RegisterAlias("metrics", "graphite_metrics")

		fn.Must(
			func() error { return provider.BindPFlag("graphite_metrics", flags.Lookup("metrics")) },
			func() error { return provider.BindEnv("app", "APP_NAME") },
			func() error { return provider.BindEnv("metrics", "GRAPHITE_METRICS") },
		)
	}
}

func WithOutputFormat() Option {
	return func(command *cobra.Command, provider *viper.Viper) {
		flags := command.Flags()
		flags.StringP("format", "f", presenter.DefaultFormat, "output format")
	}
}
