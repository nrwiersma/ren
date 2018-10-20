package main

import (
	"fmt"
	"net/http"

	"github.com/msales/pkg/clix"
	"github.com/msales/pkg/log"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/server"
	"github.com/nrwiersma/ren/server/middleware"
	"github.com/urfave/cli"
)

func runServer(c *cli.Context) {
	ctx, err := clix.NewContext(c)
	if err != nil {
		panic(err)
	}

	app, err := newApplication(ctx)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}

	port := c.String(clix.FlagPort)
	s := newServer(ctx, app)
	log.Info(ctx, fmt.Sprintf("Starting server on port %s", port))
	if err := http.ListenAndServe(":"+port, s); err != nil {
		log.Fatal(ctx, "ren: server error", "error", err.Error())
	}
}

func newServer(ctx *clix.Context, app *ren.Application) http.Handler {
	s := server.New(app)

	h := middleware.Common(s)
	return middleware.WithContext(ctx, h)
}
