package cnf_test

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.octolab.org/safe"

	. "github.com/kamilsk/grafaman/internal/cnf"
)

func TestWithConfig(t *testing.T) {
	t.Run("load from dotenv", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
			cnf = Config{File: "testdata/.env.paas"}
		)

		src := cnf // copy
		box.RegisterAlias("app", "app_name")
		box.RegisterAlias("grafana", "grafana_url")
		box.RegisterAlias("dashboard", "grafana_dashboard")
		box.RegisterAlias("graphite", "graphite_url")
		box.RegisterAlias("metrics", "graphite_metrics")

		cmd = Apply(cmd, box, WithConfig(&cnf))
		assert.NoError(t, cmd.PreRunE(cmd, nil))
		assert.NotEqual(t, src, cnf)
		assert.Equal(t, "awesome-service", cnf.App)
		assert.Equal(t, "https://grafana.api/", cnf.Grafana.URL)
		assert.Equal(t, "DTknF4rik", cnf.Grafana.Dashboard)
		assert.Equal(t, "https://graphite.api/", cnf.Graphite.URL)
		assert.Equal(t, "apps.services.awesome-service", cnf.Graphite.Prefix)
	})

	t.Run("load from old dotenv", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
			cnf = Config{File: "testdata/.env"}
		)

		src := cnf // copy
		box.RegisterAlias("app", "app_name")
		box.RegisterAlias("grafana", "grafana_url")
		box.RegisterAlias("dashboard", "grafana_dashboard")
		box.RegisterAlias("graphite", "graphite_url")
		box.RegisterAlias("metrics", "graphite_metrics")

		cmd = Apply(cmd, box, WithConfig(&cnf))
		assert.NoError(t, cmd.PreRunE(cmd, nil))
		assert.NotEqual(t, src, cnf)
		assert.Empty(t, cnf.App)
		assert.Equal(t, "https://grafana.api/", cnf.Grafana.URL)
		assert.Equal(t, "DTknF4rik", cnf.Grafana.Dashboard)
		assert.Equal(t, "https://graphite.api/", cnf.Graphite.URL)
		assert.Equal(t, "apps.services.awesome-service", cnf.Graphite.Prefix)
	})

	t.Run("load from app.toml", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
			cnf = Config{File: ".env.unknown"}
		)

		src := cnf // copy
		box.RegisterAlias("app", "app_name")
		box.RegisterAlias("grafana", "grafana_url")
		box.RegisterAlias("dashboard", "grafana_dashboard")
		box.RegisterAlias("graphite", "graphite_url")
		box.RegisterAlias("metrics", "graphite_metrics")

		cmd = Apply(cmd, box, WithConfig(&cnf))
		require.NoError(t, os.Chdir("testdata"))
		assert.NoError(t, cmd.PreRunE(cmd, nil))
		assert.NotEqual(t, src, cnf)
		assert.Equal(t, "awesome-service", cnf.App)
		assert.Equal(t, "https://grafana.api/", cnf.Grafana.URL)
		assert.Equal(t, "DTknF4rik", cnf.Grafana.Dashboard)
		assert.Equal(t, "https://graphite.api/", cnf.Graphite.URL)
		assert.Equal(t, "apps.services.awesome-service", cnf.Graphite.Prefix)
	})

	t.Run("without config file", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
			cnf = Config{}
		)

		src := cnf // copy
		box.RegisterAlias("app", "app_name")
		box.RegisterAlias("grafana", "grafana_url")
		box.RegisterAlias("dashboard", "grafana_dashboard")
		box.RegisterAlias("graphite", "graphite_url")
		box.RegisterAlias("metrics", "graphite_metrics")

		cmd = Apply(cmd, box, WithConfig(&cnf))
		assert.NoError(t, cmd.PreRunE(cmd, nil))
		assert.Equal(t, src, cnf)
	})
}

