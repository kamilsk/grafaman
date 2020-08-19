package cnf

import (
	"strings"

	"github.com/kamilsk/grafaman/internal/model"
)

// A Config contains all necessary tool configuration.
type Config struct {
	App     string `mapstructure:"app"`
	File    string `mapstructure:"-"`
	Grafana struct {
		URL       string `mapstructure:"grafana"`
		Dashboard string `mapstructure:"dashboard"`
	} `mapstructure:",squash"`
	Graphite struct {
		URL    string `mapstructure:"graphite"`
		Filter string `mapstructure:"filter"`
		Prefix string `mapstructure:"metrics"`
	} `mapstructure:",squash"`
	Debug struct {
		Enabled bool   `mapstructure:"enabled"`
		Host    string `mapstructure:"host"`
		Level   int    `mapstructure:"level"`
	} `mapstructure:"debug"`
	Output struct {
		Format string `mapstructure:"format"`
	} `mapstructure:"output"`
}

// FilterQuery returns a Query to filter metrics.
func (config *Config) FilterQuery() model.Query {
	filter, prefix := config.Graphite.Filter, config.Graphite.Prefix
	if filter == "" {
		filter = "*"
	}
	if !strings.HasPrefix(filter, prefix) {
		filter = prefix + "." + filter
	}
	return model.Query(filter)
}
