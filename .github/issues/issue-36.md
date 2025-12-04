---
id: 36
database_id: 670796520
node_id: MDU6SXNzdWU2NzA3OTY1MjA=
status: closed
title: "simplify env unpack"
labels: ["help wanted"]
url: https://github.com/kamilsk/grafaman/issues/36
created_at: 2020-08-01T11:45:42Z
updated_at: 2020-08-01T11:55:19Z
---

# simplify env unpack

```
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
```

-> replace an env by `*`
