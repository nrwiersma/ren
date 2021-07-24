package main

import (
	"context"
	"time"

	"github.com/hamba/cmd/v2"
	"github.com/hamba/cmd/v2/term"
	logCtx "github.com/hamba/logger/v2/ctx"
	"github.com/hamba/pkg/v2/http"
	"github.com/nrwiersma/ren/api"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/semconv"
)

func runServer(_ term.Term) func(c *cli.Context) error {
	return func(c *cli.Context) error {
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

		apiHdlr := api.New(app, log, stats, tracer.Tracer("api"))

		port := c.String(cmd.FlagPort)
		srv := http.NewServer(ctx, ":"+port, apiHdlr)

		log.Info("Starting server", logCtx.Str("port", port))
		srv.Serve(func(err error) {
			log.Error("Server error", logCtx.Error("error", err))
		})
		defer func() { _ = srv.Close() }()

		<-ctx.Done()

		log.Info("Shutting down")

		if err = srv.Shutdown(10 * time.Second); err != nil {
			log.Warn("Failed to shutdown server", logCtx.Error("error", err))
		}

		return nil
	}
}
