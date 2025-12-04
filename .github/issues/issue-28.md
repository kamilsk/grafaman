---
id: 28
database_id: 644932787
node_id: MDU6SXNzdWU2NDQ5MzI3ODc=
status: closed
title: "refactor coverage reporter"
labels: ["help wanted"]
url: https://github.com/kamilsk/grafaman/issues/28
created_at: 2020-06-24T20:31:29Z
updated_at: 2020-06-28T21:30:42Z
---

# refactor coverage reporter

- [x] hold queries
```go
coverage.New().Report(metrics, queries)
->
coverage.New(queries).Report(metrics)
```
- [x] simplify function signature
```go
repl.NewCoverageExecutor(config.Graphite.Prefix, metrics, queries, printer, logger)
->
repl.NewCoverageExecutor(reporter, metrics, printer, logger)
```
- [x] simplify filter
```go
Filter(metrics provider.Metrics, pattern, prefix string)
->
Filter(metrics provider.Metrics, pattern string)
```
- [x] ~find way to filter/exclude in place~
