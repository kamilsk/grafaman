> # 📈 grafaman
>
> Metrics coverage reporter for [Graphite][] and [Grafana][].

[![Build][build.icon]][build.page]
[![Template][template.icon]][template.page]

## 💡 Idea

```bash
$ grafaman coverage \
    --grafana https://grafana.api/ -d DTknF4rik \
    --graphite https://graphite.api/ \
    --subset apps.services.awesome-service
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
```

Full description of the idea is available [here][design.page].

## 🏆 Motivation

At [Avito](https://tech.avito.ru/), we develop many services built on top of our excellent
[PaaS](https://en.wikipedia.org/wiki/Platform_as_a_service) and internal modules. These services send
a lot of metrics about their internal state that are then output to [Grafana][] dashboards.

I need a tool that helps me to understand what metrics are published by services
and how many of them are presented at [Grafana][] dashboards.

## 🤼‍♂️ How to

### Metrics coverage report

```bash
$ grafaman coverage \
    --grafana https://grafana.api/ -d DTknF4rik \
    --graphite https://graphite.api/ \
    -s apps.services.awesome-service \
    --trim='complex.$env.' --trim='env.$env.' \
    --exclude='*.count' --exclude='*.max' --exclude='*.min' --exclude='*.sum'
```

Supported environment variables:

- GRAFANA_URL
- GRAFANA_DASHBOARD
- GRAPHITE_URL
- GRAPHITE_METRICS

### Fetch metrics from [Graphite][]

```bash
$ grafaman metrics -e https://graphite.api/ -s apps.services.awesome-service
```

### Fetch queries from [Grafana][]

```bash
$ grafaman queries -e https://grafana.api/ -d DTknF4rik \
    -s apps.services.awesome-service \
    --trim='complex.$env.' --trim='env.$env.' \
    --sort
```

## 🧩 Installation

### Homebrew

```bash
$ brew install kamilsk/tap/grafaman
```

### Binary

```bash
$ curl -sSfL https://raw.githubusercontent.com/kamilsk/grafaman/master/bin/install | sh
# or
$ wget -qO-  https://raw.githubusercontent.com/kamilsk/grafaman/master/bin/install | sh
```

> don't forget about [security](https://www.idontplaydarts.com/2016/04/detecting-curl-pipe-bash-server-side/).

### Source

```bash
# use standard go tools
$ go get github.com/kamilsk/grafaman@latest
# or use egg tool
$ egg tools add github.com/kamilsk/grafaman@latest
```

> [egg][] is an `extended go get`.

### Bash and Zsh completions

```bash
$ grafaman completion bash > /path/to/bash_completion.d/grafaman.sh
$ grafaman completion zsh  > /path/to/zsh-completions/_grafaman.zsh
```

---

made with ❤️ for everyone

[build.page]:       https://travis-ci.com/kamilsk/grafaman
[build.icon]:       https://travis-ci.com/kamilsk/grafaman.svg?branch=master
[design.page]:      https://www.notion.so/octolab/grafaman-06e6fcd46c924126ae134c69dafbca6c?r=0b753cbf767346f5a6fd51194829a2f3
[promo.page]:       https://github.com/kamilsk/grafaman
[template.page]:    https://github.com/octomation/go-tool
[template.icon]:    https://img.shields.io/badge/template-go--tool-blue

[egg]:              https://github.com/kamilsk/egg
[Graphite]:         https://graphiteapp.org/
[Grafana]:          https://grafana.com/
