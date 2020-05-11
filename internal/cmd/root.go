package cmd

import (
	"github.com/alexeyco/simpletable"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

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
			return nil
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
	flags.StringVarP(&format, "format", "f", formatDefault, "Output format.")
	return &command
}
