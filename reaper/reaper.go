package reaper

import (
	"rakewire.com/db"
	"rakewire.com/logging"
	m "rakewire.com/model"
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
	logger.Println("service starting...")
	z.setRunning(true)
	z.runlatch.Add(1)
	go z.run()
	logger.Println("service started.")
}

// Stop service
func (z *Service) Stop() {
	logger.Println("service stopping...")
	z.killsignal <- true
	z.runlatch.Wait()
	logger.Println("service stopped.")
}

func (z *Service) run() {

	logger.Println("run starting...")

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
	logger.Println("run exited.")

}

func (z *Service) processResponse(rsp *m.Feed) {

	//logger.Printf("saving feed: %s %s", rsp.ID, rsp.URL)

	// convert feeds
	feeds := m.NewFeeds()
	feeds.Add(rsp)
	err := z.database.SaveFeeds(feeds)
	if err != nil {
		logger.Printf("Cannot save feed %s: %s", rsp.URL, err.Error())
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
