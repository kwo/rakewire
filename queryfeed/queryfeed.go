package queryfeed

import (
	"rakewire.com/db"
	"rakewire.com/fetch"
	"rakewire.com/logging"
	"sync"
	"sync/atomic"
)

var (
	logger = logging.New("qf")
)

// Configuration for pump service
type Configuration struct {
}

// Service for pumping feeds between fetcher and database
type Service struct {
	Output  chan *fetch.Request
	running int32
	latch   sync.WaitGroup
}

// NewService create a new service
func NewService(cfg *Configuration) *Service {

	return &Service{
		Output: make(chan *fetch.Request),
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
	z.latch.Wait()
	logger.Println("service stopped.")
}

func (z *Service) run() {

	logger.Println("run starting...")

	for z.IsRunning() {

	}

	close(z.Output)

	logger.Println("run exited.")
	z.latch.Done()

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

func databaseFeedsToFetchRequests(dbfeeds *db.Feeds) []*fetch.Request {
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
