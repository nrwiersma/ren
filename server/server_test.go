package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/server"
	"github.com/stretchr/testify/assert"
)

type testApp struct {
	render    func(path string, data interface{}) ([]byte, error)
	isHealthy func() error
}

func (a testApp) Render(path string, data interface{}) ([]byte, error) {
	return a.render(path, data)
}

func (a testApp) IsHealthy() error {
	return a.isHealthy()
}

func TestServer_ImageHandler(t *testing.T) {
	tests := []struct {
		url  string
		err  error
		path string
		data interface{}
		code int
	}{
		{"/foo/bar", nil, "foo/bar.svg", map[string]string{}, 200},
		{"/foo/bar?a=b", nil, "foo/bar.svg", map[string]string{"a": "b"}, 200},
		{"/foo/bar", ren.ErrTemplateNotFound, "foo/bar.svg", map[string]string{}, 404},
		{"/foo/bar", errors.New(""), "foo/bar.svg", map[string]string{}, 500},
		{"//", nil, ".svg", map[string]string{}, 200},
	}

	for _, tt := range tests {
		app := testApp{
			render: func(p string, d interface{}) ([]byte, error) {
				assert.Equal(t, tt.path, p)
				assert.Equal(t, tt.data, d)

				return []byte{}, tt.err
			},
		}
		srv := server.New(app)

		r := httptest.NewRequest("GET", tt.url, nil)
		w := httptest.NewRecorder()
		srv.ServeHTTP(w, r)

		assert.Equal(t, tt.code, w.Code)
	}
}

func TestServer_HealthHandler(t *testing.T) {
	tests := []struct {
		err  error
		code int
	}{
		{nil, http.StatusOK},
		{errors.New(""), http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		app := testApp{
			isHealthy: func() error {
				return tt.err
			},
		}
		srv := server.New(app)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/health", nil)
		srv.ServeHTTP(w, req)

		assert.Equal(t, tt.code, w.Code)
	}
}

func TestNotFoundHandler(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	server.NotFoundHandler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
