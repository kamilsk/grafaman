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

	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/backoff"
	"github.com/kamilsk/retry/v5/strategy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.octolab.org/safe"
	"go.octolab.org/sequence"
	"go.octolab.org/unsafe"
	"golang.org/x/sync/errgroup"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

// New returns an instance of Graphite metrics provider.
func New(endpoint string, logger *logrus.Logger) (*provider, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "graphite: prepare metrics provider endpoint URL")
	}
	return &provider{
		client:   &http.Client{Timeout: time.Second},
		endpoint: *u,
		logger:   logger,
	}, nil
}

type provider struct {
	client   *http.Client
	endpoint url.URL
	logger   *logrus.Logger
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
		factor     = runtime.GOMAXPROCS(0)
		aggregator = make(chan dto, factor)
		result     = make(chan []dto, factor)
		requests   = make(chan *http.Request, 10)
		metrics    = make(entity.Metrics, 0, 512)
	)

	main, ctx := errgroup.WithContext(ctx)
	main.Go(func() error {
		defer provider.logger.Info("aggregator done")
		for {
			select {
			case node, ok := <-aggregator:
				if !ok {
					return nil
				}
				metrics = append(metrics, entity.Metric(node.ID))
			case <-ctx.Done():
				provider.logger.WithError(err).Error("aggregator process timeout")
				return errors.Wrap(ctx.Err(), "graphite: aggregator process")
			}
		}
	})
	main.Go(func() error {
		defer provider.logger.Info("request manager done")
		defer close(aggregator)
		defer close(requests)

		buffer := make([]*http.Request, 0, 1024)
		buffer = append(buffer, request)

		for len(buffer) > 0 {
			provider.logger.WithField("len", len(buffer)).Debug("loop")

			min := len(buffer)
			if limit := cap(requests); min > limit {
				min = limit
			}

			for range sequence.Simple(min) {
				request, buffer = buffer[len(buffer)-1], buffer[:len(buffer)-1]
				select {
				case requests <- request:
				case <-ctx.Done():
					provider.logger.WithError(err).Error("request manager process timeout")
					return errors.Wrap(ctx.Err(), "graphite: request manager process")
				}
			}

			for range sequence.Simple(min) {
				select {
				case nodes, ok := <-result:
					if !ok {
						return nil
					}
					for _, node := range nodes {
						if node.Leaf == 1 {
							select {
							case <-ctx.Done():
								provider.logger.WithError(err).Error("add metric timeout")
								return errors.Wrap(ctx.Err(), "graphite: add metric")
							case aggregator <- node:
							}
							continue
						}
						request := request.Clone(ctx)
						q := request.URL.Query()
						q.Set(queryParam, node.ID+".*")
						request.URL.RawQuery = q.Encode()
						buffer = append(buffer, request)
					}
				case <-ctx.Done():
					provider.logger.WithError(err).Error("request manager process timeout")
					return errors.Wrap(ctx.Err(), "graphite: request manager process")
				}
			}
		}
		return nil
	})

	pool, ctx := errgroup.WithContext(ctx)
	for range sequence.Simple(factor) {
		pool.Go(func() error {
			defer provider.logger.Info("worker done")
			for {
				select {
				case request, ok := <-requests:
					if !ok {
						return nil
					}
					data, err := provider.fetch(request)
					if err != nil {
						provider.logger.WithError(err).Error("worker crashed")
						return err
					}
					select {
					case <-ctx.Done():
						provider.logger.WithError(err).Error("worker write process timeout")
						return errors.Wrap(ctx.Err(), "graphite: worker write process")
					case result <- data:
					}
				case <-ctx.Done():
					provider.logger.WithError(err).Error("worker read process timeout")
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
	var response *http.Response

	what := func(ctx context.Context) error {
		var err error
		logger := provider.logger.WithField("url", request.URL.String())
		logger.Info("start to fetch data")
		response, err = provider.client.Do(request) // nolint:bodyclose
		if err != nil {
			logger.WithError(err).Error("fail fetch data")
			return errors.Wrap(err, "graphite: metrics fetch request")
		}
		logger.Info("success fetch data")
		return nil
	}
	how := retry.How{
		strategy.Limit(3),
		strategy.Backoff(
			backoff.Linear(50 * time.Millisecond),
		),
		strategy.CheckError(
			strategy.NetworkError(strategy.Strict),
		),
	}
	if err := retry.Do(request.Context(), what, how...); err != nil {
		provider.logger.WithError(err).Error("fail do request")
		return nil, err
	}
	defer safe.Close(response.Body, unsafe.Ignore)

	var nodes []dto
	if err := json.NewDecoder(response.Body).Decode(&nodes); err != nil {
		provider.logger.WithError(err).Error("fail decode response")
		return nil, errors.Wrap(err, "graphite: decode metrics fetch response")
	}

	return nodes, nil
}
