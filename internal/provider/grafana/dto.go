package grafana

import "github.com/kamilsk/grafaman/internal/model"

type dashboard struct {
	Panels     []panel    `json:"panels,omitempty"`
	Templating templating `json:"templating,omitempty"`
}

type panel struct {
	ID      int      `json:"id,omitempty"`
	Title   string   `json:"title,omitempty"`
	Type    string   `json:"type,omitempty"`
	Panels  []panel  `json:"panels,omitempty"`
	Targets []target `json:"targets,omitempty"`
}

type templating struct {
	List []variable `json:"list,omitempty"`
}

type target struct {
	Query string `json:"target,omitempty"`
}

type variable struct {
	Name    string        `json:"name,omitempty"`
	Options []option      `json:"options,omitempty"`
	Current currentOption `json:"current,omitempty"`
}

type option struct {
	Text  string `json:"text,omitempty"`
	Value string `json:"value,omitempty"`
}

type currentOption struct {
	Text  string      `json:"text,omitempty"`
	Value interface{} `json:"value,omitempty"` // string or []string
}

func convertTargets(in []target) []model.Query {
	out := make([]model.Query, 0, len(in))

	for _, target := range in {
		if target.Query != "" {
			out = append(out, model.Query(target.Query))
		}
	}

	return out
}

func fetchTargets(panels []panel) []target {
	targets := make([]target, 0, 4*len(panels))

	for _, panel := range panels {
		if count := len(panel.Targets); count > 0 {
			targets = append(targets, panel.Targets...)
			continue
		}
		targets = append(targets, fetchTargets(panel.Panels)...)
	}

	return targets
}

func convertVariables(in []variable) []model.Variable {
	out := make([]model.Variable, 0, len(in))

	for _, v := range in {
		variable := model.Variable{Name: v.Name, Options: make([]model.Option, 0, len(v.Options))}
		for _, opt := range v.Options {
			variable.Options = append(variable.Options, model.Option{Name: opt.Text, Value: opt.Value})
		}
		out = append(out, variable)
	}

	return out
}

func fetchVariables(dashboard dashboard) []variable {
	variables := make([]variable, 0, 4*len(dashboard.Templating.List))

	for _, variable := range dashboard.Templating.List {
		// handle current option
		switch v := variable.Current.Value.(type) {
		case []interface{}:
			for _, opt := range v {
				value, _ := opt.(string)
				variable.Options = append(variable.Options, option{Text: variable.Current.Text, Value: value})
			}
		case string:
			variable.Options = append(variable.Options, option{Text: variable.Current.Text, Value: v})
		}

		// filter duplicate options
		filtered := variable.Options[:0]
		registry := map[string]struct{}{}
		for _, option := range variable.Options {
			if _, present := registry[option.Value]; present {
				continue
			}
			registry[option.Value] = struct{}{}
			filtered = append(filtered, option)
		}
		variable.Options = filtered
		variables = append(variables, variable)
	}

	return variables
}
