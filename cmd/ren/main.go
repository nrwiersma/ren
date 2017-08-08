package main

import (
	"os"

	"gopkg.in/urfave/cli.v1"
)

const (
	FlagPort     = "port"
	FlagLogLevel = "log.level"

	FlagTemplates = "templates"
)

var CommonFlags = []cli.Flag{
	cli.StringFlag{
		Name:   FlagLogLevel,
		Value:  "info",
		Usage:  "Specify the log level. You can use this to enable debug logs by specifying `debug`.",
		EnvVar: "REN_LOG_LEVEL",
	},
}

var Commands = []cli.Command{
	{
		Name:  "server",
		Usage: "Run the ren HTTP server",
		Flags: append([]cli.Flag{
			cli.StringFlag{
				Name:   FlagPort,
				Value:  "80",
				Usage:  "The port to run the server on.",
				EnvVar: "REN_PORT",
			},
			cli.StringFlag{
				Name:   FlagTemplates,
				Value:  "./templates",
				Usage:  "The path to the templates.",
				EnvVar: "REN_TEMPLATES",
			},
		}, CommonFlags...),
		Action: runServer,
	},
}

func main() {
	app := cli.NewApp()
	app.Name = "ren"
	app.Version = Version
	app.Commands = Commands

	app.Run(os.Args)
}
