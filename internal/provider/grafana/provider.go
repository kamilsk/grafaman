package grafana

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/pkg/errors"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	entity "github.com/kamilsk/grafaman/internal/provider"
)

// New returns an instance of Grafana dashboard provider.
func New(endpoint string) (*provider, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, errors.Wrap(err, "prepare Grafana dashboard provider endpoint URL")
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

func (provider *provider) Fetch(ctx context.Context, uid string) (*entity.Dashboard, error) {
	const method = "/api/dashboards/uid/"

	u := provider.endpoint
	u.Path = path.Join(u.Path, method, uid)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "create Grafana dashboard base request")
	}

	resp, err := provider.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Grafana dashboard fetch request")
	}
	defer safe.Close(resp.Body, unsafe.Ignore)

	var payload struct {
		Dashboard dashboard `json:"dashboard"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, errors.Wrap(err, "decode Grafana dashboard fetch response")
	}

	result := entity.Dashboard{
		RawData:   convertTargets(fetchTargets(payload.Dashboard.Panels)),
		Variables: convertVariables(fetchVariables(payload.Dashboard)),
	}
	return &result, nil
}
