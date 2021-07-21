package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/hamba/cmd/v2/term"
	"github.com/hamba/logger/v2"
	"github.com/hamba/statter/v2"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/trace"
)

// UI ======================================

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

// Application =============================

func newApplication(
	c *cli.Context,
	log *logger.Logger,
	stats *statter.Statter,
	tp trace.TracerProvider,
) (*ren.Application, error) {
	r, err := newReader(c.String(flagTemplates), tp.Tracer("reader"))
	if err != nil {
		return nil, err
	}

	return ren.NewApplication(r, log, stats, tp), nil
}

// Reader ==================================

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
