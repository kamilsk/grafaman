package cmd_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "github.com/kamilsk/grafaman/internal/cmd"
)

func TestNew(t *testing.T) {
	root := New()
	require.NotNil(t, root)
	assert.NotEmpty(t, root.Use)
	assert.NotEmpty(t, root.Short)
	assert.NotEmpty(t, root.Long)
	assert.False(t, root.SilenceErrors)
	assert.True(t, root.SilenceUsage)

	t.Run("defaults", func(t *testing.T) {
		root := New()
		assert.NoError(t, root.PersistentPreRunE(root, nil))
	})

	t.Run("unknown format", func(t *testing.T) {
		root := New()
		require.NoError(t, root.Flag("format").Value.Set("unknown"))
		assert.Error(t, root.PersistentPreRunE(root, nil))
	})

	t.Run("with debug and warning level", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		root := New()
		root.SetErr(buf)
		require.NoError(t, root.Flag("debug").Value.Set("true"))
		require.NoError(t, root.Flag("verbose").Value.Set("1"))
		assert.NoError(t, root.PersistentPreRunE(root, nil))
		assert.Contains(t, buf.String(), "start listen and serve pprof")
	})

	t.Run("with debug and invalid host", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		root := New()
		root.SetErr(buf)
		require.NoError(t, root.Flag("debug").Value.Set("true"))
		require.NoError(t, root.Flag("debug-host").Value.Set("bad:host"))
		assert.Error(t, root.PersistentPreRunE(root, nil))
	})

	t.Run("with debug and info level", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		root := New()
		root.SetErr(buf)
		require.NoError(t, root.Flag("debug").Value.Set("true"))
		require.NoError(t, root.Flag("verbose").Value.Set("2"))
		assert.NoError(t, root.PersistentPreRunE(root, nil))
		assert.Contains(t, buf.String(), "start listen and serve pprof")
	})

	t.Run("with debug and verbose level", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		root := New()
		root.SetErr(buf)
		require.NoError(t, root.Flag("debug").Value.Set("true"))
		require.NoError(t, root.Flag("verbose").Value.Set("5"))
		assert.NoError(t, root.PersistentPreRunE(root, nil))
		assert.Contains(t, buf.String(), "start listen and serve pprof")
	})

	t.Run("with dotenv config", func(t *testing.T) {
		root := New()
		require.NoError(t, root.Flag("env-file").Value.Set("testdata/.env.paas"))
		assert.NoError(t, root.PersistentPreRunE(root, nil))
	})

	t.Run("with app.toml config", func(t *testing.T) {
		root := New()
		require.NoError(t, os.Chdir("testdata"))
		require.NoError(t, root.Flag("env-file").Value.Set(".env"))
		assert.NoError(t, root.PersistentPreRunE(root, nil))
	})
}
