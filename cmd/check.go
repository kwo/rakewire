package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"path/filepath"
	"rakewire/logger"
	"rakewire/model"
	"time"
)

// Check the database
func Check(c *cli.Context) {

	fmt.Printf("Rakewire %s\n", model.Version)
	fmt.Printf("Build Time: %s\n", model.BuildTime)
	fmt.Printf("Build Hash: %s\n", model.BuildHash)

	model.AppStart = time.Now()

	log := logger.New("check")

	var dbFile string
	if filename, err := filepath.Abs(c.String("f")); err == nil {
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
