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
	app.Version = ""
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start rakewire",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f",
					Value: "rakewire.db",
					Usage: "location of the database",
				},
				cli.StringFlag{
					Name:  "pidfile",
					Value: "/var/run/rakewire.pid",
					Usage: "location of the pidfile",
				},
				cli.BoolFlag{
					Name:  "debug",
					Usage: "show debug log statements",
				},
			},
			Action: cmd.Start,
		},
		{
			Name:  "check",
			Usage: "check database",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "f",
					Value: "rakewire.db",
					Usage: "location of the database",
				},
			},
			Action: cmd.Check,
		},
	}
	app.Run(os.Args)
}
