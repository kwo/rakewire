package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/kwo/rakewire/model"
)

func showVersionInformation(c *cli.Context) {
	fmt.Printf("Rakewire %s\n", c.App.Version)
}

func openDatabase(dbFile string) (model.Database, error) {

	if filename, err := filepath.Abs(dbFile); err == nil {
		dbFile = filename
	} else {
		return nil, err
	}

	db, err := model.Instance.Open(dbFile)
	if err == nil {
		return db, nil
	}
	model.Instance.Close(db)
	return nil, err

}

func closeDatabase(db model.Database) error {
	return model.Instance.Close(db)
}

func initDb(c *cli.Context) (model.Database, error) {

	dbFile := c.String("file")
	verbose := c.GlobalBool("verbose")

	if verbose {
		showVersionInformation(c)
	}

	db, errDb := openDatabase(dbFile)
	if errDb != nil {
		return nil, errDb
	}

	if verbose {
		fmt.Printf("Database: %s\n", db.Location())
	}

	return db, nil

}
