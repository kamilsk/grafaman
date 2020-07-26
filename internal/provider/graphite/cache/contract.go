package cache

import (
	"context"
	"time"

	"github.com/spf13/afero"

	"github.com/kamilsk/grafaman/internal/model"
)

//go:generate mockgen -source $GOFILE -destination mocks_test.go -package ${GOPACKAGE}_test

// Graphite defines Graphite provider interface.
type Graphite interface {
	Fetch(context.Context, string, time.Duration) (model.Metrics, error)
}

// Proxies for mocking.
type (
	File interface{ afero.File }
	FS   interface{ afero.Fs }
)
