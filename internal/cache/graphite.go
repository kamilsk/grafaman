package cache

import (
	"context"
	"time"

	"github.com/spf13/afero"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

func Wrap(provider Graphite, fs afero.Fs) Graphite {
	return &graphite{provider, fs}
}

type graphite struct {
	provider Graphite
	fs       afero.Fs
}

func (decorator *graphite) Fetch(ctx context.Context, prefix string, last time.Duration) (entity.Metrics, error) {
	return decorator.provider.Fetch(ctx, prefix, last)
}
