---
id: 6
database_id: 622005332
node_id: MDU6SXNzdWU2MjIwMDUzMzI=
status: closed
title: "implement REPL"
labels: ["help wanted"]
url: https://github.com/kamilsk/grafaman/issues/6
created_at: 2020-05-20T18:44:01Z
updated_at: 2020-06-28T21:26:45Z
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
