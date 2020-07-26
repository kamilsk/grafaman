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
)

// Decorate wraps Graphite provider by cache layer.
func Decorate(provider Graphite, fs afero.Fs, logger *logrus.Logger) Graphite {
	return &decorator{provider, fs, logger}
}

// Filename returns cache file name.
func Filename(prefix string) string {
	return filepath.Join(os.TempDir(), prefix) + ".grafaman.json"
}

type decorator struct {
	provider Graphite
	fs       afero.Fs
	logger   *logrus.Logger
}

// Fetch tries to load data from cache first or fallback
// to a decorated provider and store its success response.
func (decorator *decorator) Fetch(ctx context.Context, prefix string, last time.Duration) (model.Metrics, error) {
	filename := Filename(prefix)
	logger := decorator.logger.WithFields(logrus.Fields{"component": "cache", "file": filename})
	file, err := decorator.fs.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.WithError(err).Error("prepare storage")
		return nil, errors.Wrap(err, "cache: prepare storage")
	}
	defer safe.Close(file, func(err error) { logger.WithError(err).Warning("flush data") })

	var data struct {
		Metrics model.Metrics `json:"metrics"`
		TTL     int64         `json:"ttl"`
	}
	if err := json.NewDecoder(file).Decode(&data); err != nil && !errors.Is(err, io.EOF) {
		logger.WithError(err).Error("decode data")
		return nil, errors.Wrap(err, "cache: decode data")
	}

	now := time.Now()
	if time.Unix(data.TTL, 0).After(now) {
		logger.Info("fetch data from cache")
		return data.Metrics, nil
	}

	data.Metrics, err = decorator.provider.Fetch(ctx, prefix, last)
	if err != nil {
		logger.WithError(err).Error("fetch data")
		return nil, errors.Wrap(err, "cache: fetch data")
	}
	data.TTL = now.Add(xtime.Day).Unix()

	if err := file.Truncate(0); err != nil {
		logger.WithError(err).Error("truncate data")
		return nil, errors.Wrap(err, "cache: truncate data")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		logger.WithError(err).Error("prepare to write")
		return nil, errors.Wrap(err, "cache: prepare to write")
	}
	if err := json.NewEncoder(file).Encode(data); err != nil {
		logger.WithError(err).Error("store data")
		return nil, errors.Wrap(err, "cache: store data")
	}

	logger.Info("store data to cache for one day")
	return data.Metrics, nil
}
