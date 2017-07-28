package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/msales/ren"
	"github.com/msales/ren/server"
	"github.com/msales/ren/server/middleware"
	"gopkg.in/urfave/cli.v1"
)

func runServer(c *cli.Context) {
	// Context
	ctx, err := newContext(c)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Application
	app, err := newApplication(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Server
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
