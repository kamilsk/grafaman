package graphite

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"strconv"
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

const (
	fromKey  = "from"
	queryKey = "query"
)

type dto struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Leaf int    `json:"leaf"`
}

type provider struct {
	client   *http.Client
	endpoint url.URL
}

func (provider *provider) Fetch(ctx context.Context, subset string, fast bool) (entity.Metrics, error) {
	const method = "/metrics/find"

	u := provider.endpoint
	u.Path = path.Join(u.Path, method)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "create Graphite metrics base request")
	}
	q, week := req.URL.Query(), 7*24*time.Hour
	q.Add(fromKey, strconv.Itoa(int(time.Now().Add(-week).Unix())))
	q.Add(queryKey, subset)
	req.URL.RawQuery = q.Encode()

	// try to fetch fast by one query
	if fast { // TODO:research fastFetch returns invalid state
		metrics, err := provider.fastFetch(ctx, req)
		if err == nil && metrics.Len() > 0 {
			return metrics, nil
		}
	}

	// fallback to recursive algorithm
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

func (provider *provider) fastFetch(ctx context.Context, req *http.Request) (entity.Metrics, error) {
	req = req.Clone(ctx)
	q := req.URL.Query()
	q.Set(queryKey, q.Get(queryKey)+".~")
	req.URL.RawQuery = q.Encode()

	resp, err := provider.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Graphite metrics fast fetch request")
	}
	defer safe.Close(resp.Body, unsafe.Ignore)

	var nodes []dto
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
		return nil, errors.Wrap(err, "decode Graphite metrics fast fetch response")
	}

	metrics := make(entity.Metrics, 0, len(nodes)/2)
	for _, node := range nodes {
		if node.Leaf == 1 {
			metrics = append(metrics, entity.Metric(node.ID))
		}
	}
	return metrics, nil
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
			q.Set(queryKey, query)
			req.URL.RawQuery = q.Encode()
			return provider.fetch(ctx, out, req)
		})
	}
	return g.Wait()
}
