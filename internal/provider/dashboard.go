package provider

import (
	"sort"
	"strings"

	"github.com/go-graphite/carbonapi/pkg/parser"
	"github.com/pkg/errors"
)

type Dashboard struct {
	Subset    string
	RawData   []Query
	Variables []Variable
}

func (dashboard Dashboard) Queries(cfg Transform) (Queries, error) {
	transformed := make(Queries, 0, len(dashboard.RawData))

	for _, raw := range dashboard.RawData {
		if dashboard.Subset != "" && !strings.Contains(string(raw), dashboard.Subset) {
			continue
		}

		if cfg.SkipRaw {
			transformed = append(transformed, raw)
			continue
		}

		exp, _, err := parser.ParseExpr(string(raw))
		if err != nil {
			return nil, errors.Wrap(err, "parse query")
		}
		for _, query := range exp.Metrics() {
			for _, prefix := range cfg.TrimPrefixes {
				if strings.HasPrefix(query.Metric, prefix) {
					query.Metric = strings.TrimPrefix(query.Metric, prefix)
					break
				}
			}
			transformed = append(transformed, Query(query.Metric))
		}
	}

	if !cfg.SkipDuplicates {
		registry := map[Query]struct{}{}

		// preserve order
		iterator := transformed
		transformed = transformed[:0]
		for _, query := range iterator {
			if _, present := registry[query]; present {
				continue
			}
			transformed = append(transformed, query)
		}
	}

	if cfg.NeedSorting {
		sort.Sort(transformed)
	}
	return transformed, nil
}

type Transform struct {
	SkipRaw        bool
	SkipDuplicates bool
	NeedSorting    bool
	TrimPrefixes   []string
}

type Query string

type Queries []Query

func (queries Queries) Len() int           { return len(queries) }
func (queries Queries) Less(i, j int) bool { return queries[i] < queries[j] }
func (queries Queries) Swap(i, j int)      { queries[i], queries[j] = queries[j], queries[i] }

type Variable struct {
	Name    string
	Options []Option
}

type Option struct {
	Name  string
	Value string
}
