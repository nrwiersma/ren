package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hamba/cmd/v2/observe"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer_HandleRenderImage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		url            string
		err            error
		wantPath       string
		wantData       interface{}
		wantStatusCode int
	}{
		{
			name:           "valid request",
			url:            "/foo/bar",
			wantPath:       "foo/bar.svg",
			wantData:       map[string]string{},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "valid request with data",
			url:            "/foo/bar?a=b",
			wantPath:       "foo/bar.svg",
			wantData:       map[string]string{"a": "b"},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "handles non-existent template",
			url:            "/foo/bar",
			err:            ren.ErrTemplateNotFound,
			wantPath:       "foo/bar.svg",
			wantData:       map[string]string{},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "handles application error",
			url:            "/foo/bar",
			err:            errors.New("test"),
			wantPath:       "foo/bar.svg",
			wantData:       map[string]string{},
			wantStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obsvr := observe.NewFake()

			app := &mockApp{}
			app.On("Render", test.wantPath, test.wantData).Return([]byte{}, test.err)

			srv := api.New(app, obsvr)

			r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, test.url, nil)
			w := httptest.NewRecorder()
			srv.ServeHTTP(w, r)

			assert.Equal(t, test.wantStatusCode, w.Code)
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
