package ren_test

import (
	"fmt"
	"testing"

	"github.com/hamba/pkg/log"
	"github.com/hamba/pkg/stats"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
	"github.com/stretchr/testify/assert"
)

func TestApplication_Render(t *testing.T) {
	tests := []struct {
		path   string
		data   interface{}
		ok     bool
		expect []byte
	}{
		{"test.tmpl", map[string]string{"str": "str"}, true, []byte("str")},
		{"test.tmpl", map[string]string{}, true, []byte{}},
		{"nonexistant", map[string]string{}, false, nil},
		{"parse_err.tmpl", map[string]string{}, false, nil},
		{"exec_err.tmpl", map[string]string{}, false, nil},
	}

	for i, tt := range tests {
		a := ren.NewApplication(log.Null, stats.Null)
		a.Reader = reader.NewFileReader("testdata")
		got, err := a.Render(tt.path, tt.data)
		if ok := (err == nil); ok != tt.ok {
			if err != nil {
				assert.FailNow(t, fmt.Sprintf("test %d; unexpected error: %statter", i, err))
			}
			assert.FailNow(t, fmt.Sprintf("test %d; unexpected success", i))
		}

		assert.Equal(t, tt.expect, got)
	}
}

func TestApplication_RenderNotFound(t *testing.T) {
	a := ren.NewApplication(log.Null, stats.Null)
	a.Reader = reader.NewFileReader("")

	_, err := a.Render("", nil)

	assert.Equal(t, ren.ErrTemplateNotFound, err)
}

func TestApplication_IsHealthy(t *testing.T) {
	a := ren.NewApplication(log.Null, stats.Null)

	assert.Nil(t, a.IsHealthy())
}

func TestApplication_Logger(t *testing.T) {
	a := ren.NewApplication(log.Null, stats.Null)

	assert.Equal(t, log.Null, a.Logger())
}

func TestApplication_Statter(t *testing.T) {
	a := ren.NewApplication(log.Null, stats.Null)

	assert.Equal(t, stats.Null, a.Statter())
}
