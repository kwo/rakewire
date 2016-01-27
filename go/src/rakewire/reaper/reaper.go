package reaper

import (
	"log"
	"rakewire/model"
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

// Service for saving fetch responses back to the database
type Service struct {
	Input      chan *model.Feed
	database   model.Database
	killsignal chan bool
	running    int32
	runlatch   sync.WaitGroup
}

// NewService create a new service
func NewService(cfg *model.Configuration, database model.Database) *Service {

	return &Service{
		Input:      make(chan *model.Feed),
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

func (z *Service) processResponse(feed *model.Feed) {

	err := z.database.Update(func(tx model.Transaction) error {

		// query previous items of feed
		var guIDs []string
		for _, item := range feed.Items {
			guIDs = append(guIDs, item.GUID)
		}

		dbItems0, err := model.ItemsByGUIDs(feed.ID, guIDs, tx)
		if err != nil {
			log.Printf("%-7s %-7s Cannot get previous feed items %s: %s", logWarn, logName, feed.URL, err.Error())
			return err
		}
		dbItems := dbItems0.GroupByGUID()

		// setIDs, check dates for new items
		var mostRecent time.Time
		newItemCount := 0
		for _, item := range feed.Items {

			if dbItem, ok := dbItems[item.GUID]; !ok {

				// new item
				newItemCount++
				now := time.Now()

				// prevent items marks with a future date
				if item.Created.IsZero() || item.Created.After(now) {
					item.Created = now
				}
				if item.Updated.IsZero() || item.Updated.After(now) {
					item.Updated = item.Created
				}

			} else {

				// old item
				item.ID = dbItem.ID

				// set if zero and prevent from creeping forward
				if item.Created.IsZero() || item.Created.After(dbItem.Created) {
					item.Created = dbItem.Created
				}

				// set if zero and prevent from creeping forward
				if item.Updated.IsZero() {
					item.Updated = dbItem.Updated
				} else if item.Updated.After(dbItem.Updated) && item.Hash() == dbItem.Hash() {
					item.Updated = dbItem.Updated
				}

			}

			if item.Updated.After(mostRecent) {
				mostRecent = item.Updated
			}

		} // loop items

		feed.Transmission.LastUpdated = mostRecent

		// only bump up LastUpdated if mostRecent is after previous time
		// lastUpdated can move forward if no new items, if an existing item has been updated
		if mostRecent.After(feed.LastUpdated) {
			feed.LastUpdated = mostRecent
		}

		if feed.Transmission.Result == model.FetchResultOK {
			if feed.LastUpdated.IsZero() {
				feed.LastUpdated = time.Now() // only if new items?
			}
		}

		feed.Transmission.ItemCount = len(feed.Items)
		feed.Transmission.NewItems = newItemCount

		switch feed.Status {
		case model.FetchResultOK:
			feed.UpdateFetchTime(feed.LastUpdated)
		case model.FetchResultRedirect:
			feed.AdjustFetchTime(1 * time.Second)
		default: // errors
			feed.UpdateFetchTime(feed.StatusSince)
		}

		// save feed
		_, err = feed.Save(tx)
		if err != nil {
			log.Printf("%-7s %-7s Cannot save feed %s: %s", logWarn, logName, feed.URL, err.Error())
			return err
		}

		log.Printf("%-7s %-7s %2s  %3d  %s  %3d/%-3d  %s  %s", logDebug, logName, feed.Status, feed.Transmission.StatusCode, feed.LastUpdated.Local().Format("02.01.06 15:04"), feed.Transmission.NewItems, feed.Transmission.ItemCount, feed.URL, feed.StatusMessage)

		return nil

	})

	if err != nil {
		log.Printf("%-7s %-7s Error processing feed: %s", logWarn, logName, err.Error())
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
