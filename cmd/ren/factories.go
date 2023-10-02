package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/hamba/cmd/v2"
	"github.com/hamba/cmd/v2/observe"
	"github.com/hamba/cmd/v2/term"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
	"github.com/urfave/cli/v2"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

func newTerm() term.Term {
	return term.Prefixed{
		ErrorPrefix: "Error: ",
		Term: term.Colored{
			ErrorColor: term.Red,
			Term: term.Basic{
				Writer:      os.Stdout,
				ErrorWriter: os.Stderr,
				Verbose:     false,
			},
		},
	}
}

func newObserver(c *cli.Context) (*observe.Observer, error) {
	log, err := cmd.NewLogger(c)
	if err != nil {
		return nil, err
	}

	stats, err := cmd.NewStatter(c, log)
	if err != nil {
		return nil, err
	}

	tracer, err := cmd.NewTracer(c, log,
		semconv.ServiceNameKey.String("ren"),
		semconv.ServiceVersionKey.String(version),
	)
	if err != nil {
		return nil, err
	}
	tracerCancel := func() { _ = tracer.Shutdown(context.Background()) }

	return observe.New(log, stats, tracer, tracerCancel), nil
}

func newApplication(c *cli.Context, obsvr *observe.Observer) (*ren.Application, error) {
	r, err := newReader(c.String(flagTemplates), obsvr.Tracer("reader"))
	if err != nil {
		return nil, err
	}

	return ren.NewApplication(r, obsvr), nil
}

func newReader(uri string, tracer trace.Tracer) (ren.Reader, error) {
	if uri == "" {
		return nil, errors.New("template uri required")
	}

	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "file":
		return reader.NewFileReader(u.Path, tracer), nil

	case "http", "https":
		return reader.NewHTTPReader(uri, tracer)

	default:
		return nil, fmt.Errorf("unsupported template scheme: %s", u.Scheme)
	}
}
