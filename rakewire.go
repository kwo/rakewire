package main

import (
	"github.com/codegangsta/cli"
	"os"
	"rakewire/cmd"
	"rakewire/model"
)

// TODO: remove model as import, move version to main package?

func main() {
	app := cli.NewApp()
	app.Name = "Rakewire"
	app.Usage = "Feed Reader"
	app.Version = model.Version
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start rakewire",
			Action: func(c *cli.Context) {
				cmd.Start()
			},
		},
	}
	app.Run(os.Args)
}
