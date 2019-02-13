package reader_test

import (
	"testing"

	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
	"github.com/stretchr/testify/assert"
)

func TestNewFileReader(t *testing.T) {
	r := reader.NewFileReader("test")

	assert.Implements(t, (*ren.Reader)(nil), r)
	assert.IsType(t, &reader.FileReader{}, r)
}

func TestFileReader_Read(t *testing.T) {
	r := reader.NewFileReader("/../testdata")

	str, err := r.Read("test.tmpl")

	assert.NoError(t, err)
	assert.Equal(t, "{{ .str }}", str)
}

func TestFileReader_ReadFileNotFound(t *testing.T) {
	r := reader.NewFileReader("/../testdata")

	_, err := r.Read("wrong")

	assert.Error(t, err)
}

func TestFileReader_ReadDirectory(t *testing.T) {
	r := reader.NewFileReader("/../testdata")

	_, err := r.Read("")

	assert.Error(t, err)
}
