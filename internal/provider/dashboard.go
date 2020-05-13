package provider

import (
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/go-graphite/carbonapi/pkg/parser"
	"github.com/pkg/errors"
)

type Dashboard struct {
	Subset    string
	RawData   []Query
	Variables []Variable
}

func (dashboard *Dashboard) Queries(cfg Transform) (Queries, error) {
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
			return nil, errors.Wrapf(err, "parse expression %q", raw)
		}
		for _, query := range exp.Metrics() {
			for _, prefix := range cfg.TrimPrefixes {
				if strings.HasPrefix(query.Metric, prefix) {
					query.Metric = strings.TrimPrefix(query.Metric, prefix)
					break
				}
			}
			queries := Queries{Query(query.Metric)}
			if cfg.Unpack && strings.Contains(query.Metric, "$") {
				queries.Convert(unpack(query.Metric, dashboard.Variables))
			}
			transformed = append(transformed, queries...)
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

func unpack(metric string, variables []Variable) []string {
	for i, variable := range variables {
		env := "$" + variable.Name
		if !strings.Contains(metric, env) {
			continue
		}
		result := make([]string, 0, len(variable.Options))
		for _, option := range variable.Options {
			// TODO:research some variable has $__all option
			result = append(result, unpack(strings.ReplaceAll(metric, env, option.Value), variables[i+1:])...)
		}
		return result
	}
	return []string{metric}
}

type Transform struct {
	SkipRaw        bool
	SkipDuplicates bool
	NeedSorting    bool
	Unpack         bool
	TrimPrefixes   []string
}

type Query string

type Queries []Query

func (queries *Queries) Convert(src []string) {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	*queries = *(*[]Query)(unsafe.Pointer(header))
}

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
