package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/kwo/rakewire/model"
	"os"
	"path/filepath"
)

// Init initializes a new rakewire database.
func Init(c *cli.Context) error {

	dbFile := c.String("file")
	if filename, err := filepath.Abs(dbFile); err == nil {
		dbFile = filename
	} else {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	if stat, err := os.Stat(dbFile); stat != nil && err == nil {
		fmt.Printf("File already exists, will not overwrite: %s\n", dbFile)
		os.Exit(1)
	}

	if f, err := os.OpenFile(dbFile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644); err == nil {
		if err := f.Close(); err != nil {
			fmt.Printf("Cannot create database file: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Printf("Cannot create database file: %s\n", err.Error())
		os.Exit(1)
	}

	if db, err := model.Instance.Open(dbFile); err == nil {
		if err := model.Instance.Close(db); err != nil {
			fmt.Printf("Cannot create database file: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Printf("Cannot create database file: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("database initialized at %s\n", dbFile)

	return nil

}
