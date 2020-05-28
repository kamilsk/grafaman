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
	"go.octolab.org/sequence"
	"go.octolab.org/unsafe"
	"golang.org/x/sync/errgroup"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

// New returns an instance of Graphite metrics provider.
func New(endpoint string) (*provider, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "graphite: prepare metrics provider endpoint URL")
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
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "graphite: create metrics base request")
	}
	q := request.URL.Query()
	q.Add(formatParam, "json")
	q.Add(fromParam, fmt.Sprintf("now-%s", last))
	q.Add(untilParam, "now")
	q.Add(queryParam, prefix)
	request.URL.RawQuery = q.Encode()

	var (
		aggregator = make(chan dto, 256)
		result     = make(chan []dto, 256)
		metrics    = make(entity.Metrics, 0, 256)
		requests   = make(chan *http.Request, 256)
	)

	main, ctx := errgroup.WithContext(ctx)
	main.Go(func() error {
		for {
			select {
			case node, ok := <-aggregator:
				if !ok {
					return nil
				}
				metrics = append(metrics, entity.Metric(node.ID))
			case <-ctx.Done():
				return errors.Wrap(ctx.Err(), "graphite: aggregator process")
			}
		}
	})
	main.Go(func() error {
		defer close(aggregator)
		defer close(requests)
		var counter int

		requests <- request
		counter++

		for counter > 0 {
			select {
			case nodes, ok := <-result:
				counter--
				if !ok {
					return nil
				}
				for _, node := range nodes {
					if node.Leaf == 1 {
						select {
						case aggregator <- node:
						case <-ctx.Done():
							return errors.Wrap(ctx.Err(), "graphite: add metric")
						}
						continue
					}
					request := request.Clone(ctx)
					q := request.URL.Query()
					q.Set(queryParam, node.ID+".*")
					request.URL.RawQuery = q.Encode()
					select {
					case requests <- request:
						counter++
					case <-ctx.Done():
						return errors.Wrap(ctx.Err(), "graphite: add request")
					}
				}
			case <-ctx.Done():
				return errors.Wrap(ctx.Err(), "graphite: request manager process")
			}
		}
		return nil
	})

	pool, ctx := errgroup.WithContext(ctx)
	for range sequence.Simple(runtime.GOMAXPROCS(0)) {
		pool.Go(func() error {
			for {
				select {
				case request, ok := <-requests:
					if !ok {
						return nil
					}
					data, err := provider.fetch(request)
					if err != nil {
						return err
					}
					select {
					case result <- data:
					case <-ctx.Done():
						return errors.Wrap(ctx.Err(), "graphite: worker write process")
					}
				case <-ctx.Done():
					return errors.Wrap(ctx.Err(), "graphite: worker read process")
				}
			}
		})
	}
	if err := pool.Wait(); err != nil {
		close(result)
		return nil, err
	}

	return metrics, main.Wait()
}

func (provider *provider) fetch(request *http.Request) ([]dto, error) {
	response, err := provider.client.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "graphite: metrics fetch request")
	}
	defer safe.Close(response.Body, unsafe.Ignore)

	var nodes []dto
	if err := json.NewDecoder(response.Body).Decode(&nodes); err != nil {
		return nil, errors.Wrap(err, "graphite: decode metrics fetch response")
	}

	return nodes, nil
}
