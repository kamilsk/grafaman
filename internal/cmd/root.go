package cmd

import (
	"os"

	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"
)

// TODO:debt
//  - support tabular view (for `| column -t`) to output analyze
//  - support json view to output analyze by jq

// New returns the new root command.
func New() *cobra.Command {
	const (
		formatDefault     = "default"
		formatCompact     = "compact"
		formatCompactLite = "compact-lite"
		formatMarkdown    = "markdown"
		formatRounded     = "rounded"
		formatUnicode     = "unicode"
	)
	var (
		format  string
		formats = []string{
			formatDefault,
			formatCompact,
			formatCompactLite,
			formatMarkdown,
			formatRounded,
			formatUnicode,
		}
		valid = map[string]*simpletable.Style{
			formatDefault:     simpletable.StyleDefault,
			formatCompact:     simpletable.StyleCompact,
			formatCompactLite: simpletable.StyleCompactLite,
			formatMarkdown:    simpletable.StyleMarkdown,
			formatRounded:     simpletable.StyleRounded,
			formatUnicode:     simpletable.StyleUnicode,
		}
		style simpletable.Style
	)
	command := cobra.Command{
		Use:   "grafaman",
		Short: "Metrics coverage reporter",
		Long:  "Metrics coverage reporter for Graphite and Grafana.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			selected, is := valid[format]
			if !is {
				return errors.Errorf("invalid format %q, only %v are available", format, formats)
			}
			style = *selected

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
		NewCoverageCommand(&style),
		NewMetricsCommand(&style),
		NewQueriesCommand(&style),
	)
	flags := command.PersistentFlags()
	flags.StringVarP(&format, "format", "f", formatDefault, "output format")
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
