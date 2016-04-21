package reaper

import (
	"rakewire/logger"
	"rakewire/model"
	"sync"
	"sync/atomic"
	"time"
)

var (
	log = logger.New("reaper")
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
func NewService(database model.Database) *Service {

	return &Service{
		Input:      make(chan *model.Harvest),
		database:   database,
		killsignal: make(chan bool),
	}

}

// Start Service
func (z *Service) Start() error {
	log.Debugf("starting...")
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
	z.killsignal <- true
	z.runlatch.Wait()
	log.Infof("stopped")
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

	log.Debugf("run starting...")

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
	log.Debugf("run exited")

}

func (z *Service) reapHarvest(harvest *model.Harvest) {

	err := z.database.Update(func(tx model.Transaction) error {

		dbItems := z.getDatabaseItems(tx, harvest.Items).GroupByGUID()

		// setIDs, check dates for new items
		var mostRecent time.Time
		newItems := model.Items{}
		for _, item := range harvest.Items {

			if dbItem, ok := dbItems[item.GUID]; !ok {

				// new item
				newItems = append(newItems, item)
				now := time.Now()

				// prevent items marked with a future date
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
		harvest.Transmission.NewItems = len(newItems)

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
			log.Debugf("Cannot save transmission %s: %s", harvest.Transmission.URL, err.Error())
			return err
		}

		// save items
		if err := model.I.SaveAll(tx, harvest.Items); err != nil {
			log.Debugf("Cannot save items %s: %s", harvest.Feed.URL, err.Error())
			return err
		}

		// save feed
		if err := model.F.Save(tx, harvest.Feed); err != nil {
			log.Debugf("Cannot save feed %s: %s", harvest.Feed.URL, err.Error())
			return err
		}

		// save entries
		if err := model.E.AddItems(tx, newItems); err != nil {
			log.Debugf("Cannot save entries %s: %s", harvest.Feed.URL, err.Error())
			return err
		}

		log.Debugf("%2s  %3d  %s  %3d/%-3d  %s  %s", harvest.Feed.Status, harvest.Transmission.StatusCode, harvest.Feed.LastUpdated.Local().Format("02.01.06 15:04"), harvest.Transmission.NewItems, harvest.Transmission.ItemCount, harvest.Feed.URL, harvest.Feed.StatusMessage)

		return nil

	})

	if err != nil {
		log.Debugf("Error processing feed: %s", err.Error())
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
