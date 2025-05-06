package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/ettle/strcase"
	"github.com/hamba/cmd/v3"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v3"
)

const (
	flagAddr      = "addr"
	flagTemplates = "templates"
)

var version = "¯\\_(ツ)_/¯"

var commands = []*cli.Command{
	{
		Name:  "server",
		Usage: "Run the ren HTTP server",
		Flags: cmd.Flags{
			&cli.StringFlag{
				Name:    flagAddr,
				Value:   ":8080",
				Usage:   "The address to listen to.",
				Sources: cli.EnvVars(strcase.ToSNAKE(flagAddr)),
			},
			&cli.StringFlag{
				Name:    flagTemplates,
				Value:   "file:///./templates",
				Usage:   "The URI to the templates. Supported schemes: 'file', 'http', 'https'.",
				Sources: cli.EnvVars(strcase.ToSNAKE(flagTemplates)),
			},
		}.Merge(cmd.MonitoringFlags),
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

	app := cli.Command{
		Name:     "ren",
		Version:  version,
		Commands: commands,
		Suggest:  true,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.Run(ctx, os.Args); err != nil {
		ui.Error(err.Error())
		return 1
	}
	return 0
}
