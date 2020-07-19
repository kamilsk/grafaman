module github.com/kamilsk/grafaman

go 1.13

require (
	github.com/alexeyco/simpletable v0.0.0-20200203113705-55bd62a5b8df
	github.com/c-bata/go-prompt v0.2.3
	github.com/go-graphite/carbonapi v0.0.0-20200617193347-7bbdac316538
	github.com/gobwas/glob v0.2.3
	github.com/golang/mock v1.4.3
	github.com/kamilsk/retry/v5 v5.0.0-rc5
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mattn/go-tty v0.0.3 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/term v0.0.0-20200520122047-c3ffed290a03 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/afero v1.3.2
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	go.octolab.org v0.2.0
	go.octolab.org/toolkit/cli v0.1.1
	go.octolab.org/toolkit/config v0.0.3
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
)

// hg -> git
replace bitbucket.org/tebeka/strftime => github.com/tebeka/strftime v0.1.4
