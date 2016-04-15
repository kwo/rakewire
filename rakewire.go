package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"rakewire/fetch"
	"rakewire/httpd"
	"rakewire/logger"
	"rakewire/model"
	"rakewire/pollfeed"
	"rakewire/reaper"
	"syscall"
	"time"
)

var (
	database model.Database
	fetchd   *fetch.Service
	polld    *pollfeed.Service
	reaperd  *reaper.Service
	ws       *httpd.Service
	log      = logger.New("main")
)

func main() {

	var flagDebug = flag.Bool("debug", false, "enable debug mode")
	var flagVersion = flag.Bool("version", false, "print version and exit")
	var flagFile = flag.String("f", "rakewire.db", "file to open as the rakewire database")
	var flagCheckDatabase = flag.Bool("check", false, "check database integrity and exit")
	var flagPidFile = flag.String("pidfile", "/var/run/rakewire.pid", "PID file")
	flag.Parse()

	fmt.Printf("Rakewire %s\n", model.Version)
	fmt.Printf("Build Time: %s\n", model.BuildTime)
	fmt.Printf("Build Hash: %s\n", model.BuildHash)

	if *flagVersion {
		return
	}

	model.AppStart = time.Now()

	var dbFile string
	if filename, err := filepath.Abs(*flagFile); err == nil {
		dbFile = filename
	} else {
		log.Infof("Cannot find database file: %s", err.Error())
		return
	}
	log.Infof("Database: %s", dbFile)

	if *flagCheckDatabase {
		if err := model.Instance.Check(dbFile); err != nil {
			log.Infof("Error: %s", err.Error())
			os.Exit(1)
		}
		return
	}

	if db, err := model.Instance.Open(dbFile); db != nil && err == nil {
		database = db
	} else if err != nil {
		log.Infof(err.Error())
		model.Instance.Close(db)
		return
	} else if db == nil {
		return
	}

	var cfg *model.Configuration
	if c, err := loadConfiguration(database); err == nil {
		cfg = c
	} else {
		log.Infof("Abort! Cannot load configuration: %s", err.Error())
		model.Instance.Close(database)
		return
	}

	// initialize logging - debug statements above this point will never be logged
	// Forbid debugMode in production.
	// If model.Version is not an empty string (stamped via LDFLAGS) then we are in production mode.
	logger.DebugMode = model.Version == "" && *flagDebug

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
		log.Infof("received error: %s", err.Error())
	case <-chSignals:
		fmt.Println()
		log.Infof("caught signal")
	}

	log.Infof("stopping... ")

	// shutdown httpd
	ws.Stop()
	polld.Stop()
	fetchd.Stop()
	reaperd.Stop()
	if err := model.Instance.Close(database); err != nil {
		log.Infof("Error closing database: %s", err.Error())
	}

	pidFileRemove(pidFile)

	log.Infof("done")

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
		log.Infof("Cannot write pid file: %s", err.Error())
	}
}

func pidFileRemove(pidFile string) {
	if err := os.Remove(pidFile); err != nil {
		log.Infof("Cannot remove pid file: %s", err.Error())
	}
}
