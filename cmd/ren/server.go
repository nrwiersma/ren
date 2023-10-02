package main

import (
	"context"
	"net/http"
	"time"

	logCtx "github.com/hamba/logger/v2/ctx"
	httpx "github.com/hamba/pkg/v2/http"
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

	mux := http.NewServeMux()
	mux.Handle("/readyz", httpx.OKHandler())
	mux.Handle("/healthz", httpx.NewHealthHandler(app))
	mux.Handle("/", apiSrv)

	addr := c.String(flagAddr)
	srv := httpx.NewServer(ctx, addr, mux)

	obsvr.Log.Info("Starting server", logCtx.Str("address", addr))
	srv.Serve(func(err error) {
		obsvr.Log.Error("Server error", logCtx.Error("error", err))
		cancel()
	})
	defer func() { _ = srv.Close() }()

	<-ctx.Done()

	obsvr.Log.Info("Shutting down")
	if err = srv.Shutdown(10 * time.Second); err != nil {
		obsvr.Log.Error("Failed to shutdown server", logCtx.Error("error", err))
	}

	return nil
}
