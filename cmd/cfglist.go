package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"strings"
)

// ConfigList lists the configuration
func ConfigList(c *cli.Context) {

	db, cfg, err := initConfig(c)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer closeDatabase(db)

	if c.NArg() == 0 {
		// print all, truncate long values
		for key, value := range cfg.Values {
			fmt.Printf("%s: %s\n", key, truncateValue(value, 80))
		}
	} else {
		prefix := c.Args().First()
		// print all with prefix, truncate long values
		for key, value := range cfg.Values {
			if strings.HasPrefix(key, prefix) {
				fmt.Printf("%s: %s\n", key, truncateValue(value, 80))
			}
		}
	}

}
