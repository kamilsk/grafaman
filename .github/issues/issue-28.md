---
code: null
id: MDU6SXNzdWU2NDQ5MzI3ODc=
databaseId: 644932787
number: 28
url: https://github.com/kamilsk/grafaman/issues/28
title: "refactor coverage reporter"
issueType: null
labels:
  - "help wanted"
assignees:
  - kamilsk
milestone: "[[milestone-1]]"
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-06-24T20:31:29Z
updatedAt: 2020-06-28T21:30:42Z
lastEditedAt: 2020-06-28T21:30:34Z
closedAt: 2020-06-28T21:30:42Z
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

<!-- 2020-06-28T21:30Z https://github.com/kamilsk/grafaman/issues/28#issuecomment-650824489
implemented
-->
