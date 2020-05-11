package cmd

import (
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
	"go.octolab.org/fn"
)

// NewTargetsCommand returns command to fetch targets from a Grafana dashboard.
// TODO
//  - implement auth, if needed
//  - validate subset by regexp
//  - support raw option
//  - support duplicates option
//  - support sort option
func NewTargetsCommand(style *simpletable.Style) *cobra.Command {
	var (
		endpoint   string
		dashboard  string
		subset     string
		raw        bool
		duplicates bool
		sort       bool
	)
	command := cobra.Command{
		Use:   "targets",
		Short: "fetch targets from a Grafana dashboard",
		Long:  "Fetch targets from a Grafana dashboard.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	flags := command.Flags()
	flags.StringVarP(&endpoint, "endpoint", "e", "", "Grafana API endpoint.")
	flags.StringVarP(&dashboard, "dashboard", "d", "", "A dashboard unique identifier.")
	flags.BoolVar(&raw, "raw", false, "Leave the original values of targets.")
	flags.BoolVar(&duplicates, "allow-duplicates", false, "Allow duplicates of targets.")
	flags.BoolVar(&sort, "sort", false, "Need to sort targets.")
	flags.StringVarP(&subset, "subset", "s", "", "The required subset of metrics. Must be a simple prefix.")
	fn.Must(
		func() error { return command.MarkFlagRequired("endpoint") },
	)
	return &command
}
