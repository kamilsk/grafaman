---
code: null
id: MDU6SXNzdWU2MTgwNjg1Mzc=
databaseId: 618068537
number: 2
url: https://github.com/kamilsk/grafaman/issues/2
title: "replace source of tebeka/strftime"
issueType: null
labels:
  - "help wanted"
assignees:
  - kamilsk
milestone: null
parent: null
issueFields: {}
state: CLOSED
stateReason: COMPLETED
duplicateOf: null
createdAt: 2020-05-14T09:19:00Z
updatedAt: 2020-05-14T09:24:55Z
lastEditedAt: null
closedAt: 2020-05-14T09:24:55Z
---

# replace source of tebeka/strftime

because hg is an unnecessary dependency

<!-- 2020-05-14T09:19Z https://github.com/kamilsk/grafaman/issues/2#issuecomment-628509980
- https://github.com/tebeka/strftime/commits/master
- https://bitbucket.org/tebeka/strftime/commits/
-->

<!-- 2020-05-14T09:20Z https://github.com/kamilsk/grafaman/issues/2#issuecomment-628510231
possible patch is

```vgo
replace bitbucket.org/tebeka/strftime => github.com/tebeka/strftime v0.1.4
```
-->
