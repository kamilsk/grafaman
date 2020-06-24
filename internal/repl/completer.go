package repl

import (
	"github.com/c-bata/go-prompt"

	"github.com/kamilsk/grafaman/internal/provider"
)

func NewMetricsCompleter(metrics provider.Metrics) func(prompt.Document) []prompt.Suggest {
	return func(document prompt.Document) []prompt.Suggest {
		return nil
	}
}
