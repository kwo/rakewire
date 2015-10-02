package main

import (
	"os"
	"os/signal"
	"rakewire/config"
	"rakewire/db/bolt"
	"rakewire/fetch"
	"rakewire/httpd"
	"rakewire/logging"
	"rakewire/model"
	"rakewire/pollfeed"
	"rakewire/reaper"
	"runtime"
	"syscall"
)

var (
	database *bolt.Database
	fetchd   *fetch.Service
	logger   = logging.New("main")
	polld    *pollfeed.Service
	reaperd  *reaper.Service
	ws       *httpd.Service
)

func main() {

	logger.Printf("Rakewire %s starting with %d CPUs", model.VERSION, runtime.NumCPU())

	cfg := config.GetConfig()
	if cfg == nil {
		logger.Printf("Abort! No config file found at %s\n", config.GetConfigFileLocation())
		os.Exit(1)
		return
	}

	database = &bolt.Database{}
	err := database.Open(&cfg.Database)
	if err != nil {
		logger.Printf("Abort! Cannot access database: %s\n", err.Error())
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
		logger.Printf("Received error: %s", err.Error())
	case <-chSignals:
	}

	logging.Linefeed()
	logger.Println("Stopping... ")

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

	logger.Println("Done")

}
