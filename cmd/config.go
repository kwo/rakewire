package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"rakewire/model"
	"strings"
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

// ConfigExport exports the configuration parameters as JSON
func ConfigExport(c *cli.Context) {

	db, cfg, err := initConfig(c)
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

// ConfigImport imports the configuration parameters as JSON from stdin
func ConfigImport(c *cli.Context) {

	db, cfg, err := initConfig(c)
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

func truncateValue(value string, maxlen int) string {
	if len(value) > maxlen {
		value = value[:maxlen] + "..."
	}
	if pos := strings.IndexAny(value, "\r\n"); pos > -1 {
		value = value[:pos] + " ..."
	}
	return value
}
