---
id: 37
database_id: 670845955
node_id: MDU6SXNzdWU2NzA4NDU5NTU=
status: closed
title: "coverage didn't work properly"
labels: ["help wanted"]
url: https://github.com/kamilsk/grafaman/issues/37
created_at: 2020-08-01T12:52:26Z
updated_at: 2020-08-01T21:17:17Z
---

# coverage didn't work properly

## the report

```
+-------------------------------------------------------------+-------+
| Metric of apps.services.awesome-service                     | Hits  |
+-------------------------------------------------------------+-------+
| service.api.source.get.POST.request_time.200.count          |     2 |
| service.api.source.get.POST.request_time.200.max            |     0 |
| service.api.source.get.POST.request_time.200.mean           |     0 |
| service.api.source.get.POST.request_time.200.median         |     0 |
| service.api.source.get.POST.request_time.200.min            |     0 |
| service.api.source.get.POST.request_time.200.percentile.75  |     0 |
| service.api.source.get.POST.request_time.200.percentile.95  |     0 |
| service.api.source.get.POST.request_time.200.percentile.98  |     0 |
| service.api.source.get.POST.request_time.200.percentile.99  |     0 |
| service.api.source.get.POST.request_time.200.percentile.999 |     0 |
| service.api.source.get.POST.request_time.200.sum            |     0 |
| service.api.source.get.POST.request_time.499.count          |     2 |
| service.api.source.get.POST.request_time.499.max            |     0 |
| service.api.source.get.POST.request_time.499.mean           |     0 |
| service.api.source.get.POST.request_time.499.median         |     0 |
| service.api.source.get.POST.request_time.499.min            |     0 |
| service.api.source.get.POST.request_time.499.percentile.75  |     0 |
| service.api.source.get.POST.request_time.499.percentile.95  |     0 |
| service.api.source.get.POST.request_time.499.percentile.98  |     0 |
| service.api.source.get.POST.request_time.499.percentile.99  |     0 |
| service.api.source.get.POST.request_time.499.percentile.999 |     0 |
| service.api.source.get.POST.request_time.499.sum            |     0 |
| service.api.source.get.POST.request_time.500.count          |     2 |
| service.api.source.get.POST.request_time.500.max            |     0 |
| service.api.source.get.POST.request_time.500.mean           |     0 |
| service.api.source.get.POST.request_time.500.median         |     0 |
| service.api.source.get.POST.request_time.500.min            |     0 |
| service.api.source.get.POST.request_time.500.percentile.75  |     0 |
| service.api.source.get.POST.request_time.500.percentile.95  |     0 |
| service.api.source.get.POST.request_time.500.percentile.98  |     0 |
| service.api.source.get.POST.request_time.500.percentile.99  |     0 |
| service.api.source.get.POST.request_time.500.percentile.999 |     0 |
| service.api.source.get.POST.request_time.500.sum            |     0 |
| service.api.source.get.POST.request_time.503.count          |     2 |
| service.api.source.get.POST.request_time.503.max            |     0 |
| service.api.source.get.POST.request_time.503.mean           |     0 |
| service.api.source.get.POST.request_time.503.median         |     0 |
| service.api.source.get.POST.request_time.503.min            |     0 |
| service.api.source.get.POST.request_time.503.percentile.75  |     0 |
| service.api.source.get.POST.request_time.503.percentile.95  |     0 |
| service.api.source.get.POST.request_time.503.percentile.98  |     0 |
| service.api.source.get.POST.request_time.503.percentile.99  |     0 |
| service.api.source.get.POST.request_time.503.percentile.999 |     0 |
| service.api.source.get.POST.request_time.503.sum            |     0 |
+-------------------------------------------------------------+-------+
|                                                       Total | 9.09% |
+-------------------------------------------------------------+-------+
```

## the queries

```
+----------------------------------------------------------------------------------------------------------+
| Query                                                                                                    |
+----------------------------------------------------------------------------------------------------------+
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.$code.count                  |
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.$code.percentile.$percentile |
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.*.count                      |
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.*.count                      |
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.*.count                      |
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.200.count                    |
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.4*.count                     |
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.4*.percentile.$percentile    |
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.5*.count                     |
| apps.services.awesome-service.service.api.$source.$method.POST.request_time.5*.percentile.$percentile    |
+----------------------------------------------------------------------------------------------------------+
| Total: xx                                                                                                |
+----------------------------------------------------------------------------------------------------------+
```

## expected

percentile full covered
