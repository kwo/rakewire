package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rakewire/fetch"
	"rakewire/httpd"
	"rakewire/logging"
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

	var flagFile = flag.String("f", "rakewire.db", "file to open as the rakewire database")
	var flagCheckDatabase = flag.Bool("check", false, "check database integrity and exit")
	flag.Parse()

	log.Printf("Rakewire %s\n", model.Version)
	log.Printf("Build Time: %s\n", model.BuildTime)
	log.Printf("Build Hash: %s\n", model.BuildHash)
	log.Printf("Database:   %s\n", *flagFile)

	var err error
	database, err = model.OpenDatabase(*flagFile, *flagCheckDatabase)
	if err != nil {
		log.Println(err.Error())
		model.CloseDatabase(database)
		return
	}
	if database == nil {
		return
	}

	cfg, err := loadConfiguration(database)
	if err != nil {
		log.Printf("Abort! Cannot load configuration: %s", err.Error())
		model.CloseDatabase(database)
		return
	}

	// initialize logging
	loggingConfiguration := &logging.Configuration{
		File:    cfg.Get("logging.file", ""),
		Level:   cfg.Get("logging.level", "WARN"),
		NoColor: cfg.GetBool("logging.nocolor", false),
	}
	loggingConfiguration.Init()

	polld = pollfeed.NewService(cfg, database)
	reaperd = reaper.NewService(cfg, database)
	fetchd := fetch.NewService(cfg, polld.Output, reaperd.Input)
	ws = httpd.NewService(cfg, database)

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

func loadConfiguration(db model.Database) (*model.Configuration, error) {
	cfg := model.NewConfiguration()
	err := db.Select(func(tx model.Transaction) error {
		return cfg.Load(tx)
	})
	return cfg, err
}
