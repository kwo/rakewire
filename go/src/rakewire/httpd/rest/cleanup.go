package rest

import (
	"log"
	"net/http"
)

func (z *API) cleanup(w http.ResponseWriter, req *http.Request) {

	err := func() error {
		duplicates, err := z.db.FeedDuplicates()
		if err != nil {
			return err
		}

		for url, feedIDs := range duplicates {
			if len(feedIDs) > 1 {

				feed, err := z.db.GetFeedByID(feedIDs[0])
				if err != nil {
					return err
				}

				for _, feedID := range feedIDs[1:] {

					userfeeds, err := z.db.UserFeedGetByFeed(feedID)
					if err != nil {
						return err
					}

					for _, userfeed := range userfeeds {
						if userfeed.FeedID != feed.ID {
							log.Printf("%-7s %-7s Update userfeed %d to %d", logDebug, logName, userfeed.FeedID, feed.ID)
							userfeed.FeedID = feed.ID
							if err := z.db.UserFeedSave(userfeed); err != nil {
								return err
							}
						}
					}

					log.Printf("%-7s %-7s Remove duplicate feed %d: %s", logDebug, logName, feedID, url)
					duplicateFeed, err := z.db.GetFeedByID(feedID)
					if err != nil {
						return err
					}
					if err := z.db.FeedDelete(duplicateFeed); err != nil {
						return err
					}

				}

			}
		}

		return nil
	}()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = func() error {

		feeds, err := z.db.GetFeeds()
		if err != nil {
			return err
		}

		for _, feed := range feeds {

			userfeeds, err := z.db.UserFeedGetByFeed(feed.ID)
			if err != nil {
				return err
			}

			if len(userfeeds) == 0 {
				log.Printf("%-7s %-7s Remove unused feed %d: %s", logDebug, logName, feed.ID, feed.URL)
				if err := z.db.FeedDelete(feed); err != nil {
					return err
				}
			}

		}

		return nil
	}()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
