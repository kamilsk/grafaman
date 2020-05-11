package cmd

import (
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"go.octolab.org/fn"
	xtime "go.octolab.org/time"
)

// NewCoverageCommand returns command to calculate metrics coverage by queries.
// TODO
//  - implement auth, if needed
//  - validate subset by regexp
//  - try to fetch fast by ~, if possible
//  - support exclude option
//  - support last option
func NewCoverageCommand(style *simpletable.Style) *cobra.Command {
	var (
		grafana   string
		dashboard string
		graphite  string
		subset    string
		last      time.Duration
		exclude   []string
	)
	command := cobra.Command{
		Use:   "coverage",
		Short: "calculates metrics coverage",
		Long:  "Calculates metrics coverage by queries.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVar(&grafana, "grafana", "", "Grafana API endpoint.")
	flags.StringVarP(&dashboard, "dashboard", "d", "", "A dashboard unique identifier.")
	flags.StringVar(&graphite, "graphite", "", "Graphite API endpoint.")
	flags.StringVarP(&subset, "subset", "s", "", "The required subset of metrics. Must be a simple prefix.")
	flags.DurationVar(&last, "last", xtime.Week, "The last interval to fetch.")
	flags.StringArrayVarP(&exclude, "exclude", "e", nil, "Patterns to exclude metrics from coverage, e.g. *.median")
	fn.Must(
		func() error { return command.MarkFlagRequired("grafana") },
		func() error { return command.MarkFlagRequired("dashboard") },
		func() error { return command.MarkFlagRequired("graphite") },
		func() error { return command.MarkFlagRequired("subset") },
	)
	return &command
}
