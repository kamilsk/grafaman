package grafana_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.octolab.org/safe"
	"go.octolab.org/unsafe"

	. "github.com/kamilsk/grafaman/internal/provider/grafana"
)

func TestProvider(t *testing.T) {
	ctx := context.Background()
	_ = ctx

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	t.Run("success fetch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)
		client.EXPECT().
			Do(gomock.Any()).
			Return(response("testdata/success.json")) // nolint:bodyclose

		provider, err := New("test", client, logger)
		require.NoError(t, err)

		dashboard, err := provider.Fetch(ctx, "dashboard")
		assert.NoError(t, err)
		assert.NotNil(t, dashboard)
	})

	t.Run("bad endpoint", func(t *testing.T) {
		provider, err := New(":invalid", nil, logger)
		assert.Error(t, err)
		assert.Nil(t, provider)
	})

	t.Run("nil context", func(t *testing.T) {
		provider, err := New("test", nil, logger)
		require.NoError(t, err)

		dashboard, err := provider.Fetch(nil, "dashboard") // nolint:staticcheck
		assert.Error(t, err)
		assert.Nil(t, dashboard)
	})

	t.Run("service unavailable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)
		client.EXPECT().
			Do(gomock.Any()).
			Return(nil, errors.New(http.StatusText(http.StatusServiceUnavailable)))

		provider, err := New("test", client, logger)
		require.NoError(t, err)

		dashboard, err := provider.Fetch(ctx, "dashboard")
		assert.Error(t, err)
		assert.Nil(t, dashboard)
	})

	t.Run("bad response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)
		client.EXPECT().
			Do(gomock.Any()).
			Return(response("testdata/invalid.json")) // nolint:bodyclose

		provider, err := New("test", client, logger)
		require.NoError(t, err)

		dashboard, err := provider.Fetch(ctx, "dashboard")
		assert.Error(t, err)
		assert.Nil(t, dashboard)
	})
}

// helpers

func response(filename string) (*http.Response, error) {
	resp := new(http.Response)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer safe.Close(file, unsafe.Ignore)

	var dto struct {
		Code int             `json:"code,omitempty"`
		Body json.RawMessage `json:"body,omitempty"`
	}
	if err := json.NewDecoder(file).Decode(&dto); err != nil {
		return nil, err
	}

	resp.StatusCode = dto.Code
	resp.Body = ioutil.NopCloser(bytes.NewReader(dto.Body))
	return resp, nil
}
