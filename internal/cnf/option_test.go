package cnf_test

import (
	"errors"
	"os"
	"testing"

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
