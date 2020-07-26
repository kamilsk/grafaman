package cache_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pelletier/go-toml"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	xtime "go.octolab.org/time"

	"github.com/kamilsk/grafaman/internal/model"
	. "github.com/kamilsk/grafaman/internal/provider/graphite/cache"
)

func TestDecorate(t *testing.T) {
	ctx, prefix := context.Background(), "test"
	metrics := model.Metrics{"metric.a", "metric.b", "metric.c"}

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	t.Run("fetch data from cache", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewMockGraphite(ctrl)

		cache := map[string]interface{}{"metrics": metrics, "ttl": time.Now().Add(time.Hour).Unix()}
		fs := afero.NewMemMapFs()
		file, err := fs.Create(Filename(prefix))
		require.NoError(t, err)
		require.NoError(t, json.NewEncoder(file).Encode(cache))

		decorator := Decorate(provider, fs, logger)
		obtained, err := decorator.Fetch(ctx, prefix, xtime.Week)
		assert.NoError(t, err)
		assert.Equal(t, metrics, obtained)
	})

	t.Run("store data to cache for one day", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewMockGraphite(ctrl)
		provider.EXPECT().
			Fetch(ctx, prefix, xtime.Week).
			Return(metrics, nil)

		fs := afero.NewMemMapFs()

		decorator := Decorate(provider, fs, logger)
		obtained, err := decorator.Fetch(ctx, prefix, xtime.Week)
		assert.NoError(t, err)
		assert.Equal(t, metrics, obtained)
	})

	t.Run("decode invalid data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewMockGraphite(ctrl)

		cache := struct {
			Metrics model.Metrics
			TTL     int64
		}{
			Metrics: metrics,
			TTL:     time.Now().Add(time.Hour).Unix(),
		}
		fs := afero.NewMemMapFs()
		file, err := fs.Create(Filename(prefix))
		require.NoError(t, err)
		require.NoError(t, toml.NewEncoder(file).Encode(cache))

		decorator := Decorate(provider, fs, logger)
		obtained, err := decorator.Fetch(ctx, prefix, xtime.Week)
		assert.Error(t, err)
		assert.Nil(t, obtained)
	})

	t.Run("fail to fetch data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewMockGraphite(ctrl)
		provider.EXPECT().
			Fetch(ctx, prefix, xtime.Week).
			Return(nil, errors.New("service unavailable"))

		decorator := Decorate(provider, afero.NewMemMapFs(), logger)
		obtained, err := decorator.Fetch(ctx, prefix, xtime.Week)
		require.Error(t, err)
		assert.EqualError(t, err, "cache: fetch data: service unavailable")
		assert.Nil(t, obtained)
	})

	t.Run("fail to prepare storage", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewMockGraphite(ctrl)
		fs := NewMockFS(ctrl)
		fs.EXPECT().
			OpenFile(Filename(prefix), os.O_RDWR|os.O_CREATE, os.FileMode(0644)).
			Return(nil, errors.New("fs unhealthy"))

		decorator := Decorate(provider, fs, logger)
		obtained, err := decorator.Fetch(ctx, prefix, xtime.Week)
		require.Error(t, err)
		assert.EqualError(t, err, "cache: prepare storage: fs unhealthy")
		assert.Nil(t, obtained)
	})

	t.Run("fail to truncate data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewMockGraphite(ctrl)
		provider.EXPECT().
			Fetch(ctx, prefix, xtime.Week).
			Return(metrics, nil)
		file := NewMockFile(ctrl)
		file.EXPECT().
			Read(gomock.Any()).
			Return(0, io.EOF)
		file.EXPECT().
			Truncate(int64(0)).
			Return(errors.New("fs unhealthy"))
		file.EXPECT().
			Close().
			Return(errors.New("fs unhealthy"))
		fs := NewMockFS(ctrl)
		fs.EXPECT().
			OpenFile(Filename(prefix), os.O_RDWR|os.O_CREATE, os.FileMode(0644)).
			Return(file, nil)

		decorator := Decorate(provider, fs, logger)
		obtained, err := decorator.Fetch(ctx, prefix, xtime.Week)
		require.Error(t, err)
		assert.EqualError(t, err, "cache: truncate data: fs unhealthy")
		assert.Nil(t, obtained)
	})

	t.Run("fail to prepare to write", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewMockGraphite(ctrl)
		provider.EXPECT().
			Fetch(ctx, prefix, xtime.Week).
			Return(metrics, nil)
		file := NewMockFile(ctrl)
		file.EXPECT().
			Read(gomock.Any()).
			Return(0, io.EOF)
		file.EXPECT().
			Truncate(int64(0)).
			Return(nil)
		file.EXPECT().
			Seek(int64(0), io.SeekStart).
			Return(int64(0), errors.New("fs unhealthy"))
		file.EXPECT().
			Close().
			Return(errors.New("fs unhealthy"))
		fs := NewMockFS(ctrl)
		fs.EXPECT().
			OpenFile(Filename(prefix), os.O_RDWR|os.O_CREATE, os.FileMode(0644)).
			Return(file, nil)

		decorator := Decorate(provider, fs, logger)
		obtained, err := decorator.Fetch(ctx, prefix, xtime.Week)
		require.Error(t, err)
		assert.EqualError(t, err, "cache: prepare to write: fs unhealthy")
		assert.Nil(t, obtained)
	})

	t.Run("fail to store data", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		provider := NewMockGraphite(ctrl)
		provider.EXPECT().
			Fetch(ctx, prefix, xtime.Week).
			Return(metrics, nil)
		file := NewMockFile(ctrl)
		file.EXPECT().
			Read(gomock.Any()).
			Return(0, io.EOF)
		file.EXPECT().
			Truncate(int64(0)).
			Return(nil)
		file.EXPECT().
			Seek(int64(0), io.SeekStart).
			Return(int64(0), nil)
		file.EXPECT().
			Write(gomock.Any()).
			Return(0, errors.New("fs unhealthy"))
		file.EXPECT().
			Close().
			Return(errors.New("fs unhealthy"))
		fs := NewMockFS(ctrl)
		fs.EXPECT().
			OpenFile(Filename(prefix), os.O_RDWR|os.O_CREATE, os.FileMode(0644)).
			Return(file, nil)

		decorator := Decorate(provider, fs, logger)
		obtained, err := decorator.Fetch(ctx, prefix, xtime.Week)
		require.Error(t, err)
		assert.EqualError(t, err, "cache: store data: fs unhealthy")
		assert.Nil(t, obtained)
	})
}

func TestFilename(t *testing.T) {
	filename := Filename("test")
	assert.Equal(t, "test.grafaman.json", filepath.Base(filename))
	assert.Equal(t, ".json", filepath.Ext(filename))
}
