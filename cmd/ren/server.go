package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

	l := log15.New()
	l.SetHandler(h)

	return l, nil
}
