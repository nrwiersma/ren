package main

import (
	"context"
	"net/http"
	"time"

	"github.com/hamba/cmd/v2"
	logCtx "github.com/hamba/logger/v2/ctx"
	httpx "github.com/hamba/pkg/v2/http"
	"github.com/nrwiersma/ren/api"
	"github.com/urfave/cli/v2"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func runServer(c *cli.Context) error {
	ctx := c.Context

	log, err := cmd.NewLogger(c)
	if err != nil {
		return err
	}

	stats, err := cmd.NewStatter(c, log)
	if err != nil {
		return err
	}

	tracer, err := cmd.NewTracer(c, log,
		semconv.ServiceNameKey.String("ren"),
		semconv.ServiceVersionKey.String(version),
	)
	if err != nil {
		return err
	}
	defer func() { _ = tracer.Shutdown(context.Background()) }()

	app, err := newApplication(c, log, stats, tracer)
	if err != nil {
		return err
	}

	apiSrv := api.New(app, log, stats, tracer.Tracer("api"))

	mux := http.NewServeMux()
	mux.Handle("/readyz", http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	mux.Handle("/healthz", httpx.NewHealthHandler(app))
	mux.Handle("/", apiSrv)

	port := c.String(cmd.FlagPort)
	srv := httpx.NewServer(ctx, ":"+port, mux)

	log.Info("Starting server", logCtx.Str("port", port))
	srv.Serve(func(err error) {
		log.Error("Server error", logCtx.Error("error", err))
	})
	defer func() { _ = srv.Close() }()

	<-ctx.Done()

	log.Info("Shutting down")
	if err = srv.Shutdown(10 * time.Second); err != nil {
		log.Error("Failed to shutdown server", logCtx.Error("error", err))
	}

	return nil
}
