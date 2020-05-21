package grafana

import entity "github.com/kamilsk/grafaman/internal/provider"

type dashboard struct {
	Panels     []panel `json:"panels"`
	Templating struct {
		List []variable `json:"list"`
	} `json:"templating"`
}

type panel struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Type    string   `json:"type"`
	Panels  []panel  `json:"panels"`
	Targets []target `json:"targets"`
}

type target struct {
	Query string `json:"target"`
}

type variable struct {
	Name    string        `json:"name"`
	Options []option      `json:"options"`
	Current currentOption `json:"current"`
}

type option struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type currentOption struct {
	Text  string      `json:"text"`
	Value interface{} `json:"value"` // string or []string
}

func convertTargets(in []target) []entity.Query {
	out := make([]entity.Query, 0, len(in))

	for _, target := range in {
		if target.Query != "" {
			out = append(out, entity.Query(target.Query))
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

func convertVariables(in []variable) []entity.Variable {
	out := make([]entity.Variable, 0, len(in))

	for _, v := range in {
		variable := entity.Variable{Name: v.Name, Options: make([]entity.Option, 0, len(v.Options))}
		for _, opt := range v.Options {
			variable.Options = append(variable.Options, entity.Option{Name: opt.Text, Value: opt.Value})
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
