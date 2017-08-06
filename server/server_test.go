package server_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-zoo/bone"
	"github.com/nrwiersma/ren/server"
	"github.com/stretchr/testify/assert"
	"github.com/nrwiersma/ren"
)

func TestImageHandler_ServeHTTP(t *testing.T) {
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
		mux := bone.New()
		mux.Get("/:group/:file", server.ImageHandler{
			Render: func(p string, d interface{}) ([]byte, error) {
				assert.Equal(t, tt.path, p)
				assert.Equal(t, tt.data, d)

				return []byte{}, tt.err
			},
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", tt.url, nil)
		mux.ServeHTTP(w, req)

		assert.Equal(t, tt.code, w.Code)
	}
}

func TestHealthHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		err  error
		code int
	}{
		{nil, http.StatusOK},
		{errors.New(""), http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		h := server.HealthHandler{
			IsHealthy: func() error {
				return tt.err
			},
		}

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/health", nil)
		h.ServeHTTP(w, req)

		assert.Equal(t, tt.code, w.Code)
	}
}

func TestNotFoundHandler(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	server.NotFoundHandler().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
