package server_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/hamba/pkg/log"
	"github.com/hamba/pkg/stats"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/server"
	"github.com/stretchr/testify/assert"
)

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
		mux := server.NewMux(app)

		r := httptest.NewRequest("GET", tt.url, nil)
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, r)

		assert.Equal(t, tt.code, w.Code)
	}
}

type testApp struct {
	render func(path string, data interface{}) ([]byte, error)
}

func (a testApp) Render(path string, data interface{}) ([]byte, error) {
	return a.render(path, data)
}

func (a testApp) Logger() log.Logger {
	return log.Null
}

func (a testApp) Statter() stats.Statter {
	return stats.Null
}
