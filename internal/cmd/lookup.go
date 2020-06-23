package cmd

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"

	"github.com/kamilsk/grafaman/internal/cache"
	"github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/validator"
)

// NewCacheLookupCommand returns command to lookup cache.
func NewCacheLookupCommand(
	config *cnf.Config,
	logger *logrus.Logger,
) *cobra.Command {
	command := cobra.Command{
		Use:   "cache-lookup",
		Short: "lookup cache location",
		Long:  "Lookup cache location.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			flags := cmd.Flags()
			fn.Must(
				func() error { return viper.BindPFlag("graphite_metrics", flags.Lookup("metrics")) },
				func() error { return viper.Unmarshal(config) },
			)

			checker := validator.Metric()
			if !checker(config.Graphite.Prefix) {
				return errors.Errorf(
					"invalid metric prefix: %s; it must be simple, e.g. apps.services.name",
					config.Graphite.Prefix,
				)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			wrapper := cache.WrapGraphiteProvider(nil, afero.NewOsFs(), logger)
			cmd.Println(wrapper.Key(config.Graphite.Prefix))
			return nil
		},
	}

	flags := command.Flags()
	{
		flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")
	}

	return &command
}
