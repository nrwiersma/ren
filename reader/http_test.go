package reader_test

import (
	"net/http"
	"testing"

	httptest "github.com/hamba/testutils/http"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpReader(t *testing.T) {
	r := reader.NewHTTPReader("http://test")

	assert.Implements(t, (*ren.Reader)(nil), r)
	assert.IsType(t, &reader.HTTPReader{}, r)
}

func TestHttpReader_Read(t *testing.T) {
	srv := httptest.NewServer(t)
	srv.On("GET", "/test.tmpl").ReturnsString(http.StatusOK, "{{ .str }}")
	defer srv.Close()

	r := reader.NewHTTPReader(srv.URL())

	str, err := r.Read("test.tmpl")

	assert.NoError(t, err)
	assert.Equal(t, "{{ .str }}", str)
}

func TestHttpReader_ReadGetError(t *testing.T) {
	r := reader.NewHTTPReader("http://")

	_, err := r.Read("test.tmpl")

	assert.Error(t, err)
}

func TestHttpReader_Read404(t *testing.T) {
	srv := httptest.NewServer(t)
	srv.On("GET", "/test.tmpl").ReturnsStatus(http.StatusNotFound)
	defer srv.Close()

	r := reader.NewHTTPReader(srv.URL())

	_, err := r.Read("test.tmpl")

	assert.Error(t, err)
}
