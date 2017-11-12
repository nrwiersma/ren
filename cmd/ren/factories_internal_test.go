package main

import (
	"context"
	"flag"
	"testing"

	"github.com/msales/pkg/stats"
	"github.com/stretchr/testify/assert"
	"gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/urfave/cli.v1"
)

func TestLogger(t *testing.T) {
	ctx, fs := newTestContext()
	fs.String(FlagLogLevel, "debug", "")

	l, err := newLogger(ctx)
	assert.NoError(t, err)

	assert.IsType(t, log15.New(), l)
}

func TestLogger_InvlidLevel(t *testing.T) {
	ctx, fs := newTestContext()
	fs.String(FlagLogLevel, "test", "")

	_, err := newLogger(ctx)
	assert.Error(t, err)
}

func TestStats_DsnNull(t *testing.T) {
	ctx, _ := newTestContext()

	s, err := newStats(ctx)
	assert.NoError(t, err)

	assert.IsType(t, stats.Null, s)
}

func TestStats_Statsd(t *testing.T) {
	ctx, fs := newTestContext()
	fs.String(FlagStats, "statsd://localhost:6000", "")

	s, err := newStats(ctx)
	assert.NoError(t, err)

	assert.IsType(t, &stats.Statsd{}, s)

}

func TestStats_L2met(t *testing.T) {
	ctx, fs := newTestContext()
	fs.String(FlagStats, "l2met://", "")

	s, err := newStats(ctx)
	assert.NoError(t, err)

	assert.IsType(t, &stats.L2met{}, s)

}

func TestStats_Null(t *testing.T) {
	ctx, fs := newTestContext()
	fs.String(FlagStats, "null://", "")

	s, err := newStats(ctx)
	assert.NoError(t, err)

	assert.IsType(t, stats.Null, s)

}

func newTestContext() (*Context, *flag.FlagSet) {
	fs := new(flag.FlagSet)
	ctx := cli.NewContext(cli.NewApp(), fs, nil)

	return &Context{
		Context: ctx,
		netCtx:  context.Background(),
	}, fs
}
