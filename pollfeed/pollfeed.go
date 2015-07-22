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
	logger = logging.New("poll")
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
	killed        int32
	running       int32
	runlatch      sync.WaitGroup
	polling       int32
	polllatch     sync.WaitGroup
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
	z.runlatch.Add(1)
	go z.run()
	logger.Println("service started.")
}

// Stop service
func (z *Service) Stop() {
	logger.Println("service stopping...")
	z.kill()
	z.runlatch.Wait()
	logger.Println("service stopped.")
}

func (z *Service) run() {

	logger.Println("run starting...")

	// run once initially
	z.setPolling(true)
	z.polllatch.Add(1)
	go z.poll(nil)

	ticker := time.NewTicker(z.pollFrequency)

run:
	for {
		select {
		case tick := <-ticker.C:
			if !z.isPolling() {
				z.setPolling(true)
				z.polllatch.Add(1)
				go z.poll(&tick)
			} else {
				logger.Println("Polling still in progress, skipping.")
			}
		case <-z.killsignal:
			break run
		}
	}

	ticker.Stop()
	z.polllatch.Wait()

	close(z.Output)

	z.setRunning(false)
	z.runlatch.Done()
	logger.Println("run exited.")

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
	logger.Printf("request feeds: %d", len(requests))

	// send to output
	for i := 0; i < len(requests) && !z.isKilled(); i++ {
		z.Output <- requests[i]
	}

	z.setPolling(false)
	z.polllatch.Done()
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
		atomic.StoreInt32(&z.killed, 0)
	}
}

func (z *Service) kill() {
	z.killsignal <- true
	atomic.StoreInt32(&z.killed, 1)
}

func (z *Service) isKilled() bool {
	return atomic.LoadInt32(&z.killed) != 0
}

func (z *Service) isPolling() bool {
	return atomic.LoadInt32(&z.polling) != 0
}

func (z *Service) setPolling(polling bool) {
	if polling {
		atomic.StoreInt32(&z.polling, 1)
	} else {
		atomic.StoreInt32(&z.polling, 0)
	}
}

func feedsToRequests(feeds *db.Feeds) []*fetch.Request {
	var requests []*fetch.Request
	for _, v := range feeds.Values {
		request := &fetch.Request{
			ID:           v.ID,
			ETag:         v.ETag,
			LastFetch:    v.LastFetch,
			LastModified: v.LastModified,
			URL:          v.URL,
		}
		requests = append(requests, request)
	}
	return requests
}
