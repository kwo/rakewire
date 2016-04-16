package cmd

import (
	"github.com/codegangsta/cli"
	"os"
	"path/filepath"
	"rakewire/logger"
	"rakewire/model"
)

// Check the database
func Check(c *cli.Context) {

	dbFile := c.String("file")
	verbose := c.GlobalBool("verbose")

	if verbose {
		showVersionInformation(c)
	}

	log := logger.New("check")
	log.Silent = !verbose

	if filename, err := filepath.Abs(dbFile); err == nil {
		dbFile = filename
	} else {
		log.Infof("Cannot find database file: %s", err.Error())
		return
	}
	log.Infof("Database: %s", dbFile)

	if err := model.Instance.Check(dbFile); err != nil {
		log.Infof("Error: %s", err.Error())
		os.Exit(1)
	}

}
