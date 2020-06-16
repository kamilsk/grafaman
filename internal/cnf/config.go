package cnf

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
