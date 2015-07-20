package queryfeed

import (
	"rakewire.com/db"
	"rakewire.com/fetch"
	"rakewire.com/logging"
	"sync/atomic"
)

var (
	logger = logging.New("queryfeed")
)

// NewService create a new service
func NewService(cfg *Configuration) *Service {

	return &Service{
		countdownLatch: make(chan bool),
	}

}

// Configuration for pump service
type Configuration struct {
}

// Service for pumping feeds between fetcher and database
type Service struct {
	killsignal     int32
	countdownLatch chan bool
}

// Start Service
func (z *Service) Start(chErrors chan error) {

	logger.Println("service starting...")

	z.setRunning(true)
	go z.saveResponses()
	go z.fetchRequests()

	logger.Println("service started.")

}

// Stop service
func (z *Service) Stop() {
	logger.Println("service stopping...")
	z.setRunning(false)
	for i := 0; i < numGoRoutines; i++ {
		<-z.countdownLatch
	}
	logger.Println("service stopped.")
}

func (z *Service) fetchRequests() {

	logger.Println("fetchRequests starting...")

	for z.IsRunning() {

	}

	logger.Println("fetchRequests exited.")
	z.countdownLatch <- true

}

func (z *Service) saveResponses() {

	logger.Println("saveResponses starting...")

	for z.IsRunning() {

	}

	logger.Println("saveResponses exited.")
	z.countdownLatch <- true

}

// IsRunning status of the service
func (z *Service) IsRunning() bool {
	return atomic.LoadInt32(&z.killsignal) != 0
}

func (z *Service) setRunning(running bool) {
	if running {
		atomic.StoreInt32(&z.killsignal, 1)
	} else {
		atomic.StoreInt32(&z.killsignal, 0)
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
