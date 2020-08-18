package cnf

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"

	"github.com/kamilsk/grafaman/internal/presenter"
)

func Apply(command *cobra.Command, provider *viper.Viper, options ...Option) *cobra.Command {
	for _, configure := range options {
		configure(command, provider)
	}
	return command
}

type Option func(command *cobra.Command, provider *viper.Viper)

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
		fn.Must(
			func() error {
				provider.RegisterAlias("grafana", "grafana_url")
				return provider.BindEnv("grafana", "GRAFANA_URL")
			},
			func() error {
				provider.RegisterAlias("dashboard", "grafana_dashboard")
				return provider.BindEnv("dashboard", "GRAFANA_DASHBOARD")
			},
		)
	}
}

func WithGraphite() Option {
	return func(command *cobra.Command, provider *viper.Viper) {
		fn.Must(
			func() error {
				provider.RegisterAlias("graphite", "graphite_url")
				return provider.BindEnv("graphite", "GRAPHITE_URL")
			},
			func() error {
				provider.RegisterAlias("metrics", "graphite_metrics")
				return provider.BindEnv("metrics", "GRAPHITE_METRICS")
			},
			func() error {
				provider.RegisterAlias("name", "app_name")
				return provider.BindEnv("name", "APP_NAME")
			},
		)
	}
}

func WithOutputFormat() Option {
	return func(command *cobra.Command, provider *viper.Viper) {
		flags := command.Flags()
		flags.StringP("format", "f", presenter.DefaultFormat, "output format")
	}
}
