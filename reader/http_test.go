package reader_test

import (
	"context"
	"net/http"
	"testing"

	httptest "github.com/hamba/testutils/http"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestNewHttpReader(t *testing.T) {
	r, err := reader.NewHTTPReader("http://test", otel.Tracer("http-render"))

	require.NoError(t, err)
	assert.Implements(t, (*ren.Reader)(nil), r)
	assert.IsType(t, &reader.HTTPReader{}, r)
}

func TestHttpReader_Read(t *testing.T) {
	srv := httptest.NewServer(t)
	srv.On("GET", "/basepath/some/test.tmpl").ReturnsString(http.StatusOK, "{{ .str }}")
	defer srv.Close()

	r, err := reader.NewHTTPReader(srv.URL()+"/basepath/", otel.Tracer("http-render"))
	require.NoError(t, err)

	str, err := r.Read(context.Background(), "some/test.tmpl")

	assert.NoError(t, err)
	assert.Equal(t, "{{ .str }}", str)
}

func TestHttpReader_ReadGetError(t *testing.T) {
	r, err := reader.NewHTTPReader("http://", otel.Tracer("http-render"))
	require.NoError(t, err)

	_, err = r.Read(context.Background(), "test.tmpl")

	assert.Error(t, err)
}

func TestHttpReader_Read404(t *testing.T) {
	srv := httptest.NewServer(t)
	srv.On("GET", "/test.tmpl").ReturnsStatus(http.StatusNotFound)
	defer srv.Close()

	r, err := reader.NewHTTPReader(srv.URL(), otel.Tracer("http-render"))
	require.NoError(t, err)

	_, err = r.Read(context.Background(), "test.tmpl")

	assert.Error(t, err)
}
