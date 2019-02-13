package reader_test

import (
	"net/http"
	"testing"

	"github.com/hamba/pkg/httpx/httptest"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
	"github.com/stretchr/testify/assert"
)

func TestNewHttpReader(t *testing.T) {
	r := reader.NewHttpReader("http://test")

	assert.Implements(t, (*ren.Reader)(nil), r)
	assert.IsType(t, &reader.HttpReader{}, r)
}

func TestHttpReader_Read(t *testing.T) {
	srv := httptest.NewServer(t)
	srv.On("GET", "/test.tmpl").ReturnsString(http.StatusOK, "{{ .str }}")
	defer srv.Close()

	r := reader.NewHttpReader(srv.URL())

	str, err := r.Read("test.tmpl")

	assert.NoError(t, err)
	assert.Equal(t, "{{ .str }}", str)
}

func TestHttpReader_ReadGetError(t *testing.T) {
	r := reader.NewHttpReader("http://")

	_, err := r.Read("test.tmpl")

	assert.Error(t, err)
}

func TestHttpReader_Read404(t *testing.T) {
	srv := httptest.NewServer(t)
	srv.On("GET", "/test.tmpl").ReturnsStatus(http.StatusNotFound)
	defer srv.Close()

	r := reader.NewHttpReader(srv.URL())

	_, err := r.Read("test.tmpl")

	assert.Error(t, err)
}
