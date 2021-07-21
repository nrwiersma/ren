package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/hamba/cmd/v2"
	"github.com/hamba/cmd/v2/term"
	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli/v2"
)

const flagTemplates = "templates"

var version = "¯\\_(ツ)_/¯"

func commands(ui term.Term) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "server",
			Usage: "Run the ren HTTP server",
			Flags: cmd.Flags{
				&cli.StringFlag{
					Name:    flagTemplates,
					Value:   "file:///./templates",
					Usage:   "The URI to the templates. Supported schemes: 'file', 'http', 'https'.",
					EnvVars: []string{"TEMPLATES"},
				},
			}.Merge(cmd.MonitoringFlags, cmd.ServerFlags),
			Action: runServer(ui),
		},
	}
}

func main() {
	ui := newTerm()

	app := &cli.App{
		Name:     "ren",
		Version:  version,
		Commands: commands(ui),
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := app.RunContext(ctx, os.Args); err != nil {
		ui.Error(err.Error())
	}
}
