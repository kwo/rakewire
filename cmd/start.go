package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/kwo/rakewire/fetch"
	"github.com/kwo/rakewire/httpd"
	"github.com/kwo/rakewire/logger"
	"github.com/kwo/rakewire/model"
	"github.com/kwo/rakewire/pollfeed"
	"github.com/kwo/rakewire/reaper"
	"io/ioutil"
	"os"
	"os/signal"
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
func Start(c *cli.Context) error {

	showVersionInformation(c)

	appStart := time.Now().Unix()
	dbFile := c.String("file")
	pidFile := c.String("pid")
	verbose := c.GlobalBool("verbose")

	ctx := &startContext{
		log:     logger.New("main"),
		pidFile: pidFile,
		errors:  make(chan error, 1),
	}

	if db, err := openDatabase(dbFile); err == nil {
		ctx.database = db
	} else {
		ctx.log.Infof("Error: Cannot open database: %s", err.Error())
		return nil
	}
	ctx.log.Infof("Database: %s", ctx.database.Location())

	// initialize logging - debug statements above this point will never be logged
	// Forbid debugMode in production.
	// If model.Version is not an empty string (stamped via LDFLAGS) then we are in production mode.
	if verbose {
		if len(c.App.Version) == 0 {
			logger.DebugMode = true
		} else {
			ctx.log.Infof("verbose logging not available in production mode")
		}
	}

	pollConfig := &pollfeed.Configuration{
		BatchMax:        c.Int("poll.batchmax"),
		IntervalSeconds: c.Int("poll.intervalsecs"),
	}
	ctx.polld = pollfeed.NewService(pollConfig, ctx.database)
	ctx.reaperd = reaper.NewService(ctx.database)

	fetchConfig := &fetch.Configuration{
		TimeoutSeconds: c.Int("fetch.timeoutsecs"),
		Workers:        c.Int("fetch.workers"),
		UserAgent:      c.App.Version,
	}
	ctx.fetchd = fetch.NewService(fetchConfig, ctx.polld.Output, ctx.reaperd.Input)

	httpdConfig := &httpd.Configuration{
		DebugMode:      len(c.App.Version) == 0,
		ListenHostPort: c.String("bind"),
		PublicHostPort: c.String("host"),
		TLSCertFile:    c.String("tlscert"),
		TLSKeyFile:     c.String("tlskey"),
	}
	ctx.httpd = httpd.NewService(httpdConfig, ctx.database, c.App.Version, appStart)

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
			ctx.errors <- err
			break
		}
	}

	// we want this to run in the main goroutine
	monitorShutdown(ctx)

	return nil

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
