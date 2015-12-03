package reaper

import (
	"log"
	"rakewire/db"
	m "rakewire/model"
	"sync"
	"sync/atomic"
	"time"
)

const (
	logName  = "[reap]"
	logDebug = "[DEBUG]"
	logInfo  = "[INFO]"
	logWarn  = "[WARN]"
	logError = "[ERROR]"
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
	log.Printf("%-7s %-7s service starting...", logInfo, logName)
	z.setRunning(true)
	z.runlatch.Add(1)
	go z.run()
	log.Printf("%-7s %-7s service started.", logInfo, logName)
}

// Stop service
func (z *Service) Stop() {
	log.Printf("%-7s %-7s service stopping...", logInfo, logName)
	z.killsignal <- true
	z.runlatch.Wait()
	log.Printf("%-7s %-7s service stopped.", logInfo, logName)
}

func (z *Service) run() {

	log.Printf("%-7s %-7s run starting...", logInfo, logName)

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
	log.Printf("%-7s %-7s run exited", logInfo, logName)

}

func (z *Service) processResponse(feed *m.Feed) {

	// query previous entries of feed
	var entryIDs []string
	for _, entry := range feed.Entries {
		entryIDs = append(entryIDs, entry.ID)
	}
	dbEntries, err := z.database.GetFeedEntriesFromIDs(feed.ID, entryIDs)
	if err != nil {
		log.Printf("%-7s %-7s Cannot get previous feed entries %s: %s", logWarn, logName, feed.URL, err.Error())
		return
	}

	// setIDs, check dates for new entries
	var mostRecent time.Time
	newEntryCount := 0
	for i, entry := range feed.Entries {
		if dbEntries[i] == nil {
			// new entry
			newEntryCount++
			entry.GenerateNewID()
			if entry.Created.IsZero() {
				entry.Created = time.Now()
			}
			if entry.Updated.IsZero() {
				entry.Updated = entry.Created
			}
		} else {
			// old entry
			entry.ID = dbEntries[i].ID
		}
		if mostRecent.Before(entry.Updated) {
			mostRecent = entry.Updated
		}
	} // loop entries

	// recalc feed.LastUpdated
	if feed.LastUpdated.Before(mostRecent) {
		feed.LastUpdated = mostRecent
	}

	if feed.Attempt.Result == m.FetchResultOK {
		if feed.LastUpdated.IsZero() {
			feed.LastUpdated = time.Now()
		}
		feed.UpdateFetchTime(feed.LastUpdated)
	}

	feed.Attempt.IsUpdated = (newEntryCount > 0)

	// save feed
	err = z.database.SaveFeed(feed)
	if err != nil {
		log.Printf("%-7s %-7s Cannot save feed %s: %s", logWarn, logName, feed.URL, err.Error())
		return
	}

	log.Printf("%-7s %-7s: %2s  %3d  %5t  %2s %2d %s  %s", logInfo, logName, feed.Status, feed.Attempt.StatusCode, feed.Attempt.IsUpdated, feed.Attempt.UpdateCheck, newEntryCount, feed.URL, feed.StatusMessage)

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
