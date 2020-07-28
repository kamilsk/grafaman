package grafana

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/kamilsk/retry/v5"
	"github.com/kamilsk/retry/v5/backoff"
	"github.com/kamilsk/retry/v5/strategy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	"github.com/kamilsk/grafaman/internal/model"
)

// New returns an instance of Grafana dashboard provider.
func New(endpoint string, logger *logrus.Logger) (*provider, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "grafana: prepare dashboard provider endpoint URL")
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

// Fetch takes a dashboard JSON model and extracts queries and variables from it.
// Documentation: https://grafana.com/docs/grafana/latest/http_api/dashboard/#get-dashboard-by-uid.
func (provider *provider) Fetch(ctx context.Context, uid string) (*model.Dashboard, error) {
	const source = "/api/dashboards/uid/"

	u := provider.endpoint
	u.Path = path.Join(u.Path, source, uid)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "grafana: create dashboard base request")
	}

	var response *http.Response
	what := func(ctx context.Context) error {
		var err error
		logger := provider.logger.WithField("url", request.URL.String())
		logger.Info("start to fetch data")
		response, err = provider.client.Do(request) // nolint:bodyclose
		if err != nil {
			logger.WithError(err).Error("fail fetch data")
			return errors.Wrap(err, "grafana: dashboard fetch request")
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
		return nil, err
	}
	defer safe.Close(response.Body, unsafe.Ignore)

	var payload struct {
		Dashboard dashboard `json:"dashboard,omitempty"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, errors.Wrap(err, "grafana: decode dashboard fetch response")
	}

	result := model.Dashboard{
		RawData:   convertTargets(fetchTargets(payload.Dashboard.Panels)),
		Variables: convertVariables(fetchVariables(payload.Dashboard)),
	}
	return &result, nil
}
