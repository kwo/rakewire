package pollfeed

import (
	"rakewire.com/db"
	"rakewire.com/fetch"
	"rakewire.com/logging"
	"sync"
	"sync/atomic"
	"time"
)

var (
	logger = logging.New("pollfeed")
)

// Configuration for pump service
type Configuration struct {
	FrequencyMinutes int
}

// Service for pumping feeds between fetcher and database
type Service struct {
	Output        chan *fetch.Request
	database      db.Database
	pollFrequency time.Duration
	killsignal    chan bool
	running       int32
	latch         sync.WaitGroup
}

// NewService create a new service
func NewService(cfg *Configuration, database db.Database) *Service {

	freqMin := cfg.FrequencyMinutes
	if freqMin < 1 {
		freqMin = 5
		logger.Printf("Bad or missing FrequencyMinutes configuration parameter (%d), setting to default of 5 minutes.", cfg.FrequencyMinutes)
	}

	return &Service{
		Output:        make(chan *fetch.Request),
		database:      database,
		pollFrequency: time.Duration(freqMin) * time.Minute,
		killsignal:    make(chan bool),
	}

}

// Start Service
func (z *Service) Start() {
	logger.Println("service starting...")
	z.setRunning(true)
	z.latch.Add(1)
	go z.run()
	logger.Println("service started.")
}

// Stop service
func (z *Service) Stop() {
	logger.Println("service stopping...")
	z.setRunning(false)
	z.killsignal <- true
	z.latch.Wait()
	logger.Println("service stopped.")
}

func (z *Service) run() {

	logger.Println("run starting...")

	ticker := time.NewTicker(z.pollFrequency)

	select {
	case tick := <-ticker.C:
		z.poll(&tick)
	case <-z.killsignal:
		ticker.Stop()
		break
	}

	close(z.Output)

	logger.Println("run exited.")
	z.latch.Done()

}

func (z *Service) poll(t *time.Time) {

	logger.Println("polling...")

	// get next feeds
	feeds, err := z.database.GetFetchFeeds(t)
	if err != nil {
		logger.Printf("Cannot poll feeds: %s", err.Error())
		return
	}

	// convert feeds
	requests := feedsToRequests(feeds)

	logger.Println("sending feeds to output channel")

	// send to output
	for _, req := range requests {
		z.Output <- req
	}

	logger.Println("polling exited.")

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

func feedsToRequests(dbfeeds *db.Feeds) []*fetch.Request {
	var feeds []*fetch.Request
	for _, v := range dbfeeds.Values {
		feed := &fetch.Request{
			ID:           v.ID,
			ETag:         v.ETag,
			LastModified: v.LastModified,
			URL:          v.URL,
		}
		feeds = append(feeds, feed)
	}
	return feeds
}
