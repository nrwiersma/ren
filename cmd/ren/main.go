package main

import (
	"os"

	"github.com/hamba/cmd"
	"gopkg.in/urfave/cli.v1"
)

import _ "github.com/joho/godotenv/autoload"

const (
	flagTemplates = "templates"
)

var version = "¯\\_(ツ)_/¯"

var commands = []cli.Command{
	{
		Name:  "server",
		Usage: "Run the ren HTTP server",
		Flags: cmd.Flags{
			cli.StringFlag{
				Name:   flagTemplates,
				Value:  "file:///./templates",
				Usage:  "The URI to the templates. Supported schemes: 'file', 'http', 'https'.",
				EnvVar: "TEMPLATES",
			},
		}.Merge(cmd.CommonFlags, cmd.ServerFlags),
		Action: runServer,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "ren"
	app.Version = version
	app.Commands = commands

	app.Run(os.Args)
}
