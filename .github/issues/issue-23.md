---
code: null
id: MDU6SXNzdWU2NDMyODkzNzI=
databaseId: 643289372
number: 23
url: https://github.com/kamilsk/grafaman/issues/23
title: "add cache lookup command"
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
createdAt: 2020-06-22T18:50:31Z
updatedAt: 2020-06-23T19:35:44Z
lastEditedAt: 2020-06-22T18:53:59Z
closedAt: 2020-06-23T19:35:44Z
---

# add cache lookup command

```bash
$ grafaman cache lookup --set some.metrics
metrics    /var/tmp/some.metrics.json

$ grafaman cache lookup metrics --set some.metrics
/var/tmp/some.metrics.json
```
