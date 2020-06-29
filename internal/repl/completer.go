package repl

import (
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/gobwas/glob"

	"github.com/kamilsk/grafaman/internal/provider"
)

// NewMetricsCompleter returns naive implementation to autocomplete user input.
func NewMetricsCompleter(prefix string, metrics provider.Metrics) func(prompt.Document) []prompt.Suggest {
	prefix += "."
	return func(document prompt.Document) []prompt.Suggest {
		origin := document.TextBeforeCursor()
		pattern := origin
		if !strings.HasPrefix(pattern, prefix) {
			pattern = prefix + pattern
		}
		if !strings.HasSuffix(pattern, "*") {
			pattern += "*"
		}

		matcher, err := glob.Compile(pattern)
		if err != nil {
			return nil
		}

		// TODO:refactoring better naming
		segmentO := strings.Count(origin, ".")
		segmentP := strings.Count(pattern, ".")

		registry := make(map[string]struct{})
		for _, metric := range metrics {
			metric := string(metric)
			if matcher.Match(metric) {
				suggestion := strings.Join(strings.SplitAfterN(origin, ".", segmentO+1)[:segmentO], "")
				suggestion += strings.SplitAfterN(metric, ".", segmentP+2)[segmentP]
				registry[suggestion] = struct{}{}
			}
		}
		if len(registry) == 0 {
			return nil
		}

		suggestions := make([]prompt.Suggest, 0, len(registry))
		for suggestion := range registry {
			suggestions = append(suggestions, prompt.Suggest{Text: suggestion})
		}
		sort.Sort(sortByText(suggestions))
		return suggestions
	}
}

type sortByText []prompt.Suggest

func (list sortByText) Len() int           { return len(list) }
func (list sortByText) Less(i, j int) bool { return list[i].Text < list[j].Text }
func (list sortByText) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }
