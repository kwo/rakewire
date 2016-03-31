package model

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	logName = "[db]"
	logInfo = "[INFO]"
)

type lookupFunc func(id ...string) bool

func (z *boltInstance) Check(filename string) error {

	// backup database (rename to backup name)
	var backupFilename string
	if name, err := backupDatabase(filename); err == nil {
		backupFilename = name
	} else {
		return err
	}

	// open both old and new databases
	// check schema of both databases (in Open)
	var oldDb Database
	if db, err := z.Open(backupFilename); err == nil {
		oldDb = db
	} else {
		return err
	}
	defer z.Close(oldDb)

	var newDb Database
	if db, err := z.Open(filename); err == nil {
		newDb = db
	} else {
		return err
	}
	defer z.Close(newDb)

	// copy buckets
	if err := copyBuckets(oldDb, newDb); err != nil {
		return err
	}

	log.Printf("%-7s %-7s validating data...", logInfo, logName)

	// enforce referential integrity
	// enforce unique fields

	if err := removeBogusSubscriptions(newDb); err != nil {
		return err
	}

	if err := removeFeedsWithoutSubscription(newDb); err != nil {
		return err
	}

	// enforces uniqueness of lowercase feed.URL
	if err := migrateSubscriptionsToFirstOfDuplicateFeeds(newDb); err != nil {
		return err
	}

	if err := removeSubscriptionsToSameFeed(newDb); err != nil {
		return err
	}

	if err := removeFeedsWithoutSubscription(newDb); err != nil {
		return err
	}

	if err := removeBogusGroups(newDb); err != nil {
		return err
	}

	// TODO: test for unique group names by userID, groupName (lowercase)

	if err := removeBogusGroupsFromSubscriptions(newDb); err != nil {
		return err
	}

	if err := removeBogusItems(newDb); err != nil {
		return err
	}

	// TODO: remove non-unique Item GUIDs (FeedID, GUID lowercase)

	if err := removeBogusEntries(newDb); err != nil {
		return err
	}

	if err := removeBogusTransmissions(newDb); err != nil {
		return err
	}

	// TODO: warn on duplicate user names (lowercase)

	log.Printf("%-7s %-7s validating data done", logInfo, logName)

	// rebuild indexes
	if err := rebuildIndexes(newDb); err != nil {
		return err
	}

	return nil

}

func backupDatabase(location string) (string, error) {

	now := time.Now().Truncate(time.Second)
	timestamp := now.Format("20060102150405")

	dir := filepath.Dir(location)
	ext := filepath.Ext(location)
	filename := strings.TrimSuffix(filepath.Base(location), ext)

	newFilename := fmt.Sprintf("%s%s%s-%s%s", dir, string(os.PathSeparator), filename, timestamp, ext)
	err := os.Rename(location, newFilename)

	return newFilename, err

}

func copyBuckets(srcDb, dstDb Database) error {

	log.Printf("%-7s %-7s copying buckets...", logInfo, logName)
	err := srcDb.Select(func(srcTx Transaction) error {
		return dstDb.Update(func(dstTx Transaction) error {
			for entityName := range allEntities {
				log.Printf("%-7s %-7s   %s ...", logInfo, logName, entityName)
				srcBucket := srcTx.Bucket(bucketData, entityName)
				dstBucket := dstTx.Bucket(bucketData, entityName)
				entity := getObject(entityName)
				c := srcBucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					if errDecode := entity.decode(v); errDecode == nil {
						if data, errEncode := entity.encode(); errEncode == nil {
							if err := dstBucket.Put(k, data); err != nil {
								return err
							}
						} else {
							log.Printf("%-7s %-7s     error encoding entity (%s): %s", logInfo, logName, k, errEncode.Error())
						}
					} else {
						log.Printf("%-7s %-7s     error decoding entity (%s): %s", logInfo, logName, k, errDecode.Error())
					}
				} // cursor
			} // entities
			return nil
		}) // update
	}) // select

	if err != nil {
		return err
	}

	log.Printf("%-7s %-7s copying buckets done", logInfo, logName)

	return nil

}

func makeFeedLookup(tx Transaction) lookupFunc {

	feedsByID := F.Range(tx).ByID()

	return func(id ...string) bool {
		feedID := id[0]
		if _, ok := feedsByID[feedID]; ok {
			return true
		}
		return false
	}

}

func makeGroupLookup(tx Transaction) lookupFunc {

	groupsByID := G.Range(tx).ByID()

	return func(id ...string) bool {
		groupID := id[0]
		if _, ok := groupsByID[groupID]; ok {
			return true
		}
		return false
	}

}

func makeItemLookup(tx Transaction) lookupFunc {

	return func(id ...string) bool {
		itemID := id[0]
		if i := I.Get(tx, itemID); i != nil {
			return true
		}
		return false
	}

}

func makeUserLookup(tx Transaction) lookupFunc {

	usersByID := U.Range(tx).ByID()

	return func(id ...string) bool {
		userID := id[0]
		if _, ok := usersByID[userID]; ok {
			return true
		}
		return false
	}

}

func migrateSubscriptionsToFirstOfDuplicateFeeds(db Database) error {

	log.Printf("%-7s %-7s   remove feed duplicates...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		feedsByURL := F.Range(tx).ByURLAll() // mapped by lowercase url

		for url, feeds := range feedsByURL {
			if len(feeds) > 1 {
				log.Printf("     migrating subscriptions of duplicate feed: %s", url)
				// find feed with lowest feedID
				feeds.SortByID()
				originalFeed := feeds[0]
				// loop thru duplicates, getting subscriptions and reassigning to original feed
				for _, feed := range feeds[1:] {
					subscriptions := S.GetForFeed(tx, feed.ID)
					for _, subscription := range subscriptions {
						subscription.FeedID = originalFeed.ID
						if err := S.Save(tx, subscription); err != nil {
							return err
						}
					}
				}
			}
		}

		return nil

	})

}

func rebuildIndexes(db Database) error {

	log.Printf("%-7s %-7s rebuild indexes...", logInfo, logName)

	err := db.Update(func(tx Transaction) error {
		for entityName := range allEntities {
			log.Printf("%-7s %-7s   %s ...", logInfo, logName, entityName)
			bEntity := tx.Bucket(bucketData, entityName)
			bEntityIndex := tx.Bucket(bucketIndex, entityName)
			entity := getObject(entityName)
			c := bEntity.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				if err := entity.decode(v); err == nil {
					// save new indexes
					for indexName, indexKeys := range entity.indexes() {
						bIndex := bEntityIndex.Bucket(indexName)
						key := keyEncode(indexKeys...)
						value := entity.GetID()
						if err := bIndex.Put([]byte(key), []byte(value)); err != nil {
							return err
						}
					} // indexes
				} else {
					// This error should never occur since the entity was freshly encoded when copying the data.
					log.Printf("%-7s %-7s     error decoding entity (%s): %s", logInfo, logName, k, err.Error())
					// Note: just reporting the error means that there is an entity in data without an index
				}
			} // cursor
		} // entities
		return nil
	}) // update

	if err != nil {
		return err
	}

	log.Printf("%-7s %-7s rebuild indexes done", logInfo, logName)

	return nil

}

func removeBogusEntries(db Database) error {

	log.Printf("%-7s %-7s   remove bogus entries...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		userExists := makeUserLookup(tx)
		feedExists := makeFeedLookup(tx)
		itemExists := makeItemLookup(tx)
		badIDs := []string{}

		c := tx.Bucket(bucketData, entityEntry).Cursor()

		entry := &Entry{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			entry.clear()
			if err := entry.decode(v); err == nil {
				if !userExists(entry.UserID) {
					log.Printf("    entry without user: %s (%s)", entry.UserID, entry.GetID())
					badIDs = append(badIDs, entry.GetID())
				} else if !feedExists(entry.FeedID) {
					log.Printf("    entry without feed: %s (%s)", entry.FeedID, entry.GetID())
					badIDs = append(badIDs, entry.GetID())
				} else if !itemExists(entry.ItemID) {
					log.Printf("    entry without item: %s (%s)", entry.ItemID, entry.GetID())
					badIDs = append(badIDs, entry.GetID())
				}
			} else {
				return err
			}
		}

		// remove bad subscriptions
		for _, id := range badIDs {
			if err := E.Delete(tx, id); err != nil {
				return err
			}
		}

		return nil

	})

}

func removeBogusGroups(db Database) error {

	log.Printf("%-7s %-7s   remove bogus groups...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		userExists := makeUserLookup(tx)
		badIDs := []string{}

		c := tx.Bucket(bucketData, entityGroup).Cursor()

		group := &Group{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			group.clear()
			if err := group.decode(v); err == nil {
				if !userExists(group.UserID) {
					log.Printf("    group without user: %s (%s)", group.UserID, group.GetID())
					badIDs = append(badIDs, group.GetID())
				}
			} else {
				return err
			}
		}

		// remove bad items
		for _, id := range badIDs {
			if err := G.Delete(tx, id); err != nil {
				return err
			}
		}

		return nil

	})

}

func removeBogusGroupsFromSubscriptions(db Database) error {

	log.Printf("%-7s %-7s   remove bogus groups from subscriptions...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		groupExists := makeGroupLookup(tx)
		cleanedSubscriptions := Subscriptions{}

		c := tx.Bucket(bucketData, entitySubscription).Cursor()

		subscription := &Subscription{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			subscription.clear()
			if err := subscription.decode(v); err == nil {

				// remove invalid groups
				invalidGroupIDs := []string{}
				for _, groupID := range subscription.GroupIDs {
					if !groupExists(subscription.UserID, groupID) {
						log.Printf("    subscription with invalid group: %s (%s %s)", groupID, subscription.GetID(), subscription.Title)
						invalidGroupIDs = append(invalidGroupIDs, groupID)
					}
				}
				if len(invalidGroupIDs) > 0 {
					for _, groupID := range invalidGroupIDs {
						subscription.RemoveGroup(groupID)
					}
					cleanedSubscriptions = append(cleanedSubscriptions, subscription)
				}

				if len(subscription.GroupIDs) == 0 {
					log.Printf("    subscription without groups: %s %s", subscription.GetID(), subscription.Title)
				}

			} else {
				return err
			}
		}

		// resave cleaned subscriptions
		for _, subscription := range cleanedSubscriptions {
			if err := S.Save(tx, subscription); err != nil {
				return err
			}
		}

		return nil

	})

}

func removeBogusItems(db Database) error {

	log.Printf("%-7s %-7s   remove bogus items...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		feedExists := makeFeedLookup(tx)
		badIDs := []string{}

		c := tx.Bucket(bucketData, entityItem).Cursor()

		item := &Item{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			item.clear()
			if err := item.decode(v); err == nil {
				if !feedExists(item.FeedID) {
					log.Printf("    item without feed: %s (%s %s)", item.FeedID, item.GetID(), item.GUID)
					badIDs = append(badIDs, item.GetID())
				}
			} else {
				return err
			}
		}

		// remove bad items
		for _, id := range badIDs {
			if err := I.Delete(tx, id); err != nil {
				return err
			}
		}

		return nil

	})

}

func removeBogusSubscriptions(db Database) error {

	log.Printf("%-7s %-7s   remove bogus subscriptions...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		userExists := makeUserLookup(tx)
		feedExists := makeFeedLookup(tx)
		badIDs := []string{}

		c := tx.Bucket(bucketData, entitySubscription).Cursor()

		subscription := &Subscription{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			subscription.clear()
			if err := subscription.decode(v); err == nil {
				if !userExists(subscription.UserID) {
					log.Printf("    subscription without user: %s (%s %s)", subscription.UserID, subscription.GetID(), subscription.Title)
					badIDs = append(badIDs, subscription.GetID())
				} else if !feedExists(subscription.FeedID) {
					log.Printf("    subscription without feed: %s (%s %s)", subscription.FeedID, subscription.GetID(), subscription.Title)
					badIDs = append(badIDs, subscription.GetID())
				}
			} else {
				return err
			}
		}

		// remove bad subscriptions
		for _, id := range badIDs {
			if err := S.Delete(tx, id); err != nil {
				return err
			}
		}

		return nil

	})

}

func removeBogusTransmissions(db Database) error {

	log.Printf("%-7s %-7s   remove bogus transmissions...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		feedExists := makeFeedLookup(tx)
		badIDs := []string{}

		c := tx.Bucket(bucketData, entityTransmission).Cursor()

		transmission := &Transmission{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			transmission.clear()
			if err := transmission.decode(v); err == nil {
				if !feedExists(transmission.FeedID) {
					log.Printf("    transmission without feed: %s (%s)", transmission.FeedID, transmission.GetID())
					badIDs = append(badIDs, transmission.GetID())
				}
			} else {
				return err
			}
		}

		// remove bad items
		for _, id := range badIDs {
			if err := T.Delete(tx, id); err != nil {
				return err
			}
		}

		return nil

	})

}

func removeFeedsWithoutSubscription(db Database) error {

	log.Printf("%-7s %-7s   remove feeds without subscriptions...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		subscriptionsByFeedID := S.Range(tx).ByFeedID()
		subscriptionExists := func(id ...string) bool {
			feedID := id[0]
			if subscriptions, ok := subscriptionsByFeedID[feedID]; ok {
				return len(subscriptions) > 0
			}
			return false
		}

		badIDs := []string{}
		c := tx.Bucket(bucketData, entityFeed).Cursor()

		feed := &Feed{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			feed.clear()
			if err := feed.decode(v); err == nil {
				if !subscriptionExists(feed.ID) {
					log.Printf("    feed without subscription: %s %s", feed.ID, feed.Title)
					badIDs = append(badIDs, feed.ID)
				}
			} else {
				return err
			}
		}

		// remove bad feeds
		for _, id := range badIDs {
			if err := F.Delete(tx, id); err != nil {
				return err
			}
		}

		return nil

	})

}

func removeSubscriptionsToSameFeed(db Database) error {

	log.Printf("%-7s %-7s   remove subscription duplicates...", logInfo, logName)

	// loop thru subscriptions per user
	// groups subscriptions by lowercase feed url
	// identify duplicates
	// assign groups to original
	// remove duplicate

	return db.Update(func(tx Transaction) error {
		users := U.Range(tx)
		for _, user := range users {
			subscriptionsByFeedID := S.GetForUser(tx, user.ID).ByFeedID()
			for _, subs := range subscriptionsByFeedID {
				if len(subs) > 1 {
					subs.SortByAddedDate()
					originalSubscription := subs[0]
					for _, s := range subs[1:] {
						for _, groupID := range s.GroupIDs {
							originalSubscription.AddGroup(groupID)
						} // groupIDs
						if err := S.Delete(tx, s.GetID()); err != nil {
							return err
						}
					} // duplicate subscriptions
					if err := S.Save(tx, originalSubscription); err != nil {
						return err
					}
				} // has duplicates
			} // subscriptionsByFeedID
		} // users

		return nil

	})

}
