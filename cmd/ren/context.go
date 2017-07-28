package main

import (
	"context"

	"gopkg.in/urfave/cli.v1"
)

type netCtx context.Context

type Context struct {
	*cli.Context
	netCtx
}

func newContext(c *cli.Context) (ctx *Context, err error) {
	ctx = &Context{
		Context: c,
		netCtx:  context.Background(),
	}

	return
}
