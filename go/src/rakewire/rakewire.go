package main

import (
	"flag"
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
	"syscall"
)

// TODO: replace yaml config with key value package - give each package only subset of all KVs.

const (
	logName  = "[main]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

var (
	database *bolt.Database
	fetchd   *fetch.Service
	polld    *pollfeed.Service
	reaperd  *reaper.Service
	ws       *httpd.Service
)

func main() {

	var debug = flag.Bool("debug", false, "run in debug mode")
	var trace = flag.Bool("trace", false, "run in trace mode, implies debug")
	flag.Parse()

	cfg := config.GetConfig()
	if cfg == nil {
		log.Printf("Abort! No config file found at %s\n", config.GetConfigFileLocation())
		os.Exit(1)
		return
	}

	if *debug {
		cfg.Logging.Level = "DEBUG"
		cfg.Httpd.UseLocal = true
	}
	if *trace {
		cfg.Logging.Level = "TRACE"
		cfg.Httpd.UseLocal = true
	}

	// initialize logging
	cfg.Logging.Init()

	log.Printf("%-7s %-7s Rakewire %s\n", logInfo, logName, model.Version)
	log.Printf("%-7s %-7s Build Time: %s\n", logInfo, logName, model.BuildTime)
	log.Printf("%-7s %-7s Build Hash: %s\n", logInfo, logName, model.BuildHash)
	if *debug {
		log.Printf("%-7s %-7s Debug mode enabled.", logDebug, logName)
	}
	if *trace {
		log.Printf("%-7s %-7s Trace mode enabled.", logTrace, logName)
	}

	database = &bolt.Database{}
	err := database.Open(&cfg.Database)
	if err != nil {
		log.Printf("%-7s %-7s Abort! Cannot access database: %s", logError, logName, err.Error())
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

	ws = &httpd.Service{}
	go ws.Start(&cfg.Httpd, database, chErrors)

	monitorShutdown(chErrors)

}

func monitorShutdown(chErrors chan error) {

	chSignals := make(chan os.Signal, 1)
	signal.Notify(chSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-chErrors:
		log.Printf("%-7s %-7s Received error: %s", logError, logName, err.Error())
	case <-chSignals:
		fmt.Println()
		log.Printf("%-7s %-7s caught signal", logInfo, logName)
	}

	log.Printf("%-7s %-7s Stopping... ", logInfo, logName)

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

	log.Printf("%-7s %-7s Done", logInfo, logName)

}
