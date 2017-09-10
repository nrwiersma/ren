package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"net/url"

	"github.com/nrwiersma/pkg/stats"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/server"
	"github.com/nrwiersma/ren/server/middleware"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/urfave/cli.v1"
)

func runServer(c *cli.Context) {
	ctx, err := newContext(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	app, err := newApplication(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	port := c.String(FlagPort)
	s := newServer(ctx, app)
	fmt.Printf("Starting on port %s.\n", port)
	log.Fatal(http.ListenAndServe(":"+port, s))
}

func newApplication(c *Context) (*ren.Application, error) {
	app := ren.NewApplication(c.String(FlagTemplates))

	return app, nil
}

func newServer(ctx *Context, app *ren.Application) http.Handler {
	s := server.New(app)

	h := middleware.Common(s)
	return middleware.WithContext(h, ctx)
}

func newLogger(c *Context) (log15.Logger, error) {
	lvl := c.String(FlagLogLevel)
	v, err := log15.LvlFromString(lvl)
	if err != nil {
		return nil, err
	}

	h := log15.LvlFilterHandler(v, log15.StreamHandler(os.Stdout, log15.LogfmtFormat()))
	if lvl == "debug" {
		h = log15.CallerFileHandler(h)
	}

	l := log15.New()
	l.SetHandler(log15.LazyHandler(h))

	return l, nil
}

func newStats(c *Context) (stats.Stats, error) {
	dsn := c.String(FlagStats)
	if dsn == "" {
		return stats.Null, nil
	}

	uri, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch uri.Scheme {
	case "statsd":
		return newStatsdStats(uri.Host)

	default:
		return stats.Null, nil
	}
}

func newStatsdStats(addr string) (stats.Stats, error) {
	return stats.NewStatsd(addr, "ren")
}
