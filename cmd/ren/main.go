package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/ettle/strcase"
	"github.com/hamba/cmd/v2"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

const flagTemplates = "templates"

var version = "¯\\_(ツ)_/¯"

var commands = []*cli.Command{
	{
		Name:  "server",
		Usage: "Run the ren HTTP server",
		Flags: cmd.Flags{
			&cli.StringFlag{
				Name:    flagTemplates,
				Value:   "file:///./templates",
				Usage:   "The URI to the templates. Supported schemes: 'file', 'http', 'https'.",
				EnvVars: []string{strcase.ToSNAKE(flagTemplates)},
			},
		}.Merge(cmd.MonitoringFlags, cmd.ServerFlags),
		Action: runServer,
	},
}

func main() {
	os.Exit(realMain())
}

func realMain() (code int) {
	ui := newTerm()

	defer func() {
		if v := recover(); v != nil {
			ui.Error(fmt.Sprintf("Panic: %v\n%s", v, string(debug.Stack())))
			code = 1
			return
		}
	}()

	app := cli.NewApp()
	app.Name = "ren"
	app.Version = version
	app.Commands = commands

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.RunContext(ctx, os.Args); err != nil {
		ui.Error(err.Error())
		return 1
	}
	return 0
}
