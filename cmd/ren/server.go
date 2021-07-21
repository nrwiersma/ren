package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/hamba/cmd/v2"
	"github.com/hamba/cmd/v2/term"
	"github.com/hamba/logger/v2"
	logCtx "github.com/hamba/logger/v2/ctx"
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
		cancel := newServer(ctx, ":"+port, apiHdlr, log)
		defer func() { _ = cancel() }()

		<-ctx.Done()

		log.Info("Shutting down")

		return nil
	}
}

func newServer(ctx context.Context, addr string, h http.Handler, log *logger.Logger) func() error {
	srv := http.Server{
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		Addr:    addr,
		Handler: h,
	}

	go func() {
		log.Info("Starting server", logCtx.Str("addr", addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server error", logCtx.Error("error", err))
		}
	}()

	return func() error {
		stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(stopCtx); err != nil {
			log.Warn("Failed to shutdown server", logCtx.Error("error", err))
		}

		return srv.Close()
	}
}
