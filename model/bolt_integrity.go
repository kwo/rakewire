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
	backupFilename := z.makeFilenameBackup(filename)
	if err := os.Rename(filename, backupFilename); err != nil {
		return err
	}

	tmpFilename := z.makeFilenameTemp(filename)

	// open both old, new and tmp databases
	// check schema of all databases (in Open)
	var oldDb Database
	if db, err := z.Open(backupFilename); err == nil {
		oldDb = db
	} else {
		return err
	}
	defer z.Close(oldDb)

	var tmpDb Database
	if db, err := z.Open(tmpFilename); err == nil {
		tmpDb = db
	} else {
		return err
	}
	defer z.Close(tmpDb)

	// copy buckets (decoding and re-encoding)
	if err := z.copyBucketsEncodeDecode(oldDb, tmpDb); err != nil {
		return err
	}

	log.Printf("%-7s %-7s validating data...", logInfo, logName)

	// enforce referential integrity
	// enforce unique fields
	//   Feed:URL         - migrate subscriptions, remove
	//   Group:UserIDName - warn only
	//   Item:FeedIDGUID  - warn only
	//   User:Username    - warn only

	if err := z.removeBogusSubscriptions(tmpDb); err != nil {
		return err
	}

	if err := z.removeFeedsWithoutSubscription(tmpDb); err != nil {
		return err
	}

	// enforces uniqueness of lowercase feed.URL
	if err := z.migrateSubscriptionsToFirstOfDuplicateFeeds(tmpDb); err != nil {
		return err
	}

	if err := z.removeSubscriptionsToSameFeed(tmpDb); err != nil {
		return err
	}

	if err := z.removeFeedsWithoutSubscription(tmpDb); err != nil {
		return err
	}

	if err := z.removeBogusGroups(tmpDb); err != nil {
		return err
	}

	if err := z.removeBogusGroupsFromSubscriptions(tmpDb); err != nil {
		return err
	}

	if err := z.removeBogusItems(tmpDb); err != nil {
		return err
	}

	if err := z.removeBogusEntries(tmpDb); err != nil {
		return err
	}

	if err := z.removeBogusTransmissions(tmpDb); err != nil {
		return err
	}

	if err := z.warnUsersWithSameUsername(tmpDb); err != nil {
		return err
	}

	if err := z.warnGroupsWithSameName(tmpDb); err != nil {
		return err
	}

	if err := z.warnItemsWithSameGUID(tmpDb); err != nil {
		return err
	}

	log.Printf("%-7s %-7s validating data done", logInfo, logName)

	// open new database
	var newDb Database
	if db, err := z.Open(filename); err == nil {
		newDb = db
	} else {
		return err
	}
	defer z.Close(newDb)

	// copy buckets
	if err := z.copyBuckets(tmpDb, newDb); err != nil {
		return err
	}

	// rebuild indexes
	if err := z.rebuildIndexes(newDb); err != nil {
		return err
	}

	// close tmpDb, delete
	if err := z.Close(tmpDb); err != nil {
		return err
	}
	if err := os.Remove(tmpFilename); err != nil {
		return err
	}

	return nil

}

