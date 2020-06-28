package cnf

import "go.octolab.org/toolkit/config"

// Features defines a list of available features.
var Features = config.Features{
	{
		Name:    "cache",
		Enabled: true,
	},
	{
		Name:    "paas",
		Enabled: true,
	},
	{
		Name:    "repl",
		Enabled: true,
	},
}
