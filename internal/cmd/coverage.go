package cmd

import (
	"github.com/alexeyco/simpletable"
	"github.com/spf13/cobra"
)

// NewCoverageCommand returns command to calculate metrics coverage by targets.
func NewCoverageCommand(style *simpletable.Style) *cobra.Command {
	command := cobra.Command{
		Use:   "coverage",
		Short: "calculates metrics coverage",
		Long:  "Calculates metrics coverage by targets.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return &command
}
