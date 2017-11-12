package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_Health(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 200, res.StatusCode);
}

func newTestServer(t *testing.T) *httptest.Server {
	ctx, fs := newTestContext()
	fs.String(FlagStats, "statsd://host", "")

	app, err := newApplication(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	srv := newServer(ctx, app)

	return httptest.NewServer(srv)
}
