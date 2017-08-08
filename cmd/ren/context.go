package main

import (
	"context"

	"github.com/msales/pkg/log"
	"gopkg.in/urfave/cli.v1"
)

type netCtx context.Context

type Context struct {
	*cli.Context
	netCtx

	logger log.Logger
}

func newContext(c *cli.Context) (ctx *Context, err error) {
	ctx = &Context{
		Context: c,
		netCtx:  context.Background(),
	}

	ctx.logger, err = newLogger(ctx)
	if err != nil {
		return
	}

	if ctx.logger != nil {
		ctx.netCtx = log.WithLogger(ctx.netCtx, ctx.logger)
	}

	return
}
