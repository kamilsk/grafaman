package cnf

import "go.octolab.org/toolkit/config"

// Features defines a list of available features.
var Features = config.Features{
	{
		Name:    "cache",
		Enabled: true,
		Brief:   "Store state of metrics between calls.",
		Docs:    "https://www.notion.so/octolab/Cache-layer-9766d3fbe07d4eb1808f6165fff8c5d0?r=0b753cbf767346f5a6fd51194829a2f3",
	},
	{
		Name:    "debug",
		Enabled: true,
		Brief:   "Enable logger and run server with pprof.",
		Docs:    "https://www.notion.so/octolab/Debugger-95877ccae2ee4168845c54066a443fe3?r=0b753cbf767346f5a6fd51194829a2f3",
	},
	{
		Name:    "grafonnet",
		Enabled: false,
		Brief:   "Declare dashboard as code.",
		Docs:    "https://www.notion.so/octolab/Grafonnet-3cd366ab76e146db82fd28f520bbdf68?r=0b753cbf767346f5a6fd51194829a2f3",
	},
	{
		Name:    "paas",
		Enabled: true,
		Brief:   "Support app.toml and .env.paas to fetch configuration.",
		Docs:    "https://www.notion.so/octolab/Avito-integration-066415cfcd914e20af7408e45272e98c?r=0b753cbf767346f5a6fd51194829a2f3",
	},
	{
		Name:    "repl",
		Enabled: true,
		Brief:   "Enable interactive REPL mode to analyze data.",
		Docs:    "https://www.notion.so/octolab/REPL-support-62f1301f41df4e8681527e268e10f539?r=0b753cbf767346f5a6fd51194829a2f3",
	},
}
