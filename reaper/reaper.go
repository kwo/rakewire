package reaper

import (
	"rakewire.com/db"
	"rakewire.com/fetch"
	"rakewire.com/logging"
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
	Input      chan *fetch.Response
	database   db.Database
	killsignal chan bool
	running    int32
	runlatch   sync.WaitGroup
}

// NewService create a new service
func NewService(cfg *Configuration, database db.Database) *Service {

	return &Service{
		Input:      make(chan *fetch.Response),
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

func (z *Service) processResponse(rsp *fetch.Response) {

	logger.Printf("saving feed: %s", rsp.URL)

	// convert feeds
	feeds := responseToFeeds(rsp)
	err := z.database.SaveFeeds(feeds)
	if err != nil {
		logger.Printf("Error saving feed %s: %s", rsp.URL, err.Error())
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

func responseToFeeds(response *fetch.Response) *db.Feeds {
	var rsps []*fetch.Response
	rsps = append(rsps, response)
	return responsesToFeeds(rsps)
}

func responsesToFeeds(responses []*fetch.Response) *db.Feeds {
	feeds := db.NewFeeds()
	for _, v := range responses {
		feed := &db.Feed{
			ETag:   v.ETag,
			Failed: v.Failed,
			// Flavor: TODO
			// Frequency - intentionally skipping
			// Generator: TODO
			// Hub: TODO
			// Icon: TODO
			ID:           v.ID,
			LastAttempt:  v.AttemptTime,
			LastFetch:    v.FetchTime,
			LastModified: v.LastModified,
			// LastUpdated: TODO
			// Title: TODO
			URL: v.URL,
		}
		feeds.Add(feed)
	}
	return feeds
}
