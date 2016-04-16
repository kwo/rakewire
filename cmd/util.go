package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"path/filepath"
	"rakewire/model"
	"strings"
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

func initConfig(c *cli.Context) (model.Database, *model.Configuration, error) {

	dbFile := c.Parent().String("file")
	verbose := c.GlobalBool("verbose")

	if verbose {
		showVersionInformation(c)
	}

	db, errDb := openDatabase(dbFile)
	if errDb != nil {
		return nil, nil, errDb
	}
	if verbose {
		fmt.Printf("Database: %s\n", db.Location())
	}

	cfg, errCfg := loadConfiguration(db)
	if errCfg != nil {
		return nil, nil, errCfg
	}

	return db, cfg, nil

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
