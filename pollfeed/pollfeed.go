package pollfeed

import (
	"log"
	"rakewire/model"
	"sync"
	"sync/atomic"
	"time"
)

const (
	pollInterval        = "poll.interval"
	pollIntervalDefault = time.Second * 5
	pollLimit           = "poll.limit"
	pollLimitDefault    = 10
)

const (
	logName  = "[poll]"
	logTrace = "[TRACE]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
)

// Service for pumping feeds between fetcher and database
type Service struct {
	Output       chan *model.Feed
	database     model.Database
	limit        int
	pollInterval time.Duration
	killsignal   chan bool
	killed       int32
	running      int32
	runlatch     sync.WaitGroup
	polling      int32
	polllatch    sync.WaitGroup
}

// NewService create a new service
func NewService(cfg *model.Configuration, database model.Database) *Service {

	interval, err := time.ParseDuration(cfg.GetStr(pollInterval, pollIntervalDefault.String()))
	if err != nil {
		interval = pollIntervalDefault
		log.Printf("%-7s %-7s Bad or missing poll interval configuration parameter, setting to default of %s.", logWarn, logName, pollIntervalDefault.String())
	}

	return &Service{
		Output:       make(chan *model.Feed),
		limit:        cfg.GetInt(pollLimit, pollLimitDefault),
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

	err := z.database.Select(func(tx model.Transaction) error {

		// get next feeds
		feeds := model.F.GetNext(tx, t)

		// limit runs to X feeds
		if z.limit > 0 && len(feeds) > z.limit {
			feeds = feeds[:z.limit]
		}

		// convert feeds
		if numFeeds := len(feeds); numFeeds > 0 {
			log.Printf("%-7s %-7s polling feeds: %d", logInfo, logName, numFeeds)
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
		log.Printf("%-7s %-7s Error poll feeds: %s", logWarn, logName, err.Error())
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
