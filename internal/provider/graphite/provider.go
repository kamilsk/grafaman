package graphite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"

	"github.com/pkg/errors"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"
	"golang.org/x/sync/errgroup"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

// New returns an instance of Graphite metrics provider.
func New(endpoint string) (*provider, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "prepare Graphite metrics provider endpoint URL")
	}
	return &provider{
		client:   &http.Client{},
		endpoint: *u,
	}, nil
}

type provider struct {
	client   *http.Client
	endpoint url.URL
}

// Fetch walks through the endpoint and takes all metrics with the specified prefix.
// Documentation: https://graphite-api.readthedocs.io/en/latest/api.html#metrics-find.
func (provider *provider) Fetch(ctx context.Context, prefix string, last time.Duration) (entity.Metrics, error) {
	const source = "/metrics/find"

	u := provider.endpoint
	u.Path = path.Join(u.Path, source)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "create Graphite metrics base request")
	}
	q := req.URL.Query()
	q.Add(formatParam, "json")
	q.Add(fromParam, fmt.Sprintf("now-%s", last))
	q.Add(untilParam, "now")
	q.Add(queryParam, prefix)
	req.URL.RawQuery = q.Encode()

	var (
		aggregator = make(chan dto, runtime.GOMAXPROCS(0))
		metrics    = make(entity.Metrics, 0, 8)
	)

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		for node := range aggregator {
			metrics = append(metrics, entity.Metric(node.ID))
		}
		return nil
	})
	g.Go(func() error {
		defer close(aggregator)
		return provider.fetch(ctx, aggregator, req)
	})

	return metrics, g.Wait()
}

func (provider *provider) fetch(ctx context.Context, out chan<- dto, req *http.Request) error {
	resp, err := provider.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "Graphite metrics recursive fetch request")
	}
	defer safe.Close(resp.Body, unsafe.Ignore)

	var nodes []dto
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return errors.Wrap(err, "decode Graphite metrics fetch response")
	}

	g, ctx := errgroup.WithContext(ctx)
	for _, node := range nodes {
		if node.Leaf == 1 {
			out <- node
			continue
		}
		query := node.ID + ".*"
		g.Go(func() error {
			req := req.Clone(ctx)
			q := req.URL.Query()
			q.Set(queryParam, query)
			req.URL.RawQuery = q.Encode()
			return provider.fetch(ctx, out, req)
		})
	}
	return g.Wait()
}
