package main

import (
	"context"
	"time"

	logCtx "github.com/hamba/logger/v2/ctx"
	"github.com/hamba/pkg/v2/http"
	"github.com/hamba/pkg/v2/http/healthz"
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
	srv := http.NewHealthServer(ctx, http.HealthServerConfig{
		Addr:    addr,
		Handler: apiSrv,
		Stats:   obsvr.Stats,
		Log:     obsvr.Log,
	})

	if err = srv.AddHealthzChecks(healthz.PingHealth); err != nil {
		return err
	}

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
