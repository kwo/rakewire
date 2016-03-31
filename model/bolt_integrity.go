package model

import (
	"fmt"
	"github.com/murphysean/cache"
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

	if err := migrateSubscriptionsToFirstOfDuplicateFeeds(newDb); err != nil {
		return err
	}

	if err := removeDuplicateSubscriptions(newDb); err != nil {
		return err
	}

	if err := removeFeedsWithoutSubscription(newDb); err != nil {
		return err
	}

	if err := removeBogusGroups(newDb); err != nil {
		return err
	}

	// TODO: remove duplicate group names

	if err := removeBogusGroupsFromSubscriptions(newDb); err != nil {
		return err
	}

	if err := removeBogusItems(newDb); err != nil {
		return err
	}

	if err := removeBogusEntries(newDb); err != nil {
		return err
	}

	if err := removeBogusTransmissions(newDb); err != nil {
		return err
	}

	// TODO: warn on duplicate user names

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

func groupSubscriptionsByLowercaseFeedURL(tx Transaction, subscriptions Subscriptions) map[string]Subscriptions {

	result := make(map[string]Subscriptions)
	feedsByID := F.GetBySubscriptions(tx, subscriptions).ByID()

	for _, subscription := range subscriptions {
		feed := feedsByID[subscription.FeedID]
		url := strings.ToLower(feed.URL)
		subs := result[url]
		subs = append(subs, subscription)
		result[url] = subs
	}

	return result

}

func makeFeedCache(tx Transaction, max int) lookupFunc {

	c := &cache.PowerCache{}
	c.MaxKeys = max
	c.ValueLoader = func(key string) (interface{}, error) {
		if f := F.Get(tx, key); f != nil {
			return f, nil
		}
		return nil, cache.ErrNotPresent
	}
	c.Initialize()

	return func(id ...string) bool {
		feedID := id[0]
		if f, err := c.Get(feedID); f != nil && err == nil {
			return true
		}
		return false
	}

}

func makeItemCache(tx Transaction, max int) lookupFunc {

	c := &cache.PowerCache{}
	c.MaxKeys = max
	c.ValueLoader = func(key string) (interface{}, error) {
		if i := I.Get(tx, key); i != nil {
			return i, nil
		}
		return nil, cache.ErrNotPresent
	}
	c.Initialize()

	return func(id ...string) bool {
		itemID := id[0]
		if i, err := c.Get(itemID); i != nil && err == nil {
			return true
		}
		return false
	}

}

func makeGroupCache(tx Transaction, max int) lookupFunc {

	c := &cache.PowerCache{}
	c.MaxKeys = max
	c.ValueLoader = func(key string) (interface{}, error) {
		if g := G.Get(tx, key); g != nil {
			return g, nil
		}
		return nil, cache.ErrNotPresent
	}
	c.Initialize()

	return func(id ...string) bool {
		userID := id[0]
		groupID := id[1]
		if g, err := c.Get(groupID); g != nil && err == nil {
			return g.(*Group).UserID == userID
		}
		return false
	}

}

func makeUserCache(tx Transaction, max int) lookupFunc {

	c := &cache.PowerCache{}
	c.MaxKeys = max
	c.ValueLoader = func(key string) (interface{}, error) {
		if u := U.Get(tx, key); u != nil {
			return u, nil
		}
		return nil, cache.ErrNotPresent
	}
	c.Initialize()

	return func(id ...string) bool {
		userID := id[0]
		if u, err := c.Get(userID); u != nil && err == nil {
			return true
		}
		return false
	}

}

func migrateSubscriptionsToFirstOfDuplicateFeeds(db Database) error {

	log.Printf("%-7s %-7s   remove feeds duplicates...", logInfo, logName)

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

		userExists := makeUserCache(tx, 50)
		feedExists := makeFeedCache(tx, 500)
		itemExists := makeItemCache(tx, 50)
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

func removeBogusItems(db Database) error {

	log.Printf("%-7s %-7s   remove bogus items...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		feedExists := makeFeedCache(tx, 500)
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

func removeBogusGroups(db Database) error {

	log.Printf("%-7s %-7s   remove bogus groups...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		userExists := makeUserCache(tx, 50)
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

		groupExists := makeGroupCache(tx, 500)
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

func removeBogusSubscriptions(db Database) error {

	log.Printf("%-7s %-7s   remove bogus subscriptions...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		userExists := makeUserCache(tx, 50)
		feedExists := makeFeedCache(tx, 500)
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

		feedExists := makeFeedCache(tx, 500)
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

func removeDuplicateSubscriptions(db Database) error {

	log.Printf("%-7s %-7s   remove subscription duplicates...", logInfo, logName)

	// loop thru subscriptions per user
	// groups subscriptions by lowercase feed url
	// identify duplicates
	// assign groups to original
	// remove duplicate

	return db.Update(func(tx Transaction) error {
		users := U.Range(tx)
		for _, user := range users {
			subscriptions := S.GetForUser(tx, user.ID)
			subscriptionsByURL := groupSubscriptionsByLowercaseFeedURL(tx, subscriptions)
			for _, subs := range subscriptionsByURL {
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
			} // subscriptionsByURL
		} // users

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
