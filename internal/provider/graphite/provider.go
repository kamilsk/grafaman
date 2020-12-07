package graphite

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"sync"
	"time"

	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/backoff"
	"github.com/kamilsk/retry/v5/strategy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"
	"golang.org/x/sync/errgroup"

	"github.com/kamilsk/grafaman/internal/model"
)

// New returns an instance of Graphite metrics provider.
func New(endpoint string, client Client, logger *logrus.Logger, listener ProgressListener) (*provider, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "graphite: prepare metrics provider endpoint URL")
	}
	return &provider{
		client:   client,
		endpoint: *u,
		logger:   logger,
		listener: listener,
	}, nil
}

type provider struct {
	client   Client
	endpoint url.URL
	logger   *logrus.Logger
	listener ProgressListener
}

// Fetch walks through the endpoint and takes all metrics with the specified prefix.
// Documentation: https://graphite-api.readthedocs.io/en/latest/api.html#metrics-find.
func (provider *provider) Fetch(ctx context.Context, prefix string, last time.Duration) (model.Metrics, error) {
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

	return provider.recursive(request)
}

func (provider *provider) recursive(request *http.Request) (model.Metrics, error) {
	nodes, err := provider.fetch(request)
	if err != nil {
		return nil, err
	}

	var guard sync.Mutex
	metrics := make(model.Metrics, 0, 1<<4)

	group, ctx := errgroup.WithContext(request.Context())
	for _, node := range nodes {
		if node.Leaf == 1 {
			guard.Lock()
			metrics = append(metrics, model.Metric(node.ID))
			guard.Unlock()

			continue
		}

		node := node
		group.Go(func() error {
			request := request.Clone(ctx)
			q := request.URL.Query()
			q.Set(queryParam, node.ID+".*")
			request.URL.RawQuery = q.Encode()

			data, err := provider.recursive(request)
			if err != nil {
				return err
			}

			guard.Lock()
			metrics = append(metrics, data...)
			guard.Unlock()

			return nil
		})
	}

	return metrics, group.Wait()
}

func (provider *provider) fetch(request *http.Request) ([]dto, error) {
	provider.listener.OnStepQueued()
	defer provider.listener.OnStepDone()

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
