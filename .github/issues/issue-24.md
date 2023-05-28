---
code: null
id: MDU6SXNzdWU2NDMyOTA3NDM=
databaseId: 643290743
number: 24
url: https://github.com/kamilsk/grafaman/issues/24
title: "add cache drop command"
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
createdAt: 2020-06-22T18:53:00Z
updatedAt: 2020-06-22T20:53:50Z
lastEditedAt: null
closedAt: 2020-06-22T20:53:50Z
---

# add cache drop command

```bash
$ grafaman cache drop
# drops all cache

$ grafaman cache drop metrics ...
# drops cache for metrics ...
```

<!-- 2020-06-22T20:53Z https://github.com/kamilsk/grafaman/issues/24#issuecomment-647764183
won't do. alternatives:

```bash
$ rm -f $(grafaman cache lookup)

$ grafaman cache lookup | xargs rm -f
```
-->
