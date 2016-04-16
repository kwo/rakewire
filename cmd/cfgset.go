package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"rakewire/model"
	"strings"
)

// ConfigSet sets a configuration parameter
func ConfigSet(c *cli.Context) {

	db, cfg, err := initConfig(c)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer closeDatabase(db)

	if c.NArg() == 2 {
		name := c.Args().First()
		value := c.Args().Get(1)

		if strings.HasPrefix(value, "@") {
			filename := value[1:]
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Printf("Cannot find file: %s\n", filename)
				os.Exit(1)
			}
			cfg.Values[name] = string(data)
		} else {
			cfg.Values[name] = value
		}

		if err := db.Update(func(tx model.Transaction) error {
			return model.C.Put(tx, cfg)
		}); err != nil {
			fmt.Printf("Error saving configuration: %s", err.Error())
		}

	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(1)
	}

}
