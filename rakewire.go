package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"rakewire/fetch"
	"rakewire/httpd"
	"rakewire/logging"
	"rakewire/model"
	"rakewire/pollfeed"
	"rakewire/reaper"
	"syscall"
	"time"
)

const (
	logName  = "[main]"
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

	var flagVersion = flag.Bool("version", false, "print version and exit")
	var flagFile = flag.String("f", "rakewire.db", "file to open as the rakewire database")
	var flagCheckDatabase = flag.Bool("check", false, "check database integrity and exit")
	var flagPidFile = flag.String("pidfile", "/var/run/rakewire.pid", "PID file")
	flag.Parse()

	log.Printf("Rakewire %s\n", model.Version)
	log.Printf("Commit Time: %s\n", model.CommitTime)
	log.Printf("Commit Hash: %s\n", model.CommitHash)

	if *flagVersion {
		return
	}

	model.AppStart = time.Now()

	var dbFile string
	if filename, err := filepath.Abs(*flagFile); err == nil {
		dbFile = filename
	} else {
		log.Printf("Cannot find database file: %s\n", err.Error())
		return
	}
	log.Printf("Database:   %s\n", dbFile)

	if *flagCheckDatabase {
		if err := model.Instance.Check(dbFile); err != nil {
			log.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
		return
	}

	if db, err := model.Instance.Open(dbFile); db != nil && err == nil {
		database = db
	} else if err != nil {
		log.Println(err.Error())
		model.Instance.Close(db)
		return
	} else if db == nil {
		return
	}

	var cfg *model.Configuration
	if c, err := loadConfiguration(database); err == nil {
		cfg = c
	} else {
		log.Printf("Abort! Cannot load configuration: %s", err.Error())
		model.Instance.Close(database)
		return
	}

	// initialize logging
	loggingConfiguration := &logging.Configuration{
		File:    cfg.GetStr("logging.file"),
		Level:   cfg.GetStr("logging.level", "WARN"),
		NoColor: cfg.GetBool("logging.nocolor"),
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
	monitorShutdown(chErrors, *flagPidFile)

}

func monitorShutdown(chErrors chan error, pidFile string) {

	pidFileWrite(pidFile)

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
	if err := model.Instance.Close(database); err != nil {
		log.Printf("%-7s %-7s Error closing database: %s", logWarn, logName, err.Error())
	}

	pidFileRemove(pidFile)

	log.Printf("%-7s %-7s done", logInfo, logName)

}

func loadConfiguration(db model.Database) (*model.Configuration, error) {
	cfg := model.C.New()
	err := db.Select(func(tx model.Transaction) error {
		cfg = model.C.Get(tx)
		return nil
	})
	return cfg, err
}

func pidFileWrite(pidFile string) {
	if err := ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), os.FileMode(int(0644))); err != nil {
		log.Printf("%-7s %-7s Cannot write pid file: %s", logError, logName, err.Error())
	}
}

func pidFileRemove(pidFile string) {
	if err := os.Remove(pidFile); err != nil {
		log.Printf("%-7s %-7s Cannot remove pid file: %s", logError, logName, err.Error())
	}
}
