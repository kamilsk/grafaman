package cnf

import (
	"strings"

	"github.com/kamilsk/grafaman/internal/model"
)

// Config contains all necessary tool configuration.
type Config struct {
	Grafana struct {
		URL       string `mapstructure:"grafana"`
		Dashboard string `mapstructure:"dashboard"`
	} `mapstructure:",squash"`
	Graphite struct {
		URL    string `mapstructure:"graphite"`
		Filter string `mapstructure:"filter"`
		Prefix string `mapstructure:"metrics"`
	} `mapstructure:",squash"`
}

// FilterQuery returns a Query to filter metrics.
func (config Config) FilterQuery() model.Query {
	filter, prefix := config.Graphite.Filter, config.Graphite.Prefix
	if filter == "" {
		filter = "*"
	}
	if !strings.HasPrefix(filter, prefix) {
		filter = prefix + "." + filter
	}
	return model.Query(filter)
}
