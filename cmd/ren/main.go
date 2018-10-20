package main

import (
	"os"

	"github.com/msales/pkg/clix"
	"github.com/urfave/cli"
)

import _ "github.com/joho/godotenv/autoload"

// Flag constants declared for CLI use.
const (
	FlagTemplates = "templates"
)

var version = "¯\\_(ツ)_/¯"

var commands = []cli.Command{
	{
		Name:  "server",
		Usage: "Run the ren HTTP server",
		Flags: clix.Flags{
			cli.StringFlag{
				Name:   FlagTemplates,
				Value:  "./templates",
				Usage:  "The path to the templates.",
				EnvVar: "TEMPLATES",
			},
		}.Merge(clix.CommonFlags, clix.ServerFlags),
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
