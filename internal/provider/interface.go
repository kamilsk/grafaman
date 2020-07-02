package provider

import (
	"context"
	"time"

	"github.com/kamilsk/grafaman/internal/model"
)

type Graphite interface {
	Fetch(context.Context, string, time.Duration) (model.Metrics, error)
}
