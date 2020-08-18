package cmd

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.octolab.org/fn"

	"github.com/kamilsk/grafaman/internal/cnf"
	"github.com/kamilsk/grafaman/internal/model"
	"github.com/kamilsk/grafaman/internal/provider/graphite/cache"
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

			if config.Graphite.Prefix == "" {
				return errors.New("please provide metric prefix")
			}
			if !model.Metric(config.Graphite.Prefix).Valid() {
				return errors.Errorf(
					"invalid metric prefix: %s; it must be simple, e.g. apps.services.name",
					config.Graphite.Prefix,
				)
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(cache.Filename(config.Graphite.Prefix))
		},
	}

	flags := command.Flags()
	flags.StringP("metrics", "m", "", "the required subset of metrics (must be a simple prefix)")

	return &command
}
