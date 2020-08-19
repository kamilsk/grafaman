package cmd

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

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
			if config.Graphite.Prefix == "" {
				return errors.New("please provide metric prefix")
			}
			if prefix := config.Graphite.Prefix; !model.Metric(prefix).Valid() {
				return errors.Errorf("invalid metric prefix: %s; it must be simple, e.g. apps.services.name", prefix)
			}
			return nil
		},

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(cache.Filename(config.Graphite.Prefix))
		},
	}

	return &command
}
