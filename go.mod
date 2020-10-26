module github.com/kamilsk/grafaman

go 1.15

require (
	github.com/alexeyco/simpletable v0.0.0-20200730140406-5bb24159ccfb
	github.com/c-bata/go-prompt v0.2.5
	github.com/go-graphite/carbonapi v0.0.0-20200617193347-7bbdac316538
	github.com/gobwas/glob v0.2.3
	github.com/golang/mock v1.4.4
	github.com/kamilsk/retry/v5 v5.0.0-rc5
	github.com/mitchellh/mapstructure v1.3.3
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	github.com/pelletier/go-toml v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/afero v1.4.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	go.octolab.org v0.5.4
	go.octolab.org/toolkit/cli v0.2.0
	go.octolab.org/toolkit/config v0.0.4
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9
)

// hg -> git
replace bitbucket.org/tebeka/strftime => github.com/tebeka/strftime v0.1.5
