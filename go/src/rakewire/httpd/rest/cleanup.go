package rest

import (
	"log"
	"net/http"
	"rakewire/model"
)

func (z *API) cleanup(w http.ResponseWriter, req *http.Request) {

	log.Printf("%-7s %-7s starting cleanup ...", logDebug, logName)

	err := z.db.Update(func(tx model.Transaction) error {

		log.Printf("%-7s %-7s duplicate adjustment...", logDebug, logName)

		duplicates, err := model.FeedDuplicates(tx)
		if err != nil {
			return err
		}

		for _, feedIDs := range duplicates {

			feed, err := model.FeedByID(feedIDs[0], tx)
			if err != nil {
				return err
			}

			if len(feedIDs) > 1 {
				for _, feedID := range feedIDs[1:] {

					userfeeds, err := model.UserFeedsByFeed(feedID, tx)
					if err != nil {
						return err
					}

					for _, userfeed := range userfeeds {
						if userfeed.FeedID != feed.ID {
							log.Printf("%-7s %-7s Repointing subscriptions of duplicate %d to %d", logDebug, logName, userfeed.FeedID, feed.ID)
							userfeed.FeedID = feed.ID
							if err := userfeed.Save(tx); err != nil {
								return err
							}
						}
					}

				} // loop thru duplicate IDs

			} // duplicates found

			// save this feed to be certain indexes are pointing to it
			if _, err := feed.Save(tx); err != nil {
				return err
			}

		} // duplicates

		log.Printf("%-7s %-7s duplicate adjustment finished", logDebug, logName)

		log.Printf("%-7s %-7s unused removal...", logDebug, logName)

		feeds, err := model.FeedDuplicates(tx)
		if err != nil {
			return err
		}

		for _, feedIDs := range feeds {

			//log.Printf("%-7s %-7s feed %d: %s", logDebug, logName, feedIDs[0], url)

			feed, err := model.FeedByID(feedIDs[0], tx)
			if err != nil {
				return err
			}

			userfeeds, err := model.UserFeedsByFeed(feed.ID, tx)
			if err != nil {
				return err
			}

			if len(userfeeds) == 0 {
				log.Printf("%-7s %-7s Remove unused feed %d: %s", logDebug, logName, feed.ID, feed.URL)
				if err := feed.Delete(tx); err != nil {
					return err
				}
			}

		}

		log.Printf("%-7s %-7s unused removal complete", logDebug, logName)

		return nil

	}) // transaction

	log.Printf("%-7s %-7s cleanup complete: %v", logDebug, logName, err)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
