package pollfeed

import (
	"github.com/kwo/rakewire/logger"
	"github.com/kwo/rakewire/model"
	"sync"
	"sync/atomic"
	"time"
)

var (
	log = logger.New("pollfeed")
)

// Configuration contains all parameters for the PollFeed service
type Configuration struct {
	BatchMax        int
	IntervalSeconds int
}

// Service for pumping feeds between fetcher and database
type Service struct {
	Output       chan *model.Feed
	database     model.Database
	batchMax     int
	pollInterval time.Duration
	killsignal   chan bool
	killed       int32
	running      int32
	runlatch     sync.WaitGroup
	polling      int32
	polllatch    sync.WaitGroup
}

// NewService create a new service
func NewService(cfg *Configuration, database model.Database) *Service {

	return &Service{
		Output:       make(chan *model.Feed),
		batchMax:     cfg.BatchMax,
		database:     database,
		pollInterval: time.Duration(cfg.IntervalSeconds) * time.Second,
		killsignal:   make(chan bool),
	}

}

// Start Service
func (z *Service) Start() error {

	log.Infof("starting...")
	log.Infof("batch max: %d", z.batchMax)
	log.Infof("interval:  %s", z.pollInterval.String())

	z.setRunning(true)
	z.runlatch.Add(1)
	go z.run()
	log.Infof("started")
	return nil
}

// Stop service
func (z *Service) Stop() {

	if !z.IsRunning() {
		log.Debugf("service already stopped, exiting...")
		return
	}

	log.Debugf("stopping...")
	log.Debugf("killing...")
	z.kill()
	log.Debugf("waiting on latch")
	z.runlatch.Wait()
	log.Infof("stopped")
}

func (z *Service) run() {

	log.Debugf("run starting...")

	// run once initially
	z.setPolling(true)
	z.polllatch.Add(1)
	go z.poll(time.Time{})

	ticker := time.NewTicker(z.pollInterval)

run:
	for {
		select {
		case tick := <-ticker.C:
			if !z.isPolling() {
				z.setPolling(true)
				z.polllatch.Add(1)
				go z.poll(tick)
			} else {
				log.Debugf("Polling still in progress, skipping")
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
	log.Debugf("run exited")

}

func (z *Service) poll(t time.Time) {

	err := z.database.Select(func(tx model.Transaction) error {

		// get next feeds
		feeds := model.F.GetNext(tx, t)

		// limit runs to X feeds
		if z.batchMax > 0 && len(feeds) > z.batchMax {
			feeds = feeds[:z.batchMax]
		}

		// convert feeds
		if numFeeds := len(feeds); numFeeds > 0 {
			log.Infof("polling feeds: %d", numFeeds)
		}

		// send to output
		for i := 0; i < len(feeds) && !z.isKilled(); i++ {
			z.Output <- feeds[i]
		}

		z.setPolling(false)
		z.polllatch.Done()

		return nil

	})

	if err != nil {
		log.Infof("Error poll feeds: %s", err.Error())
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
