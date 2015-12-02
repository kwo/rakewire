package reaper

import (
	"log"
	"rakewire/db"
	m "rakewire/model"
	"sync"
	"sync/atomic"
)

const (
	logName  = "[reap]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

// Configuration for reaper service
type Configuration struct {
}

// Service for saving fetch responses back to the database
type Service struct {
	Input      chan *m.Feed
	database   db.Database
	killsignal chan bool
	running    int32
	runlatch   sync.WaitGroup
}

// NewService create a new service
func NewService(cfg *Configuration, database db.Database) *Service {

	return &Service{
		Input:      make(chan *m.Feed),
		database:   database,
		killsignal: make(chan bool),
	}

}

// Start Service
func (z *Service) Start() {
	log.Printf("%-7s %-7s service starting...", logInfo, logName)
	z.setRunning(true)
	z.runlatch.Add(1)
	go z.run()
	log.Printf("%-7s %-7s service started.", logInfo, logName)
}

// Stop service
func (z *Service) Stop() {
	log.Printf("%-7s %-7s service stopping...", logInfo, logName)
	z.killsignal <- true
	z.runlatch.Wait()
	log.Printf("%-7s %-7s service stopped.", logInfo, logName)
}

func (z *Service) run() {

	log.Printf("%-7s %-7s run starting...", logInfo, logName)

run:
	for {
		select {
		case rsp := <-z.Input:
			z.processResponse(rsp)
		case <-z.killsignal:
			break run
		}
	}

	close(z.Input)

	z.setRunning(false)
	z.runlatch.Done()
	log.Printf("%-7s %-7s run exited.", logInfo, logName)

}

func (z *Service) processResponse(feed *m.Feed) {

	// TODO: save entries, noting which ones are new
	// TODO: for new entries with no date, assign time.Now()
	// TODO: recalc feed.Updated
	// save feed

	err := z.database.SaveFeed(feed)
	if err != nil {
		log.Printf("%-7s %-7s Cannot save feed %s: %s", logWarn, logName, feed.URL, err.Error())
	}

	log.Printf("%-7s %-7s: %2s  %3d  %5t  %2s  %s  %s %s", logInfo, logName, feed.Status, feed.Attempt.StatusCode, feed.Attempt.IsUpdated, feed.Attempt.UpdateCheck, feed.URL, feed.StatusMessage, feed.Attempt.Flavor)

}

// IsRunning status of the service
func (z *Service) IsRunning() bool {
	return atomic.LoadInt32(&z.running) != 0
}

func (z *Service) setRunning(running bool) {
	if running {
		atomic.StoreInt32(&z.running, 1)
	} else {
		atomic.StoreInt32(&z.running, 0)
	}
}
