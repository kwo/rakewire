package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"rakewire/config"
	"rakewire/fetch"
	"rakewire/httpd"
	"rakewire/model"
	"rakewire/pollfeed"
	"rakewire/reaper"
	"syscall"
)

const (
	logName  = "[main]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

var (
	database model.Database
	fetchd   *fetch.Service
	polld    *pollfeed.Service
	reaperd  *reaper.Service
	ws       *httpd.Service
)

func main() {

	cfg := config.GetConfig()
	if cfg == nil {
		log.Printf("abort! no config file found at %s\n", config.GetConfigFileLocation())
		os.Exit(1)
		return
	}

	// initialize logging
	cfg.Logging.Init()

	log.Printf("Rakewire %s\n", model.Version)
	log.Printf("Build Time: %s\n", model.BuildTime)
	log.Printf("Build Hash: %s\n", model.BuildHash)

	var err error
	database, err = model.OpenDatabase(cfg.Database.Location)
	if err != nil {
		log.Printf("Cannot open database: %s", err.Error())
		model.CloseDatabase(database)
		return
	}
	polld = pollfeed.NewService(&cfg.Poll, database)
	reaperd = reaper.NewService(&cfg.Reaper, database)
	fetchd := fetch.NewService(&cfg.Fetcher, polld.Output, reaperd.Input)
	ws = httpd.NewService(&cfg.Httpd, database)

	chErrors := make(chan error, 1)
	for i := 0; i < 4; i++ {
		var err error
		switch i {
		case 0:
			err = polld.Start()
		case 1:
			err = fetchd.Start()
		case 2:
			err = reaperd.Start()
		case 3:
			err = ws.Start()
		} // select
		if err != nil {
			chErrors <- err
			break
		}
	}

	// we want this to run in the main goroutine
	monitorShutdown(chErrors)

}

func monitorShutdown(chErrors chan error) {

	chSignals := make(chan os.Signal, 1)
	signal.Notify(chSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-chErrors:
		log.Printf("%-7s %-7s received error: %s", logError, logName, err.Error())
	case <-chSignals:
		fmt.Println()
		log.Printf("%-7s %-7s caught signal", logInfo, logName)
	}

	log.Printf("%-7s %-7s stopping... ", logInfo, logName)

	// shutdown httpd
	ws.Stop()
	polld.Stop()
	fetchd.Stop()
	reaperd.Stop()
	if err := model.CloseDatabase(database); err != nil {
		log.Printf("%-7s %-7s Error closing database: %s", logWarn, logName, err.Error())
	}

	log.Printf("%-7s %-7s done", logInfo, logName)

}
