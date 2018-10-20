package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/msales/pkg/clix"
	"github.com/msales/pkg/log"
	"github.com/msales/pkg/stats"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
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

func newTestContext() (*clix.Context, *flag.FlagSet) {
	fs := new(flag.FlagSet)
	c := cli.NewContext(cli.NewApp(), fs, nil)

	ctx, _ := clix.NewContext(c, clix.Stats(stats.Null), clix.Logger(log.Null))

	return ctx, fs
}

func newTestServer(t *testing.T) *httptest.Server {
	ctx, _ := newTestContext()

	app, err := newApplication(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}

	srv := newServer(ctx, app)

	return httptest.NewServer(srv)
}