func TestWithDebug(t *testing.T) {
	t.Run("flags and bindings", func(t *testing.T) {
		var (
			box = viper.New()
			buf = bytes.NewBuffer(nil)
			cmd = new(cobra.Command)
			cnf = new(Config)
		)

		logger := logrus.New()
		cmd.SetErr(buf)

		cmd = Apply(cmd, box, WithDebug(cnf, logger))
		assert.NoError(t, cmd.ParseFlags([]string{
			"--debug",
			"--debug-host", "127.0.0.1:1234",
			"-vvv",
		}))
		assert.True(t, box.GetBool("debug.enabled"))
		assert.Equal(t, "127.0.0.1:1234", box.GetString("debug.host"))
		assert.Equal(t, 3, box.GetInt("debug.level"))
	})

	t.Run("debug with defaults", func(t *testing.T) {
		var (
			box = viper.New()
			buf = bytes.NewBuffer(nil)
			cmd = new(cobra.Command)
			cnf = new(Config)
		)

		logger := logrus.New()
		cmd.SetErr(buf)

		cmd = Apply(cmd, box, WithDebug(cnf, logger))
		assert.NoError(t, cmd.ParseFlags([]string{"--debug"}))
		assert.NoError(t, box.Unmarshal(cnf))
		assert.NoError(t, cmd.PreRunE(cmd, nil))
		assert.Empty(t, buf.String())
	})

	t.Run("debug with warnings", func(t *testing.T) {
		var (
			box = viper.New()
			buf = bytes.NewBuffer(nil)
			cmd = new(cobra.Command)
			cnf = new(Config)
		)

		logger := logrus.New()
		cmd.SetErr(buf)

		cmd = Apply(cmd, box, WithDebug(cnf, logger))
		assert.NoError(t, cmd.ParseFlags([]string{"--debug", "-v"}))
		assert.NoError(t, box.Unmarshal(cnf))
		assert.NoError(t, cmd.PreRunE(cmd, nil))
		assert.Contains(t, buf.String(), "start listen and serve pprof")
	})

	t.Run("debug with infos", func(t *testing.T) {
		var (
			box = viper.New()
			buf = bytes.NewBuffer(nil)
			cmd = new(cobra.Command)
			cnf = new(Config)
		)

		logger := logrus.New()
		cmd.SetErr(buf)

		cmd = Apply(cmd, box, WithDebug(cnf, logger))
		assert.NoError(t, cmd.ParseFlags([]string{"--debug", "-vv"}))
		assert.NoError(t, box.Unmarshal(cnf))
		assert.NoError(t, cmd.PreRunE(cmd, nil))
		assert.Contains(t, buf.String(), "start listen and serve pprof")
	})

	t.Run("verbose debug", func(t *testing.T) {
		var (
			box = viper.New()
			buf = bytes.NewBuffer(nil)
			cmd = new(cobra.Command)
			cnf = new(Config)
		)

		logger := logrus.New()
		cmd.SetErr(buf)

		cmd = Apply(cmd, box, WithDebug(cnf, logger))
		assert.NoError(t, cmd.ParseFlags([]string{"--debug", "-vvv"}))
		assert.NoError(t, box.Unmarshal(cnf))
		assert.NoError(t, cmd.PreRunE(cmd, nil))
		assert.Contains(t, buf.String(), "start listen and serve pprof")
	})

	t.Run("invalid host", func(t *testing.T) {
		var (
			box = viper.New()
			buf = bytes.NewBuffer(nil)
			cmd = new(cobra.Command)
			cnf = new(Config)
		)

		logger := logrus.New()
		cmd.SetErr(buf)

		cmd = Apply(cmd, box, WithDebug(cnf, logger))
		assert.NoError(t, cmd.ParseFlags([]string{"--debug", "--debug-host", "invalid:host"}))
		assert.NoError(t, box.Unmarshal(cnf))
		assert.Error(t, cmd.PreRunE(cmd, nil))
		assert.Empty(t, buf.String())
	})
}

func TestWithGrafana(t *testing.T) {
	t.Run("configure by flags", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
		)

		cmd = Apply(cmd, box, WithGrafana())
		assert.NoError(t, cmd.ParseFlags([]string{
			"--grafana", "https://grafana.api/",
			"-d", "DTknF4rik",
		}))
		assert.Equal(t, "https://grafana.api/", box.GetString("grafana"))
		assert.Equal(t, "https://grafana.api/", box.GetString("grafana_url"))
		assert.Equal(t, "DTknF4rik", box.GetString("dashboard"))
		assert.Equal(t, "DTknF4rik", box.GetString("grafana_dashboard"))
	})

	t.Run("configure by environment", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
		)

		release, err := setEnvs(
			"GRAFANA_URL", "https://grafana.api/",
			"GRAFANA_DASHBOARD", "DTknF4rik",
		)
		require.NoError(t, err)
		defer safe.Do(release, func(err error) { require.NoError(t, err) })

		cmd = Apply(cmd, box, WithGrafana())
		assert.NoError(t, cmd.ParseFlags(nil))
		assert.Equal(t, "https://grafana.api/", box.GetString("grafana"))
		assert.Equal(t, "https://grafana.api/", box.GetString("grafana_url"))
		assert.Equal(t, "DTknF4rik", box.GetString("dashboard"))
		assert.Equal(t, "DTknF4rik", box.GetString("grafana_dashboard"))
	})
}

