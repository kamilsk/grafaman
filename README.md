> # 📈 grafaman
>
> Metrics coverage reporter for Graphite and Grafana.

[![Build][build.icon]][build.page]
[![Template][template.icon]][template.page]

## 💡 Idea

...

Full description of the idea is available [here][design.page].

## 🏆 Motivation

...

## 🤼‍♂️ How to

...

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

### Source

```bash
# use standard go tools
$ go get github.com/kamilsk/grafaman@latest
# or use egg tool
$ egg tools add github.com/kamilsk/grafaman@latest
```

> [egg][]<sup id="anchor-egg">[1](#egg)</sup> is an `extended go get`.

### Bash and Zsh completions

```bash
$ grafaman completion bash > /path/to/bash_completion.d/grafaman.sh
$ grafaman completion zsh  > /path/to/zsh-completions/_grafaman.zsh
```

---

made with ❤️ for everyone

[build.page]:       https://travis-ci.com/kamilsk/grafaman
[build.icon]:       https://travis-ci.com/kamilsk/grafaman.svg?branch=master
[design.page]:      https://www.notion.so/33715348cc114ea79dd350a25d16e0b0?r=0b753cbf767346f5a6fd51194829a2f3
[promo.page]:       https://github.com/kamilsk/grafaman
[template.page]:    https://github.com/octomation/go-tool
[template.icon]:    https://img.shields.io/badge/template-go--tool-blue

[egg]:              https://github.com/kamilsk/egg
