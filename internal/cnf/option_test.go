package cnf_test

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"

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
