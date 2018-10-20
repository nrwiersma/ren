package main

import (
	"github.com/msales/pkg/clix"
	"github.com/nrwiersma/ren"
)

// Application =============================

func newApplication(c *clix.Context) (*ren.Application, error) {
	app := ren.NewApplication(c.String(FlagTemplates))

	return app, nil
}
