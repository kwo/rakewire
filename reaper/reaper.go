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
)

// Service for saving fetch responses back to the database
type Service struct {
	Input      chan *model.Harvest
	database   model.Database
	killsignal chan bool
	running    int32
	runlatch   sync.WaitGroup
}

// NewService create a new service
func NewService(cfg *model.Config, database model.Database) *Service {

	return &Service{
		Input:      make(chan *model.Harvest),
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

func (z *Service) run() {

	log.Printf("%-7s %-7s run starting...", logDebug, logName)

run:
	for {
		select {
		case harvest := <-z.Input:
			z.reapHarvest(harvest)
		case <-z.killsignal:
			break run
		}
	}

	close(z.Input)

	z.setRunning(false)
	z.runlatch.Done()
	log.Printf("%-7s %-7s run exited", logDebug, logName)

}

func (z *Service) reapHarvest(harvest *model.Harvest) {

	err := z.database.Update(func(tx model.Transaction) error {

		dbItems := z.getDatabaseItems(tx, harvest.Items).GroupByGUID()

		// setIDs, check dates for new items
		var mostRecent time.Time
		newItemCount := 0
		for _, item := range harvest.Items {

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

		harvest.Transmission.LastUpdated = mostRecent

		// only bump up LastUpdated if mostRecent is after previous time
		// lastUpdated can move forward if no new items, if an existing item has been updated
		if mostRecent.After(harvest.Feed.LastUpdated) {
			harvest.Feed.LastUpdated = mostRecent
		}

		if harvest.Transmission.Result == model.FetchResultOK {
			if harvest.Feed.LastUpdated.IsZero() {
				harvest.Feed.LastUpdated = time.Now() // only if new items?
			}
		}

		harvest.Transmission.ItemCount = len(harvest.Items)
		harvest.Transmission.NewItems = newItemCount

		switch harvest.Feed.Status {
		case model.FetchResultOK:
			harvest.Feed.UpdateFetchTime(harvest.Feed.LastUpdated)
		case model.FetchResultRedirect:
			harvest.Feed.AdjustFetchTime(1 * time.Second)
		default: // errors
			harvest.Feed.UpdateFetchTime(harvest.Feed.StatusSince)
		}

		// save transmission
		if err := model.T.Save(tx, harvest.Transmission); err != nil {
			log.Printf("%-7s %-7s Cannot save transmission %s: %s", logWarn, logName, harvest.Transmission.URL, err.Error())
			return err
		}

		// save items
		if err := model.I.SaveAll(tx, harvest.Items); err != nil {
			log.Printf("%-7s %-7s Cannot save items %s: %s", logWarn, logName, harvest.Feed.URL, err.Error())
			return err
		}

		// save feed
		if err := model.F.Save(tx, harvest.Feed); err != nil {
			log.Printf("%-7s %-7s Cannot save feed %s: %s", logWarn, logName, harvest.Feed.URL, err.Error())
			return err
		}

		log.Printf("%-7s %-7s %2s  %3d  %s  %3d/%-3d  %s  %s", logDebug, logName, harvest.Feed.Status, harvest.Transmission.StatusCode, harvest.Feed.LastUpdated.Local().Format("02.01.06 15:04"), harvest.Transmission.NewItems, harvest.Transmission.ItemCount, harvest.Feed.URL, harvest.Feed.StatusMessage)

		return nil

	})

	if err != nil {
		log.Printf("%-7s %-7s Error processing feed: %s", logWarn, logName, err.Error())
	}

}

func (z *Service) getDatabaseItems(tx model.Transaction, items model.Items) model.Items {

	result := model.Items{}

	for _, item := range items {
		if dbItem := model.I.GetByGUID(tx, item.FeedID, item.GUID); dbItem != nil {
			result = append(result, dbItem)
		}
	}

	return result

}
