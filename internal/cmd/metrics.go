package cmd

import (
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"go.octolab.org/fn"
)

// NewMetricsCommand returns command to fetch metrics from Graphite.
func NewMetricsCommand(style *simpletable.Style) *cobra.Command {
	var (
		endpoint string
		subset   string
		collapse int
	)
	command := cobra.Command{
		Use:   "metrics",
		Short: "fetch metrics from Graphite",
		Long:  "Fetch metrics from Graphite.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO
			//  - auth
			//  - validate subset by regexp
			//  - try to fetch fast
			//  - operates by nodes instead of strings
			//  - sort by id, cause async fetching
			//  - implement collapse mechanism
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVarP(&endpoint, "endpoint", "e", "", "Graphite API endpoint.")
	flags.StringVarP(&subset, "subset", "s", "", "The required subset of metrics. Must be a simple prefix.")
	flags.IntVarP(&collapse, "collapse", "c", 0, "How many levels from the right to collapse by wildcard.")
	fn.Must(
		func() error { return command.MarkFlagRequired("endpoint") },
		func() error { return command.MarkFlagRequired("subset") },
	)
	return &command
}
