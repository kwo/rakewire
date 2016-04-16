package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

// ConfigExport exports the configuration parameters as JSON
func ConfigExport(c *cli.Context) {

	db, cfg, _, err := initConfig(c)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer closeDatabase(db)

	data, err := json.MarshalIndent(cfg.Values, "", "  ")
	if err != nil {
		fmt.Printf("Cannot export JSON: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(data))

}
