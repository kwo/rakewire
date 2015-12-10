package pollfeed

import (
	"log"
	"rakewire/db"
	m "rakewire/model"
	"sync"
	"sync/atomic"
	"time"
)

const (
	minPollInterval = time.Minute * 1
)

const (
	logName  = "[poll]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
)

// Configuration for pump service
type Configuration struct {
	Interval string
}

// Service for pumping feeds between fetcher and database
type Service struct {
	Output       chan *m.Feed
	database     db.Database
	pollInterval time.Duration
	killsignal   chan bool
	killed       int32
	running      int32
	runlatch     sync.WaitGroup
	polling      int32
	polllatch    sync.WaitGroup
}

// NewService create a new service
func NewService(cfg *Configuration, database db.Database) *Service {

	interval, err := time.ParseDuration(cfg.Interval)
	if err != nil || interval < minPollInterval {
		interval = minPollInterval
		log.Printf("%-7s %-7s Bad or missing interval configuration parameter (%s), setting to default of %s.", logWarn, logName, cfg.Interval, minPollInterval.String())
	}

	return &Service{
		Output:       make(chan *m.Feed),
		database:     database,
		pollInterval: interval,
		killsignal:   make(chan bool),
	}

}

// Start Service
func (z *Service) Start() error {
	log.Printf("%-7s %-7s service starting...", logDebug, logName)
	z.setRunning(true)
	z.runlatch.Add(1)
	go z.run()
	log.Printf("%-7s %-7s service started", logInfo, logName)
	return nil
}

// Stop service
func (z *Service) Stop() {

	if !z.IsRunning() {
		log.Printf("%-7s %-7s service already stopped, exiting...", logWarn, logName)
		return
	}

	log.Printf("%-7s %-7s service stopping...", logDebug, logName)
	log.Printf("%-7s %-7s killing...", logTrace, logName)
	z.kill()
	log.Printf("%-7s %-7s waiting on latch", logTrace, logName)
	z.runlatch.Wait()
	log.Printf("%-7s %-7s service stopped", logInfo, logName)
}

func (z *Service) run() {

	log.Printf("%-7s %-7s run starting...", logDebug, logName)

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
				log.Printf("%-7s %-7s Polling still in progress, skipping", logDebug, logName)
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
	log.Printf("%-7s %-7s run exited", logDebug, logName)

}

func (z *Service) poll(t time.Time) {

	log.Printf("%-7s %-7s polling...", logDebug, logName)

	// get next feeds
	feeds, err := z.database.GetFetchFeeds(t)
	if err != nil {
		log.Printf("%-7s %-7s Cannot poll feeds: %s", logWarn, logName, err.Error())
		z.setPolling(false)
		z.polllatch.Done()
		return
	}

	// convert feeds
	log.Printf("%-7s %-7s polling feeds: %d", logInfo, logName, len(feeds))

	// send to output
	for i := 0; i < len(feeds) && !z.isKilled(); i++ {
		z.Output <- feeds[i]
	}

	z.setPolling(false)
	z.polllatch.Done()

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
