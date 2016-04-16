package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"rakewire/model"
)

// ConfigImport imports the configuration parameters as JSON from stdin
func ConfigImport(c *cli.Context) {

	db, cfg, _, err := initConfig(c)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer closeDatabase(db)

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Printf("Cannot read JSON: %s\n", err.Error())
		os.Exit(1)
	}

	var values map[string]string
	if err := json.Unmarshal(data, &values); err != nil {
		fmt.Printf("Cannot decode JSON: %s\n", err.Error())
		os.Exit(1)
	}
	cfg.Values = values

	if err := db.Update(func(tx model.Transaction) error {
		return model.C.Put(tx, cfg)
	}); err != nil {
		fmt.Printf("Error saving configuration: %s", err.Error())
	}

}
