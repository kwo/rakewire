package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"rakewire/config"
	"rakewire/db/bolt"
	"rakewire/fetch"
	"rakewire/httpd"
	"rakewire/model"
	"rakewire/pollfeed"
	"rakewire/reaper"
	"runtime"
	"syscall"
)

const (
	logName  = "main "
	logDebug = "DEBUG"
	logInfo  = "INFO "
	logWarn  = "WARN "
	logError = "ERROR"
)

var (
	database *bolt.Database
	fetchd   *fetch.Service
	polld    *pollfeed.Service
	reaperd  *reaper.Service
	ws       *httpd.Service
)

func main() {

	fmt.Printf("Rakewire %s starting with %d CPUs\n", model.VERSION, runtime.NumCPU())

	cfg := config.GetConfig()
	if cfg == nil {
		fmt.Printf("Abort! No config file found at %s\n", config.GetConfigFileLocation())
		os.Exit(1)
		return
	}

	database = &bolt.Database{}
	err := database.Open(&cfg.Database)
	if err != nil {
		log.Printf("%s %s Abort! Cannot access database: %s\n", logError, logName, err.Error())
		os.Exit(1)
		return
	}

	polld = pollfeed.NewService(&cfg.Poll, database)
	reaperd = reaper.NewService(&cfg.Reaper, database)
	fetchd := fetch.NewService(&cfg.Fetcher, polld.Output, reaperd.Input)

	polld.Start()
	fetchd.Start()
	reaperd.Start()

	chErrors := make(chan error)

	ws = &httpd.Service{
		Database: database,
	}
	go ws.Start(&cfg.Httpd, chErrors)

	monitorShutdown(chErrors)

}

func monitorShutdown(chErrors chan error) {

	chSignals := make(chan os.Signal, 1)
	signal.Notify(chSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-chErrors:
		log.Printf("%s %s Received error: %s", logError, logName, err.Error())
	case <-chSignals:
		fmt.Println()
		log.Printf("%s %s caught signal", logInfo, logName)
	}

	log.Printf("%s %s Stopping... ", logInfo, logName)

	// shutdown httpd
	ws.Stop()
	ws = nil

	polld.Stop()
	polld = nil

	fetchd.Stop()
	fetchd = nil

	reaperd.Stop()
	reaperd = nil

	// close database
	database.Close()
	database = nil

	log.Printf("%s %s Done", logInfo, logName)

}
