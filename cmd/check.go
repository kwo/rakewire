package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/kwo/rakewire/model"
	"os"
	"path/filepath"
)

// Check the database
func Check(c *cli.Context) error {

	dbFile := c.String("file")
	verbose := c.GlobalBool("verbose")

	if verbose {
		showVersionInformation(c)
	}

	if filename, err := filepath.Abs(dbFile); err == nil {
		dbFile = filename
	} else {
		fmt.Printf("Cannot find database file: %s\n", err.Error())
		return nil
	}
	if verbose {
		fmt.Printf("Database: %s\n", dbFile)
	}

	if err := model.Instance.Check(dbFile); err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	return nil

}
