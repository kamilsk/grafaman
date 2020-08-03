package model

import (
	"strings"

	"github.com/go-graphite/carbonapi/pkg/parser"
	"github.com/pkg/errors"
)

// A Config contains configuration to process raw queries of a Dashboard.
type Config struct {
	SkipRaw        bool
	SkipDuplicates bool
	NeedSorting    bool
	Unpack         bool
}

// A Dashboard represents Grafana dashboard.
type Dashboard struct {
	Prefix    string
	RawData   []Query
	Variables []Variable
}

// Queries applies variables to raw queries to transform them.
func (dashboard *Dashboard) Queries(cfg Config) (Queries, error) {
	transformed := make(Queries, 0, len(dashboard.RawData))
	prefix := dashboard.Prefix

	for _, raw := range dashboard.RawData {
		if prefix != "" && !strings.Contains(string(raw), prefix) {
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
			if prefix != "" {
				if !strings.Contains(query.Metric, prefix) {
					continue
				}
				if !strings.HasPrefix(query.Metric, prefix) {
					query.Metric = query.Metric[strings.Index(query.Metric, prefix):]
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
			registry[query] = struct{}{}
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

// An Option represents a possible value of the Variable.
type Option struct {
	Name  string
	Value string
}

// A Variable represents a Grafana dashboard variable.
type Variable struct {
	Name    string
	Options []Option
}
