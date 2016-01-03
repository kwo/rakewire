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

const (
	logName  = "[main]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

var (
	database *bolt.Service
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
		log.Printf("abort! no config file found at %s\n", config.GetConfigFileLocation())
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

	log.Printf("Rakewire %s\n", model.Version)
	log.Printf("Build Time: %s\n", model.BuildTime)
	log.Printf("Build Hash: %s\n", model.BuildHash)
	if *debug {
		log.Printf("%-7s %-7s debug mode enabled", logDebug, logName)
	}
	if *trace {
		log.Printf("%-7s %-7s trace mode enabled", logTrace, logName)
	}

	database = bolt.NewService(&cfg.Database)
	polld = pollfeed.NewService(&cfg.Poll, database)
	reaperd = reaper.NewService(&cfg.Reaper, database)
	fetchd := fetch.NewService(&cfg.Fetcher, polld.Output, reaperd.Input)
	ws = httpd.NewService(&cfg.Httpd, database)

	chErrors := make(chan error, 1)
	for i := 0; i < 5; i++ {
		var err error
		switch i {
		case 0:
			err = database.Start()
		case 1:
			err = polld.Start()
		case 2:
			err = fetchd.Start()
		case 3:
			err = reaperd.Start()
		case 4:
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
	database.Stop()

	log.Printf("%-7s %-7s done", logInfo, logName)

}
