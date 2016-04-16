package main

import (
	"github.com/codegangsta/cli"
	"os"
	"rakewire/cmd"
)

// TODO: remove model as import, move version to main package?

func main() {
	app := cli.NewApp()
	app.Name = "Rakewire"
	app.Usage = "Feed Reader"
	app.HideVersion = true
	app.Version = ""
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "v, verbose",
			EnvVar: "VERBOSE",
			Usage:  "log more information to console",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start rakewire",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f, file",
					Value: "rakewire.db",
					Usage: "location of the database file",
				},
				cli.StringFlag{
					Name:  "p, pid",
					Value: "rakewire.pid",
					Usage: "location of the pid file",
				},
			},
			Action: cmd.Start,
		},
		{
			Name:  "check",
			Usage: "check database",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f, file",
					Value: "rakewire.db",
					Usage: "location of the database file",
				},
			},
			Action: cmd.Check,
		},
		{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "manage application configuration",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f, file",
					Value: "rakewire.db",
					Usage: "location of the database file",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:        "ls",
					Aliases:     []string{"list"},
					Usage:       "list the configuration parameters.",
					ArgsUsage:   "[prefix]",
					Description: "List the configuration parameters. If an optional prefix is given, restrict listing to parameters names beginning with prefix.",
					Action:      cmd.ConfigList,
				},
				{
					Name:      "get",
					Usage:     "get a configuration parameter.",
					ArgsUsage: "<name>",
					Action:    cmd.ConfigGet,
				},
				{
					Name:        "set",
					Usage:       "set a configuration parameter.",
					ArgsUsage:   "<name> <value>",
					Description: "Set a configuration parameter. If the value begins with a @ character, it will be read in from @filename.",
					Action:      cmd.ConfigSet,
				},
				{
					Name:      "rm",
					Aliases:   []string{"remove"},
					Usage:     "remove a configuration parameter.",
					ArgsUsage: "<name>",
					Action:    cmd.ConfigRemove,
				},
				{
					Name:   "export",
					Usage:  "export the configuration as JSON to stdout.",
					Action: cmd.ConfigExport,
				},
				{
					Name:   "import",
					Usage:  "import the configuration as JSON from stdin.",
					Action: cmd.ConfigImport,
				},
			},
		},
		{
			Name:   "version",
			Usage:  "print version and exit",
			Action: cmd.Version,
		},
	}
	app.Run(os.Args)
}
