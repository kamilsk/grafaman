package provider

import (
	"context"
	"time"
)

type Graphite interface {
	Fetch(context.Context, string, time.Duration) (Metrics, error)
}
