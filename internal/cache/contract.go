package cache

import (
	"context"
	"time"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

type Graphite interface {
	Fetch(context.Context, string, time.Duration) (entity.Metrics, error)
}
