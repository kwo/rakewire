package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
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

type startContext struct {
	database model.Database
	fetchd   *fetch.Service
	polld    *pollfeed.Service
	reaperd  *reaper.Service
	httpd    *httpd.Service
	log      *logger.Logger
	errors   chan error
	pidFile  string
}

// Start the app
func Start(c *cli.Context) {

	fmt.Printf("Rakewire %s\n", model.Version)
	fmt.Printf("Build Time: %s\n", model.BuildTime)
	fmt.Printf("Build Hash: %s\n", model.BuildHash)

	model.AppStart = time.Now()

	ctx := &startContext{
		log:     logger.New("start"),
		pidFile: c.String("pidfile"),
	}

	var dbFile string
	if filename, err := filepath.Abs(c.String("f")); err == nil {
		dbFile = filename
	} else {
		ctx.log.Infof("Cannot find database file: %s", err.Error())
		return
	}
	ctx.log.Infof("Database: %s", dbFile)

	if db, err := model.Instance.Open(dbFile); db != nil && err == nil {
		ctx.database = db
	} else if err != nil {
		ctx.log.Infof(err.Error())
		model.Instance.Close(db)
		return
	} else if db == nil {
		return
	}

	var cfg *model.Configuration
	if c, err := loadConfiguration(ctx.database); err == nil {
		cfg = c
	} else {
		ctx.log.Infof("Abort! Cannot load configuration: %s", err.Error())
		model.Instance.Close(ctx.database)
		return
	}

	// initialize logging - debug statements above this point will never be logged
	// Forbid debugMode in production.
	// If model.Version is not an empty string (stamped via LDFLAGS) then we are in production mode.
	logger.DebugMode = model.Version == "" && c.Bool("debug")

	ctx.polld = pollfeed.NewService(cfg, ctx.database)
	ctx.reaperd = reaper.NewService(cfg, ctx.database)
	ctx.fetchd = fetch.NewService(cfg, ctx.polld.Output, ctx.reaperd.Input)
	ctx.httpd = httpd.NewService(cfg, ctx.database)

	chErrors := make(chan error, 1)
	for i := 0; i < 4; i++ {
		var err error
		switch i {
		case 0:
			err = ctx.polld.Start()
		case 1:
			err = ctx.fetchd.Start()
		case 2:
			err = ctx.reaperd.Start()
		case 3:
			err = ctx.httpd.Start()
		} // select
		if err != nil {
			chErrors <- err
			break
		}
	}

	// we want this to run in the main goroutine
	monitorShutdown(ctx)

}

func monitorShutdown(ctx *startContext) {

	// write pidfile
	if err := ioutil.WriteFile(ctx.pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), os.FileMode(int(0644))); err != nil {
		ctx.log.Infof("Cannot write pid file: %s", err.Error())
	}

	chSignals := make(chan os.Signal, 1)
	signal.Notify(chSignals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-ctx.errors:
		ctx.log.Infof("received error: %s", err.Error())
	case <-chSignals:
		fmt.Println()
		ctx.log.Infof("caught signal")
	}

	ctx.log.Infof("stopping... ")

	// shutdown httpd
	ctx.httpd.Stop()
	ctx.polld.Stop()
	ctx.fetchd.Stop()
	ctx.reaperd.Stop()
	if err := model.Instance.Close(ctx.database); err != nil {
		ctx.log.Infof("Error closing database: %s", err.Error())
	}

	if err := os.Remove(ctx.pidFile); err != nil {
		ctx.log.Infof("Cannot remove pid file: %s", err.Error())
	}

	ctx.log.Infof("done")

}

func loadConfiguration(db model.Database) (*model.Configuration, error) {
	cfg := model.C.New()
	err := db.Select(func(tx model.Transaction) error {
		cfg = model.C.Get(tx)
		return nil
	})
	return cfg, err
}
