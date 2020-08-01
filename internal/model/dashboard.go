package model

import (
	"strings"

	"github.com/go-graphite/carbonapi/pkg/parser"
	"github.com/pkg/errors"
)

type Config struct {
	SkipRaw        bool
	SkipDuplicates bool
	NeedSorting    bool
	Unpack         bool
	TrimPrefixes   []string
}

type Dashboard struct {
	Prefix    string
	RawData   []Query
	Variables []Variable
}

func (dashboard *Dashboard) Queries(cfg Config) (Queries, error) {
	transformed := make(Queries, 0, len(dashboard.RawData))

	for _, raw := range dashboard.RawData {
		if dashboard.Prefix != "" && !strings.Contains(string(raw), dashboard.Prefix) {
			continue
		}

		if cfg.SkipRaw {
			transformed = append(transformed, raw)
			continue
		}

		exp, _, err := parser.ParseExpr(string(raw))
		if err != nil {
			return nil, errors.Wrapf(err, "dashboard: parse expression %q", raw)
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
		transformed.Sort()
	}
	return transformed, nil
}

func unpack(metric string, variables []Variable) []string {
	for _, variable := range variables {
		env := "$" + variable.Name
		if !strings.Contains(metric, env) {
			continue
		}
		// simplify logic: replace a variable by wildcard
		// motivation: variable can use dynamic source and that fact increase the complexity of the algorithm
		metric = strings.ReplaceAll(metric, env, "*")
	}
	return []string{metric}
}

type Variable struct {
	Name    string
	Options []Option
}

type Option struct {
	Name  string
	Value string
}
