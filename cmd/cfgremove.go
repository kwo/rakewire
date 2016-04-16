package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"rakewire/model"
)

// ConfigRemove removes a configuration parameter
func ConfigRemove(c *cli.Context) {

	db, cfg, err := initConfig(c)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer closeDatabase(db)

	if c.NArg() == 1 {

		name := c.Args().First()

		if err := db.Update(func(tx model.Transaction) error {
			delete(cfg.Values, name)
			return model.C.Put(tx, cfg)
		}); err != nil {
			fmt.Printf("Error saving configuration: %s", err.Error())
		}

	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(1)
	}

}
