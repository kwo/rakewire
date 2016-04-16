package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

// ConfigGet gets a configuration parameter
func ConfigGet(c *cli.Context) {

	db, cfg, err := initConfig(c)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer closeDatabase(db)

	if c.NArg() == 1 {
		name := c.Args().First()
		if value, ok := cfg.Values[name]; ok {
			fmt.Printf("%s: %s\n", name, value)
		}
	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(1)
	}

}
