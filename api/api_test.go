package api_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hamba/logger/v2"
	"github.com/hamba/statter/v2"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel"
)

func TestServer_ImageHandler(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		err      error
		wantPath string
		wantData interface{}
		wantCode int
	}{
		{
			name:     "valid request",
			url:      "/foo/bar",
			wantPath: "foo/bar.svg",
			wantData: map[string]string{},
			wantCode: http.StatusOK,
		},
		{
			name:     "valid request with data",
			url:      "/foo/bar?a=b",
			wantPath: "foo/bar.svg",
			wantData: map[string]string{"a": "b"},
			wantCode: http.StatusOK,
		},
		{
			name:     "handles non-existent template",
			url:      "/foo/bar",
			err:      ren.ErrTemplateNotFound,
			wantPath: "foo/bar.svg",
			wantData: map[string]string{},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "handles application error",
			url:      "/foo/bar",
			err:      errors.New("test"),
			wantPath: "foo/bar.svg",
			wantData: map[string]string{},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "handles not path",
			url:      "//",
			err:      nil,
			wantPath: ".svg",
			wantData: map[string]string{},
			wantCode: http.StatusOK,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			log := logger.New(io.Discard, logger.LogfmtFormat(), logger.Error)
			stats := statter.New(statter.DiscardReporter, time.Second)
			tracer := otel.Tracer("app")

			app := &mockApp{}
			app.On("Render", test.wantPath, test.wantData).Return([]byte{}, test.err)
			mux := api.New(app, log, stats, tracer)

			r := httptest.NewRequest("GET", test.url, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, r)

			assert.Equal(t, test.wantCode, w.Code)
			app.AssertExpectations(t)
		})
	}
}

type mockApp struct {
	mock.Mock
}

func (a *mockApp) Render(_ context.Context, path string, data map[string]string) ([]byte, error) {
	args := a.Called(path, data)

	b := args.Get(0)
	if b == nil {
		return nil, args.Error(1)
	}
	return b.([]byte), args.Error(1)
}

func (a *mockApp) IsHealthy() error {
	args := a.Called()

	return args.Error(0)
}
