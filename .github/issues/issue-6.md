---
code: null
id: MDU6SXNzdWU2MjIwMDUzMzI=
databaseId: 622005332
number: 6
url: https://github.com/kamilsk/grafaman/issues/6
title: "implement REPL"
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
createdAt: 2020-05-20T18:44:01Z
updatedAt: 2020-06-28T21:26:45Z
lastEditedAt: null
closedAt: 2020-06-28T21:26:45Z
---

# implement REPL

```bash
$ grafaman coverage --repl

apps.services.*

# +--------------------------------------------------------------------+--------+
# | Metric                                                             | Hits   |
# +--------------------------------------------------------------------+--------+
# | apps.services.awesome-service.jaeger.finished_spans_sampled_n      |      0 |
# | apps.services.awesome-service.rpc.client.success.ok.percentile.75  |      1 |
# | apps.services.awesome-service.rpc.client.success.ok.percentile.95  |      1 |
# | apps.services.awesome-service.rpc.client.success.ok.percentile.99  |      2 |
# | apps.services.awesome-service.rpc.client.success.ok.percentile.999 |      1 |
# | ...                                                                |    ... |
# | apps.services.awesome-service.go.pod-5dbdcd5dbb-6z58f.threads      |      0 |
# +--------------------------------------------------------------------+--------+
# |                                                              Total | 65.77% |
# +--------------------------------------------------------------------+--------+

...

apps.services.rpc.client.*.*.95

# +--------------------------------------------------------------------+--------+
# | Metric                                                             | Hits   |
# +--------------------------------------------------------------------+--------+
# | apps.services.awesome-service.rpc.client.success.ok.percentile.95  |      1 |
# | ...                                                                |    ... |
# +--------------------------------------------------------------------+--------+
# |                                                              Total | 99.99% |
# +--------------------------------------------------------------------+--------+

```

Requirements:

- refresh table body
- recalculate total coverage

<!-- 2020-05-23T06:46Z https://github.com/kamilsk/grafaman/issues/6#issuecomment-632996548
add a possibility to add the output pattern to a list of "metric clusters", e.g.

the output list
- apps.services.rpc.client.*.error
- apps.services.internal.*.error
-->

<!-- 2020-06-28T21:26Z https://github.com/kamilsk/grafaman/issues/6#issuecomment-650824066
implemented
-->
