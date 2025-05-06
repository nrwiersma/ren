package reader_test

import (
	"testing"

	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestNewFileReader(t *testing.T) {
	r := reader.NewFileReader("test", otel.Tracer("file-render"))

	assert.Implements(t, (*ren.Reader)(nil), r)
	assert.IsType(t, &reader.FileReader{}, r)
}

func TestFileReader_Read(t *testing.T) {
	r := reader.NewFileReader("/../testdata", otel.Tracer("file-render"))

	str, err := r.Read(t.Context(), "test.tmpl")

	require.NoError(t, err)
	assert.Equal(t, "{{ .str }}", str)
}

func TestFileReader_ReadFileNotFound(t *testing.T) {
	r := reader.NewFileReader("/../testdata", otel.Tracer("file-render"))

	_, err := r.Read(t.Context(), "wrong")

	assert.Error(t, err)
}

func TestFileReader_ReadDirectory(t *testing.T) {
	r := reader.NewFileReader("/../testdata", otel.Tracer("file-render"))

	_, err := r.Read(t.Context(), "")

	assert.Error(t, err)
}
