package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"path/filepath"
	"rakewire/logger"
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

func initConfig(c *cli.Context) (model.Database, *model.Configuration, *logger.Logger, error) {

	dbFile := c.Parent().String("file")
	verbose := c.GlobalBool("verbose")

	if verbose {
		showVersionInformation(c)
	}

	log := logger.New("config")
	log.Silent = !verbose

	db, errDb := openDatabase(dbFile)
	if errDb != nil {
		return nil, nil, nil, errDb
	}
	log.Infof("Database: %s", db.Location())

	cfg, errCfg := loadConfiguration(db)
	if errCfg != nil {
		return nil, nil, nil, errCfg
	}

	return db, cfg, log, nil

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
