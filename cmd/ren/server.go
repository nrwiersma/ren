package main

import (
	"context"
	"fmt"

	"github.com/hamba/cmd/v3/observe"
	lctx "github.com/hamba/logger/v2/ctx"
	"github.com/hamba/pkg/v2/http/server"
	"github.com/nrwiersma/ren/api"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func runServer(ctx context.Context, cmd *cli.Command) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	obsvr, err := observe.New(ctx, cmd, "ren", &observe.Options{
		StatsRuntime: true,
		TracingAttrs: []attribute.KeyValue{semconv.ServiceVersionKey.String(version)},
	})
	if err != nil {
		return err
	}
	defer obsvr.Close()

	app, err := newApplication(cmd, obsvr)
	if err != nil {
		return err
	}

	apiSrv := api.New(app, obsvr)

	addr := cmd.String(flagAddr)
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
