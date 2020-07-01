package repl

import (
	"sort"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/gobwas/glob"

	"github.com/kamilsk/grafaman/internal/model"
)

// NewMetricsCompleter returns naive implementation to autocomplete user input.
func NewMetricsCompleter(prefix string, metrics model.Metrics) func(prompt.Document) []prompt.Suggest {
	prefix += "."
	return func(document prompt.Document) []prompt.Suggest {
		input := document.TextBeforeCursor()

		pattern := input
		if !strings.HasPrefix(pattern, prefix) {
			pattern = prefix + pattern
		}
		if !strings.HasSuffix(pattern, "*") {
			pattern += "*"
		}
		matcher := glob.MustCompile(pattern)

		segmentsInInput := strings.Count(input, ".")
		segmentsInFullPattern := strings.Count(pattern, ".")

		registry := make(map[string]struct{})
		for _, metric := range metrics {
			metric := string(metric)
			if matcher.Match(metric) {
				// input: s*
				//  - 0 segments
				//  - SplitAfter -> [s*]
				//  - Join -> ""
				// input: some.specific.m*
				//  - 2 segments [some, specific]
				//  - SplitAfter -> [some., specific., m*]
				//  - Join -> some.specific.
				suggestion := strings.Join(strings.SplitAfter(input, ".")[:segmentsInInput], "")
				// metric: prefix.some.specific.metric.name.and.value, pattern: prefix.some.specific.m*
				//  - 3 segments [prefix, some, specific]
				//  - SplitAfterN -> [prefix., some., specific., metric., name.and.value]
				//  - Index -> metric.
				suggestion += strings.SplitAfterN(metric, ".", segmentsInFullPattern+2)[segmentsInFullPattern]
				// result: some.specific.metric.
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
