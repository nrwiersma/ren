package main

import (
	"context"
	"fmt"

	lctx "github.com/hamba/logger/v2/ctx"
	"github.com/hamba/pkg/v2/http/server"
	"github.com/nrwiersma/ren/api"
	"github.com/urfave/cli/v2"
)

func runServer(c *cli.Context) error {
	ctx, cancel := context.WithCancel(c.Context)
	defer cancel()

	obsvr, err := newObserver(c)
	if err != nil {
		return err
	}
	defer obsvr.Close()

	app, err := newApplication(c, obsvr)
	if err != nil {
		return err
	}

	apiSrv := api.New(app, obsvr)

	addr := c.String(flagAddr)
	srv := &server.GenericServer[context.Context]{
		Addr:    addr,
		Handler: apiSrv,
		Stats:   obsvr.Stats,
		Log:     obsvr.Log,
	}

	obsvr.Log.Info("Starting server", lctx.Str("address", addr))

	if err = srv.Run(ctx); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}
