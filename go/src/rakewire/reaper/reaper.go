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
	z.killsignal <- true
	z.runlatch.Wait()
	log.Printf("%-7s %-7s service stopped", logInfo, logName)
}

func (z *Service) run() {

	log.Printf("%-7s %-7s run starting...", logDebug, logName)

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
	log.Printf("%-7s %-7s run exited", logDebug, logName)

}

func (z *Service) processResponse(feed *m.Feed) {

	// query previous entries of feed
	var guIDs []string
	for _, entry := range feed.Entries {
		guIDs = append(guIDs, entry.GUID)
	}
	dbEntries, err := z.database.GetFeedEntriesFromIDs(feed.ID, guIDs)
	if err != nil {
		log.Printf("%-7s %-7s Cannot get previous feed entries %s: %s", logWarn, logName, feed.URL, err.Error())
		return
	}

	// setIDs, check dates for new entries
	var mostRecent time.Time
	newEntryCount := 0
	for _, entry := range feed.Entries {

		if dbEntry := dbEntries[entry.GUID]; dbEntry == nil {

			// new entry
			newEntryCount++
			if entry.Created.IsZero() {
				entry.Created = time.Now()
			}
			if entry.Updated.IsZero() {
				entry.Updated = entry.Created
			}

		} else {

			// old entry
			entry.ID = dbEntry.ID

			// set if zero and prevent from creeping forward
			if entry.Created.IsZero() || entry.Created.After(dbEntry.Created) {
				entry.Created = dbEntry.Created
			}

			// set if zero and prevent from creeping forward
			if entry.Updated.IsZero() {
				entry.Updated = dbEntry.Updated
			} else if entry.Updated.After(dbEntry.Updated) && entry.Hash() == dbEntry.Hash() {
				entry.Updated = dbEntry.Updated
			}

		}

		if entry.Updated.After(mostRecent) {
			mostRecent = entry.Updated
		}

	} // loop entries

	feed.Attempt.LastUpdated = mostRecent

	// only bump up LastUpdated if mostRecent is after previous time
	// lastUpdated can move forward if no new entries, if an existing entry has been updated
	if mostRecent.After(feed.LastUpdated) {
		feed.LastUpdated = mostRecent
	}

	if feed.Attempt.Result == m.FetchResultOK {
		if feed.LastUpdated.IsZero() {
			feed.LastUpdated = time.Now() // only if new entries?
		}
	}

	feed.Attempt.EntryCount = len(feed.Entries)
	feed.Attempt.NewEntries = newEntryCount

	switch feed.Status {
	case m.FetchResultOK:
		feed.UpdateFetchTime(feed.LastUpdated)
	case m.FetchResultRedirect:
		feed.AdjustFetchTime(1 * time.Second)
	default: // errors
		feed.UpdateFetchTime(feed.StatusSince)
	}

	// save feed
	err = z.database.SaveFeed(feed)
	if err != nil {
		log.Printf("%-7s %-7s Cannot save feed %s: %s", logWarn, logName, feed.URL, err.Error())
		return
	}

	log.Printf("%-7s %-7s %2s  %3d  %s  %3d/%-3d  %s  %s", logDebug, logName, feed.Status, feed.Attempt.StatusCode, feed.LastUpdated.Local().Format("02.01.06 15:04"), feed.Attempt.NewEntries, feed.Attempt.EntryCount, feed.URL, feed.StatusMessage)

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
