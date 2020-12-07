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

		listener := NewMockProgressListener(ctrl)
		listener.EXPECT().OnStepDone().Times(3)
		listener.EXPECT().OnStepQueued().Times(3)

		provider, err := New("test", client, logger, listener)
		require.NoError(t, err)

		metrics, err := provider.Fetch(ctx, "apps.services.awesome-service", xtime.Day)
		assert.NoError(t, err)
		assert.Equal(t, model.Metrics{
			"apps.services.awesome-service.metric.a",
			"apps.services.awesome-service.metric.b",
			"apps.services.awesome-service.metric.c",
		}, metrics)
	})

	t.Run("bad endpoint", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		listener := NewMockProgressListener(ctrl)
		listener.EXPECT().OnStepDone().Times(0)
		listener.EXPECT().OnStepQueued().Times(0)

		provider, err := New(":invalid", nil, logger, listener)
		assert.Error(t, err)
		assert.Nil(t, provider)
	})

	t.Run("nil context", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		listener := NewMockProgressListener(ctrl)
		listener.EXPECT().OnStepDone().Times(0)
		listener.EXPECT().OnStepQueued().Times(0)

		provider, err := New("test", nil, logger, listener)
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

		listener := NewMockProgressListener(ctrl)
		listener.EXPECT().OnStepDone().Times(1)
		listener.EXPECT().OnStepQueued().Times(1)

		provider, err := New("test", client, logger, listener)
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

		listener := NewMockProgressListener(ctrl)
		listener.EXPECT().OnStepDone().Times(1)
		listener.EXPECT().OnStepQueued().Times(1)

		provider, err := New("test", client, logger, listener)
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

		listener := NewMockProgressListener(ctrl)
		listener.EXPECT().OnStepDone().Times(2)
		listener.EXPECT().OnStepQueued().Times(2)

		provider, err := New("test", client, logger, listener)
		require.NoError(t, err)

		metrics, err := provider.Fetch(ctx, "apps.services.awesome-service", xtime.Day)
		assert.Error(t, err)
		assert.Len(t, metrics, 0)
		time.Sleep(15 * time.Millisecond)
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
