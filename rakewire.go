package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"rakewire/cmd"
	"rakewire/cmd/remote"
	"strings"
)

// application level variables
var (
	Version   = ""
	BuildTime = ""
	BuildHash = ""
)

func main() {
	app := cli.NewApp()
	app.Name = "Rakewire"
	app.Usage = "Feed Reader"
	app.HideVersion = true
	app.Version = strings.TrimSpace(fmt.Sprintf("%s %s %s", Version, BuildTime, BuildHash))
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "v, verbose",
			EnvVar: "RAKEWIRE_VERBOSE,VERBOSE",
			Usage:  "log more information to console",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start rakewire",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "f, file",
					Value:  "rakewire.db",
					EnvVar: "RAKEWIRE_FILE",
					Usage:  "location of the database file",
				},
				cli.StringFlag{
					Name:   "p, pid",
					Value:  "rakewire.pid",
					EnvVar: "RAKEWIRE_PID",
					Usage:  "location of the pid file",
				},
				cli.IntFlag{
					Name:   "fetch.timeoutsecs",
					Value:  20,
					EnvVar: "RAKEWIRE_FETCH_TIMEOUTSECS",
					Usage:  "fetcher timeout",
				},
				cli.IntFlag{
					Name:   "fetch.workers",
					Value:  10,
					EnvVar: "RAKEWIRE_FETCH_WORKERS",
					Usage:  "fetcher workers",
				},
				cli.StringFlag{
					Name:   "fetch.useragent",
					Value:  strings.TrimSpace("Rakewire " + Version),
					EnvVar: "RAKEWIRE_FETCH_USERAGENT",
					Usage:  "fetcher useragent string",
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
					Usage:       "list the configuration parameters",
					ArgsUsage:   "[prefix]",
					Description: "List the configuration parameters. If an optional prefix is given, restrict listing to parameters names beginning with prefix.",
					Action:      cmd.ConfigList,
				},
				{
					Name:      "get",
					Usage:     "get a configuration parameter",
					ArgsUsage: "<name>",
					Action:    cmd.ConfigGet,
				},
				{
					Name:        "set",
					Usage:       "set a configuration parameter",
					ArgsUsage:   "<name> <value>",
					Description: "Set a configuration parameter. If the value begins with a @ character, it will be read in from @filename.",
					Action:      cmd.ConfigSet,
				},
				{
					Name:      "rm",
					Aliases:   []string{"remove"},
					Usage:     "remove a configuration parameter",
					ArgsUsage: "<name>",
					Action:    cmd.ConfigRemove,
				},
				{
					Name:   "export",
					Usage:  "export the configuration as JSON to stdout",
					Action: cmd.ConfigExport,
				},
				{
					Name:   "import",
					Usage:  "import the configuration as JSON from stdin",
					Action: cmd.ConfigImport,
				},
			},
		},
		{
			Name:    "user",
			Aliases: []string{"u"},
			Usage:   "manage application users",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f, file",
					Value: "rakewire.db",
					Usage: "location of the database file",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:    "ls",
					Aliases: []string{"list"},
					Usage:   "list users",
					Action:  cmd.UserList,
				},
				{
					Name:      "add",
					Usage:     "add user",
					ArgsUsage: "<username> <roles>",
					Action:    cmd.UserAdd,
				},
				{
					Name:      "passwd",
					Usage:     "change user password",
					ArgsUsage: "<username>",
					Action:    cmd.UserPasswordChange,
				},
				{
					Name:      "roles",
					Usage:     "update user roles",
					ArgsUsage: "<username> <roles>",
					Action:    cmd.UserRoles,
				},
				{
					Name:      "rm",
					Aliases:   []string{"remove"},
					Usage:     "remove user",
					ArgsUsage: "<username>",
					Action:    cmd.UserRemove,
				},
			},
		},
		{
			Name:    "remote",
			Aliases: []string{"r"},
			Usage:   "manage remote rakewire instance",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "i, instance",
					EnvVar: "RAKEWIRE_INSTANCE",
					Usage:  "name of remote rakewire instance in the format host:port",
				},
				cli.StringFlag{
					Name:   "u, username",
					EnvVar: "RAKEWIRE_USERNAME",
					Usage:  "name of rakewire user",
				},
				cli.StringFlag{
					Name:   "p, password",
					EnvVar: "RAKEWIRE_PASSWORD",
					Usage:  "password for the rakewire user",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:   "status",
					Usage:  "get instance status",
					Action: remote.Status,
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
