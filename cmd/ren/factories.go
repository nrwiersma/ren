package main

import (
	"fmt"
	"net/url"

	"github.com/hamba/cmd"
	"github.com/nrwiersma/ren"
	"github.com/nrwiersma/ren/reader"
)

// Application =============================

func newApplication(c *cmd.Context) (*ren.Application, error) {
	app := ren.NewApplication(
		c.Logger(),
		c.Statter(),
	)

	r, err := newReader(c.String(flagTemplates))
	if err != nil {
		return nil, err
	}
	app.Reader = r

	return app, nil
}

func newReader(dsn string) (ren.Reader, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "file":
		return reader.NewFileReader(u.Path), nil

	case "http", "https":
		return reader.NewHttpReader(dsn), nil

	default:
		return nil, fmt.Errorf("ren: unsupported template sheme: %s", u.Scheme)
	}
}
