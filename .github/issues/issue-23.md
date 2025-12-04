---
id: 23
database_id: 643289372
node_id: MDU6SXNzdWU2NDMyODkzNzI=
status: closed
title: "add cache lookup command"
labels: ["help wanted"]
url: https://github.com/kamilsk/grafaman/issues/23
created_at: 2020-06-22T18:50:31Z
updated_at: 2020-06-23T19:35:44Z
---

# add cache lookup command

```bash
$ grafaman cache lookup --set some.metrics
metrics    /var/tmp/some.metrics.json

$ grafaman cache lookup metrics --set some.metrics
/var/tmp/some.metrics.json
```
