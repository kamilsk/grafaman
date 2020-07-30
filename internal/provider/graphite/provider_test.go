package graphite_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.octolab.org/safe"
	xtime "go.octolab.org/time"
	"go.octolab.org/unsafe"

	"github.com/kamilsk/grafaman/internal/model"
	. "github.com/kamilsk/grafaman/internal/provider/graphite"
)

func TestProvider(t *testing.T) {
	ctx := context.Background()

	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	t.Run("bad endpoint", func(t *testing.T) {
		provider, err := New(":invalid", nil, logger)
		assert.Error(t, err)
		assert.Nil(t, provider)
	})

	t.Run("success fetch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := NewMockClient(ctrl)
		client.EXPECT().
			Do(gomock.Any()).
			Return(response("testdata/success.1.json")) // nolint:bodyclose
		client.EXPECT().
			Do(gomock.Any()).
			Return(response("testdata/success.2.json")) // nolint:bodyclose
		client.EXPECT().
			Do(gomock.Any()).
			Return(response("testdata/success.3.json")) // nolint:bodyclose

		provider, err := New("test", client, logger)
		require.NoError(t, err)

		metrics, err := provider.Fetch(ctx, "apps.services.awesome-service", xtime.Day)
		assert.NoError(t, err)
		assert.Equal(t, model.Metrics{
			"apps.services.awesome-service.metric.a",
			"apps.services.awesome-service.metric.b",
			"apps.services.awesome-service.metric.c",
		}, metrics)
	})

	t.Run("nil context", func(t *testing.T) {
		provider, err := New("test", nil, logger)
		require.NoError(t, err)

		metrics, err := provider.Fetch(nil, "apps.services.awesome-service", xtime.Day) // nolint:staticcheck
		assert.Error(t, err)
		assert.Nil(t, metrics)
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

		metrics, err := provider.Fetch(ctx, "apps.services.awesome-service", xtime.Day)
		assert.Error(t, err)
		assert.Nil(t, metrics)
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

		metrics, err := provider.Fetch(ctx, "apps.services.awesome-service", xtime.Day)
		assert.Error(t, err)
		assert.Nil(t, metrics)
	})

	t.Run("context deadline exceeded", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(ctx)

		client := NewMockClient(ctrl)
		client.EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(*http.Request) (*http.Response, error) {
				time.Sleep(15 * time.Millisecond)
				cancel()
				return response("testdata/success.1.json")
			})

		provider, err := New("test", client, logger)
		require.NoError(t, err)

		metrics, err := provider.Fetch(ctx, "apps.services.awesome-service", xtime.Day)
		assert.Error(t, err)
		assert.Nil(t, metrics)
		time.Sleep(15 * time.Millisecond)
	})

	t.Run("parallel with deadline", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx, cancel := context.WithCancel(ctx)

		client := NewMockClient(ctrl)
		client.EXPECT().
			Do(gomock.Any()).
			Return(response("testdata/parallel.1.json")) // nolint:bodyclose
		client.EXPECT().
			Do(gomock.Any()).
			Return(response("testdata/parallel.2.json")) // nolint:bodyclose
		client.EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(*http.Request) (*http.Response, error) {
				defer cancel()
				return response("testdata/parallel.3-1.json")
			}).
			AnyTimes()
		client.EXPECT().
			Do(gomock.Any()).
			DoAndReturn(func(*http.Request) (*http.Response, error) {
				defer cancel()
				return response("testdata/parallel.3-2.json")
			}).
			AnyTimes()

		provider, err := New("test", client, logger)
		require.NoError(t, err)

		metrics, err := provider.Fetch(ctx, "apps.services.awesome-service", xtime.Day)
		assert.Error(t, err)
		assert.Nil(t, metrics)
		time.Sleep(20 * time.Millisecond)
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
