package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

const (
	FlagPort = "port"

	FlagTemplates = "templates"
)

var Commands = []cli.Command{
	{
		Name:  "server",
		Usage: "Run the bidder HTTP server",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   FlagPort,
				Value:  "80",
				Usage:  "The port to run the server on.",
				EnvVar: "SVG_PORT",
			},
			cli.StringFlag{
				Name:   FlagTemplates,
				Value:  "./templates",
				Usage:  "The path to the templates.",
				EnvVar: "SVG_TEMPLATES",
			},
		},
		Action: runServer,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "svg"
	app.Version = Version
	app.Commands = Commands

	app.Run(os.Args)
}
