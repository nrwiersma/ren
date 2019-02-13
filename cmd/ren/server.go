package main

import (
	"fmt"
	"net/http"

	"github.com/hamba/cmd"
	"github.com/hamba/pkg/httpx"
	"github.com/hamba/pkg/log"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/server"
	"github.com/nrwiersma/ren/server/middleware"
	"gopkg.in/urfave/cli.v1"
)

func runServer(c *cli.Context) {
	ctx, err := cmd.NewContext(c)
	if err != nil {
		panic(err)
	}

	app, err := newApplication(ctx)
	if err != nil {
		log.Fatal(ctx, err.Error())
	}

	port := c.String(cmd.FlagPort)
	s := newServer(app)
	log.Info(ctx, fmt.Sprintf("Starting server on port %s", port))
	if err := http.ListenAndServe(":"+port, s); err != nil {
		log.Fatal(ctx, "ren: server error", "error", err.Error())
	}
}

func newServer(app *ren.Application) http.Handler {
	health := httpx.NewHealthMux(app)
	srv := server.NewMux(app)
	mux := httpx.CombineMuxes(health, srv)

	return middleware.Common(mux, app)
}
