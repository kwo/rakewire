package reaper

import (
	"rakewire/db"
	"rakewire/logging"
	m "rakewire/model"
	"sync"
	"sync/atomic"
)

var (
	logger = logging.New("reap")
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
	logger.Info("service starting...")
	z.setRunning(true)
	z.runlatch.Add(1)
	go z.run()
	logger.Info("service started.")
}

// Stop service
func (z *Service) Stop() {
	logger.Info("service stopping...")
	z.killsignal <- true
	z.runlatch.Wait()
	logger.Info("service stopped.")
}

func (z *Service) run() {

	logger.Info("run starting...")

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
	logger.Info("run exited.")

}

func (z *Service) processResponse(rsp *m.Feed) {

	//logger.Debugf("saving feed: %s %s", rsp.ID, rsp.URL)

	// convert feeds
	err := z.database.SaveFeed(rsp)
	if err != nil {
		logger.Warnf("Cannot save feed %s: %s", rsp.URL, err.Error())
	}

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