func TestWithGraphite(t *testing.T) {
	t.Run("configure by flags", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
		)

		cmd = Apply(cmd, box, WithGraphite())
		assert.NoError(t, cmd.ParseFlags([]string{
			"--filter", "metric.*",
			"--graphite", "https://graphite.api/",
		}))
		assert.Equal(t, "metric.*", box.GetString("filter"))
		assert.Equal(t, "https://graphite.api/", box.GetString("graphite"))
		assert.Equal(t, "https://graphite.api/", box.GetString("graphite_url"))
	})

	t.Run("configure by environment", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
		)

		release, err := setEnvs("GRAPHITE_URL", "https://graphite.api/")
		require.NoError(t, err)
		defer safe.Do(release, func(err error) { require.NoError(t, err) })

		cmd = Apply(cmd, box, WithGraphite())
		assert.NoError(t, cmd.ParseFlags(nil))
		assert.Empty(t, box.GetString("filter"))
		assert.Equal(t, "https://graphite.api/", box.GetString("graphite"))
		assert.Equal(t, "https://graphite.api/", box.GetString("graphite_url"))
	})
}

func TestWithGraphiteMetrics(t *testing.T) {
	t.Run("configure by flags", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
		)

		cmd = Apply(cmd, box, WithGraphiteMetrics())
		assert.NoError(t, cmd.ParseFlags([]string{"-m", "apps.services.awesome-service"}))
		assert.Empty(t, box.GetString("app"))
		assert.Empty(t, box.GetString("app_name"))
		assert.Equal(t, "apps.services.awesome-service", box.GetString("metrics"))
		assert.Equal(t, "apps.services.awesome-service", box.GetString("graphite_metrics"))
	})

	t.Run("configure by environment", func(t *testing.T) {
		var (
			box = viper.New()
			cmd = new(cobra.Command)
		)

		release, err := setEnvs(
			"APP_NAME", "awesome-service",
			"GRAPHITE_METRICS", "apps.services.awesome-service",
		)
		require.NoError(t, err)
		defer safe.Do(release, func(err error) { require.NoError(t, err) })

		cmd = Apply(cmd, box, WithGraphiteMetrics())
		assert.NoError(t, cmd.ParseFlags(nil))
		assert.Equal(t, "awesome-service", box.GetString("app"))
		assert.Equal(t, "awesome-service", box.GetString("app_name"))
		assert.Equal(t, "apps.services.awesome-service", box.GetString("metrics"))
		assert.Equal(t, "apps.services.awesome-service", box.GetString("graphite_metrics"))
	})
}

func TestWithOutputFormat(t *testing.T) {
	var (
		box = viper.New()
		cmd = new(cobra.Command)
	)

	cmd = Apply(cmd, box, WithOutputFormat())
	assert.NoError(t, cmd.ParseFlags([]string{"-f", "json"}))
	assert.Equal(t, "json", box.GetString("output.format"))
}

// helpers

var mx sync.Mutex

// issue: https://github.com/octolab/pkg/issues/22
func setEnvs(envs ...string) (func() error, error) {
	mx.Lock()

	before := make([]*string, len(envs))
	for i := 0; i < len(envs); i += 2 {
		var prev *string
		val, present := os.LookupEnv(envs[i])
		if present {
			prev = &val
		}
		before[i], before[i+1] = &envs[i], prev
		if err := os.Setenv(envs[i], envs[i+1]); err != nil {
			mx.Unlock()
			return nil, fmt.Errorf("cannot set environment variable %s=%q", envs[i], envs[i+1])
		}
	}

	return func() error {
		defer mx.Unlock()
		for i := 0; i < len(before); i += 2 {
			var err error
			env, val := before[i], before[i+1]
			if val == nil {
				err = os.Unsetenv(*env)
			} else {
				err = os.Setenv(*env, *val)
			}
			if err != nil {
				return fmt.Errorf("cannot restore previos environment variable %s", *env)
			}
		}
		return nil
	}, nil
}
