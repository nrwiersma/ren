package ren_test

import (
	"testing"

	"github.com/hamba/cmd/v2/observe"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestApplication_Render(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		path    string
		data    map[string]string
		want    []byte
		wantErr require.ErrorAssertionFunc
	}{
		{
			name:    "renders a template",
			path:    "test.tmpl",
			data:    map[string]string{"str": "str"},
			want:    []byte("str"),
			wantErr: require.NoError,
		},
		{
			name:    "handles no data",
			path:    "test.tmpl",
			data:    map[string]string{},
			want:    []byte{},
			wantErr: require.NoError,
		},
		{
			name:    "handles non-existent template",
			path:    "nonexistent",
			data:    map[string]string{},
			wantErr: require.Error,
		},
		{
			name:    "handles template error",
			path:    "parse_err.tmpl",
			data:    map[string]string{},
			wantErr: require.Error,
		},
		{
			name:    "handles template exec error",
			path:    "exec_err.tmpl",
			data:    map[string]string{},
			wantErr: require.Error,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			r := reader.NewFileReader("testdata", otel.Tracer("reader"))
			app := newTestApplication(r)

			got, err := app.Render(t.Context(), test.path, test.data)

			test.wantErr(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestApplication_RenderNotFound(t *testing.T) {
	r := reader.NewFileReader("something-that-doesnt-exist", otel.Tracer("reader"))
	app := newTestApplication(r)

	_, err := app.Render(t.Context(), "", nil)

	assert.Equal(t, ren.ErrTemplateNotFound, err)
}

func TestApplication_IsHealthy(t *testing.T) {
	app := newTestApplication(nil)

	assert.NoError(t, app.IsHealthy())
}

func newTestApplication(r ren.Reader) *ren.Application {
	obsvr := observe.NewFake()

	return ren.NewApplication(r, obsvr)
}