func (z *boltInstance) copyBuckets(srcDb, dstDb Database) error {

	log.Printf("%-7s %-7s copying buckets...", logInfo, logName)
	err := srcDb.Select(func(srcTx Transaction) error {
		return dstDb.Update(func(dstTx Transaction) error {
			for entityName := range allEntities {
				log.Printf("%-7s %-7s   %s ...", logInfo, logName, entityName)
				srcBucket := srcTx.Bucket(bucketData, entityName)
				dstBucket := dstTx.Bucket(bucketData, entityName)
				c := srcBucket.Cursor()
				for k, v := c.First(); k != nil; k, v = c.Next() {
					if err := dstBucket.Put(k, v); err != nil {
						return err
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

func (z *boltInstance) copyBucketsEncodeDecode(srcDb, dstDb Database) error {

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

// findFeedsWithSameLowercaseURL returns a map of feeds with same lowercase URL keyed by lowercase URL
// TODO: do not return map but iterate over bucket
// TODO: rename to indexFeedsByLowercaseURL or similar
// TODO: take transaction as parameter
func (z *boltInstance) findFeedsWithSameLowercaseURL(tx Transaction) (map[string]Feeds, error) {

	if err := z.createTempBucket(tx, entityFeed); err != nil {
		return nil, err
	}

	bTmpFeeds := tx.Bucket(bucketTmp, entityFeed)

	feed := &Feed{}
	c := tx.Bucket(bucketData, entityFeed).Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {

		if err := feed.decode(v); err != nil {
			return nil, err
		}

		key := []byte(strings.ToLower(feed.URL))

		feeds := Feeds{}
		if value := bTmpFeeds.Get(key); value != nil {
			if err := feeds.decode(value); err != nil {
				return nil, err
			}
		}
		feeds = append(feeds, feed)

		if value, err := feeds.encode(); err == nil {
			if err := bTmpFeeds.Put(key, value); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}

	} // cursor

	result := make(map[string]Feeds)

	feeds := Feeds{}
	c = bTmpFeeds.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		if err := feeds.decode(v); err != nil {
			return nil, err
		}
		if len(feeds) > 1 {
			lowercaseURL := string(k)
			result[lowercaseURL] = feeds
		}
	}

	return result, nil

}

// findGroupsWithSameName returns a map of groups with same names keyed by UserID|Name
// TODO: redo by creating index by UserID|Name in tmp bucket
// TODO: do not return map but iterate over bucket
// TODO: rename to indexGroupsByUserIDName or similar
// TODO: take transaction as parameter
func (z *boltInstance) findGroupsWithSameName(db Database) (map[string]Groups, error) {

	allGroups := make(map[string]Groups)

	err := db.Select(func(tx Transaction) error {
		group := &Group{}
		c := tx.Bucket(bucketData, entityGroup).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := group.decode(v); err != nil {
				return err
			}
			key := keyEncode(group.UserID, group.Name)
			groups := allGroups[key]
			groups = append(groups, group)
			allGroups[key] = groups
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	result := make(map[string]Groups)
	for key, groups := range allGroups {
		if len(groups) > 1 {
			result[key] = groups
		}
	}
	return result, err

}

// findItemsWithSameGUID returns a map of items with same guids keyed by FeedID|GUID
// TODO: redo by creating index by FeedID|GUID in tmp bucket
// TODO: do not return map but iterate over bucket
// TODO: rename to indexItemsByFeedIDGUID or similar
// TODO: take transaction as parameter
func (z *boltInstance) findItemsWithSameGUID(db Database) (map[string]Items, error) {

	err := db.Update(func(tx Transaction) error {

		// tmp bucket = FeedID|GUID : ItemIDs (space separated)
		if err := z.createTempBucket(tx, "items"); err != nil {
			return err
		}

		bTmp := tx.Bucket(bucketTmp, "items")

		item := &Item{}
		c := tx.Bucket(bucketData, entityItem).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {

			if err := item.decode(v); err != nil {
				return err
			}

			key := []byte(keyEncode(item.FeedID, item.GUID))
			var itemIDs string

			if value := bTmp.Get(key); value != nil {
				itemIDs = string(value)
			}
			if len(itemIDs) == 0 {
				itemIDs = item.GetID()
			} else {
				itemIDs = itemIDs + " " + item.GetID()
			}

			if err := bTmp.Put(key, []byte(itemIDs)); err != nil {
				return err
			}

		}

		return nil

	})

	if err != nil {
		return nil, err
	}

	result := make(map[string]Items)

	err = db.Select(func(tx Transaction) error {

		c := tx.Bucket(bucketTmp, "items").Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			itemIDs := strings.Fields(string(v))
			if len(itemIDs) > 1 {
				key := string(k)
				items := result[key]
				for _, itemID := range itemIDs {
					if item := I.Get(tx, itemID); item != nil {
						items = append(items, item)
					}
				}
				result[key] = items
			}
		}

		return nil

	})

	return result, err

}

func (z *boltInstance) indexSubscriptionsByFeedID(bSubscriptions, bTmp Bucket) error {

	subscription := &Subscription{}
	c := bSubscriptions.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {

		if err := subscription.decode(v); err != nil {
			return err
		}

		key := []byte(subscription.FeedID)

		subscriptions := Subscriptions{}
		if value := bTmp.Get(key); value != nil {
			if err := subscriptions.decode(value); err != nil {
				return err
			}
		}
		subscriptions = append(subscriptions, subscription)

		if value, err := subscriptions.encode(); err == nil {
			if err := bTmp.Put(key, value); err != nil {
				return err
			}
		} else {
			return err
		}

	} // cursor

	return nil

}

// findUsersWithSameUsername returns a map of users with same usernames keyed by username (lowercase)
// TODO: redo by creating index by username|Users in tmp bucket
// TODO: do not return map but iterate over bucket
// TODO: rename to indexUsersByUsernameLowercase or similar
// TODO: take transaction as parameter
func (z *boltInstance) findUsersWithSameUsername(db Database) (map[string]Users, error) {

	allUsers := make(map[string]Users)

	err := db.Select(func(tx Transaction) error {
		user := &User{}
		c := tx.Bucket(bucketData, entityUser).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := user.decode(v); err != nil {
				return err
			}
			username := strings.ToLower(user.Username)
			users := allUsers[username]
			users = append(users, user)
			allUsers[username] = users
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	result := make(map[string]Users)
	for username, users := range allUsers {
		if len(users) > 1 {
			result[username] = users
		}
	}
	return result, err

}

func (z *boltInstance) makeFilenameBackup(location string) string {

	now := time.Now().Truncate(time.Second)
	timestamp := now.Format("20060102150405")

	dir := filepath.Dir(location)
	ext := filepath.Ext(location)
	filename := strings.TrimSuffix(filepath.Base(location), ext)

	return fmt.Sprintf("%s%s%s-%s%s", dir, string(os.PathSeparator), filename, timestamp, ext)

}

func (z *boltInstance) makeFilenameTemp(location string) string {

	ext := filepath.Ext(location)
	filename := strings.TrimSuffix(location, ext)

	return fmt.Sprintf("%s%s", filename, ".tmp")

}

func (z *boltInstance) makeLookupFeed(tx Transaction) lookupFunc {

	return func(id ...string) bool {
		feedID := id[0]
		if i := F.Get(tx, feedID); i != nil {
			return true
		}
		return false
	}

}

func (z *boltInstance) makeLookupGroup(tx Transaction) lookupFunc {

	return func(id ...string) bool {
		groupID := id[0]
		if i := G.Get(tx, groupID); i != nil {
			return true
		}
		return false
	}

}

func (z *boltInstance) makeLookupItem(tx Transaction) lookupFunc {

	return func(id ...string) bool {
		itemID := id[0]
		if i := I.Get(tx, itemID); i != nil {
			return true
		}
		return false
	}

}

func (z *boltInstance) makeLookupSubscription(srcTx, tmpTx Transaction) lookupFunc {

	bSubscriptions := srcTx.Bucket(bucketData, entitySubscription)
	bTmp := tmpTx.Bucket(bucketTmp, entitySubscription)

	if err := z.indexSubscriptionsByFeedID(bSubscriptions, bTmp); err != nil {
		return nil
	}

	return func(id ...string) bool {
		feedIDKey := []byte(id[0])
		subscriptions := Subscriptions{}
		if value := bTmp.Get(feedIDKey); value != nil {
			if err := subscriptions.decode(value); err == nil {
				return len(subscriptions) > 0
			}
		}
		return false
	}

}

func (z *boltInstance) makeLookupUser(tx Transaction) lookupFunc {

	return func(id ...string) bool {
		userID := id[0]
		if i := U.Get(tx, userID); i != nil {
			return true
		}
		return false
	}

}

func (z *boltInstance) migrateSubscriptionsToFirstOfDuplicateFeeds(db Database) error {

	log.Printf("%-7s %-7s   remove feed duplicates...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		// mapped by lowercase url
		feedsByLowercaseURL, errFind := z.findFeedsWithSameLowercaseURL(tx)
		if errFind != nil {
			return errFind
		}

		for url, feeds := range feedsByLowercaseURL {
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

func (z *boltInstance) rebuildIndexes(db Database) error {

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

func (z *boltInstance) removeBogusEntries(db Database) error {

	log.Printf("%-7s %-7s   remove bogus entries...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		userExists := z.makeLookupUser(tx)
		feedExists := z.makeLookupFeed(tx)
		itemExists := z.makeLookupItem(tx)
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

func (z *boltInstance) removeBogusGroups(db Database) error {

	log.Printf("%-7s %-7s   remove bogus groups...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		userExists := z.makeLookupUser(tx)
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

func (z *boltInstance) removeBogusGroupsFromSubscriptions(db Database) error {

	log.Printf("%-7s %-7s   remove bogus groups from subscriptions...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		groupExists := z.makeLookupGroup(tx)
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

func (z *boltInstance) removeBogusItems(db Database) error {

	log.Printf("%-7s %-7s   remove bogus items...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		feedExists := z.makeLookupFeed(tx)
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

func (z *boltInstance) removeBogusSubscriptions(db Database) error {

	log.Printf("%-7s %-7s   remove bogus subscriptions...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		userExists := z.makeLookupUser(tx)
		feedExists := z.makeLookupFeed(tx)
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

func (z *boltInstance) removeBogusTransmissions(db Database) error {

	log.Printf("%-7s %-7s   remove bogus transmissions...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		feedExists := z.makeLookupFeed(tx)
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

func (z *boltInstance) removeFeedsWithoutSubscription(db Database) error {

	log.Printf("%-7s %-7s   remove feeds without subscriptions...", logInfo, logName)

	return db.Update(func(tx Transaction) error {

		if err := z.createTempBucket(tx, entitySubscription); err != nil {
			return err
		}

		subscriptionExists := z.makeLookupSubscription(tx, tx)

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

func (z *boltInstance) removeSubscriptionsToSameFeed(db Database) error {

	log.Printf("%-7s %-7s   remove subscription duplicates...", logInfo, logName)

	// loop thru subscriptions per user
	// groups subscriptions by lowercase feed url
	// identify duplicates
	// assign groups to original
	// remove duplicate

	return db.Update(func(tx Transaction) error {

		user := &User{}
		cUsers := tx.Bucket(bucketData, entityUser).Cursor()

		for userID, userData := cUsers.First(); userID != nil; userID, userData = cUsers.Next() {
			if err := user.decode(userData); err != nil {
				return err
			}
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

func (z *boltInstance) warnItemsWithSameGUID(db Database) error {

	log.Printf("%-7s %-7s   warn item with same guid...", logInfo, logName)

	itemsByFeedIDGUID, err := z.findItemsWithSameGUID(db)
	if err != nil {
		return err
	}

	for key, items := range itemsByFeedIDGUID {
		ids := strings.Split(key, chSep)
		feedID := ids[0]
		guid := ids[1]
		for _, item := range items {
			log.Printf("    item with duplicate guid: %s %s %s", feedID, item.ID, guid)
		}
	}

	return nil

}

func (z *boltInstance) warnGroupsWithSameName(db Database) error {

	log.Printf("%-7s %-7s   warn groups with same name...", logInfo, logName)

	groupsByUserIDName, err := z.findGroupsWithSameName(db)
	if err != nil {
		return err
	}

	for key, groups := range groupsByUserIDName {
		ids := strings.Split(key, chSep)
		userID := ids[0]
		name := ids[1]
		for _, group := range groups {
			log.Printf("    group with same name: %s %s %s", userID, group.ID, name)
		}
	}

	return nil

}

func (z *boltInstance) warnUsersWithSameUsername(db Database) error {

	log.Printf("%-7s %-7s   warn users with same username...", logInfo, logName)

	usersByUsername, err := z.findUsersWithSameUsername(db)
	if err != nil {
		return err
	}

	for username, users := range usersByUsername {
		for _, user := range users {
			log.Printf("    user with same username: %s %s", user.ID, username)
		}
	}

	return nil

}
