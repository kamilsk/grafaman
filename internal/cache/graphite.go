package cache

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"go.octolab.org/safe"
	xtime "go.octolab.org/time"

	"github.com/kamilsk/grafaman/internal/model"
	. "github.com/kamilsk/grafaman/internal/provider"
)

func WrapGraphiteProvider(provider Graphite, fs afero.Fs, logger *logrus.Logger) *graphite {
	return &graphite{provider, fs, logger}
}

type graphite struct {
	provider Graphite
	fs       afero.Fs
	logger   *logrus.Logger
}

func (decorator *graphite) Fetch(ctx context.Context, prefix string, last time.Duration) (model.Metrics, error) {
	key := decorator.Key(prefix)
	logger := decorator.logger.WithField("file", key)
	file, err := decorator.fs.OpenFile(key, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.WithError(err).Error("read data from cache")
		return nil, errors.Wrap(err, "cache: read data")
	}
	defer safe.Close(file, func(err error) { logger.WithError(err).Warning("close cache file") })

	var data struct {
		Metrics model.Metrics `json:"metrics"`
		TTL     int64         `json:"ttl"`
	}
	if err := json.NewDecoder(file).Decode(&data); err != nil && !errors.Is(err, io.EOF) {
		logger.WithError(err).Error("decode data from cache")
		return nil, errors.Wrap(err, "cache: decode data")
	}

	now := time.Now()
	if time.Unix(data.TTL, 0).After(now) {
		logger.Info("return data from cache")
		return data.Metrics, nil
	}

	data.Metrics, err = decorator.provider.Fetch(ctx, prefix, last)
	if err != nil {
		logger.WithError(err).Error("cannot fetch data by provider")
		return nil, errors.Wrap(err, "cache: fetch data")
	}
	data.TTL = now.Add(xtime.Day).Unix()
	{
		if err := file.Truncate(0); err != nil {
			logger.WithError(err).Error("reset cache")
			return nil, errors.Wrap(err, "cache: reset")
		}
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			logger.WithError(err).Error("prepare cache to write")
			return nil, errors.Wrap(err, "cache: prepare to write")
		}
	}
	if err := json.NewEncoder(file).Encode(data); err != nil {
		logger.WithError(err).Error("write data to cache")
		return nil, errors.Wrap(err, "cache: write data")
	}
	logger.Info("store cache for one day")
	return data.Metrics, nil
}

func (decorator *graphite) Key(prefix string) string {
	return filepath.Join(os.TempDir(), prefix) + ".grafaman.json"
}
