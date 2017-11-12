package main

import (
	"flag"
	"testing"

	"github.com/msales/pkg/stats"
	"github.com/stretchr/testify/assert"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/urfave/cli.v1"
)

func TestContext(t *testing.T) {
	fs := new(flag.FlagSet)
	fs.String(FlagLogLevel, "debug", "")
	fs.String(FlagStats, "statsd://localhost:6000", "")
	ctx := cli.NewContext(cli.NewApp(), fs, nil)

	c, err := newContext(ctx)
	assert.NoError(t, err)

	assert.IsType(t, log15.New(), c.logger)
	assert.IsType(t, &stats.Statsd{}, c.stats)
}
