package cnf_test

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/kamilsk/grafaman/internal/cnf"
)

func TestExtendCommand(t *testing.T) {
	t.Run("after run", func(t *testing.T) {
		var command cobra.Command

		After(&command.Run, func(cmd *cobra.Command, args []string) { cmd.Use = "first" })
		assert.Empty(t, command.Use)
		command.Run(&command, nil)
		assert.Equal(t, command.Use, "first")

		After(&command.Run, func(cmd *cobra.Command, args []string) { cmd.Use = "last" })
		assert.Equal(t, command.Use, "first")
		command.Run(&command, nil)
		assert.Equal(t, command.Use, "last")
	})

	t.Run("after run with error", func(t *testing.T) {
		var command cobra.Command

		AfterE(&command.RunE, func(cmd *cobra.Command, args []string) error {
			cmd.Use = "first"
			return nil
		})
		assert.Empty(t, command.Use)
		assert.NoError(t, command.RunE(&command, nil))
		assert.Equal(t, command.Use, "first")

		AfterE(&command.RunE, func(cmd *cobra.Command, args []string) error {
			cmd.Use = "last"
			return errors.New("test")
		})
		assert.Equal(t, command.Use, "first")
		assert.Error(t, command.RunE(&command, nil))
		assert.Equal(t, command.Use, "last")

		AfterE(&command.RunE, func(cmd *cobra.Command, args []string) error {
			cmd.Use = "first"
			return errors.New("test")
		})
		assert.Equal(t, command.Use, "last")
		assert.Error(t, command.RunE(&command, nil))
		assert.Equal(t, command.Use, "last")
	})

	t.Run("before run", func(t *testing.T) {
		var command cobra.Command

		Before(&command.Run, func(cmd *cobra.Command, args []string) { cmd.Use = "first" })
		assert.Empty(t, command.Use)
		command.Run(&command, nil)
		assert.Equal(t, command.Use, "first")

		Before(&command.Run, func(cmd *cobra.Command, args []string) { cmd.Use = "last" })
		assert.Equal(t, command.Use, "first")
		command.Run(&command, nil)
		assert.Equal(t, command.Use, "first")
	})

	t.Run("before run with error", func(t *testing.T) {
		var command cobra.Command

		BeforeE(&command.RunE, func(cmd *cobra.Command, args []string) error {
			cmd.Use = "first"
			return nil
		})
		assert.Empty(t, command.Use)
		assert.NoError(t, command.RunE(&command, nil))
		assert.Equal(t, command.Use, "first")

		BeforeE(&command.RunE, func(cmd *cobra.Command, args []string) error {
			cmd.Use = "last"
			return errors.New("test")
		})
		assert.Equal(t, command.Use, "first")
		assert.Error(t, command.RunE(&command, nil))
		assert.Equal(t, command.Use, "last")

		BeforeE(&command.RunE, func(cmd *cobra.Command, args []string) error {
			cmd.Use = "first"
			return errors.New("test")
		})
		assert.Equal(t, command.Use, "last")
		assert.Error(t, command.RunE(&command, nil))
		assert.Equal(t, command.Use, "first")
	})
}

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
		assert.Equal(t, "DTknF4rik", cnf.Grafana.Dashboard)
		assert.Equal(t, "https://grafana.api/", cnf.Grafana.URL)
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
		assert.Equal(t, "", cnf.App)
		assert.Equal(t, "DTknF4rik", cnf.Grafana.Dashboard)
		assert.Equal(t, "https://grafana.api/", cnf.Grafana.URL)
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
		assert.Equal(t, "DTknF4rik", cnf.Grafana.Dashboard)
		assert.Equal(t, "https://grafana.api/", cnf.Grafana.URL)
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
