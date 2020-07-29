package graphite_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

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
			Return(response("testdata/success.1.json"))
		client.EXPECT().
			Do(gomock.Any()).
			Return(response("testdata/success.2.json"))
		client.EXPECT().
			Do(gomock.Any()).
			Return(response("testdata/success.3.json"))

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

		metrics, err := provider.Fetch(nil, "apps.services.awesome-service", xtime.Day)
		assert.Error(t, err)
		assert.Nil(t, metrics)
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
