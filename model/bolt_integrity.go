package model

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	defer os.Remove(tmpFilename)

	// copy buckets (decoding and re-encoding)
	if err := z.copyBucketsEncodeDecode(oldDb, tmpDb); err != nil {
		return err
	}

	z.log.Infof("validating data...")

	if err := z.createTempTopLevelBucket(tmpDb); err != nil {
		return err
	}

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

	if err := z.removeTempTopLevelBucket(tmpDb); err != nil {
		return err
	}

	z.log.Infof("validating data done")

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

	return nil

}

func (z *boltInstance) copyBuckets(srcDb, dstDb Database) error {

	z.log.Infof("copying buckets...")
	err := srcDb.Select(func(srcTx Transaction) error {
		return dstDb.Update(func(dstTx Transaction) error {
			for entityName := range allEntities {
				z.log.Infof("  %s ...", entityName)
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

	z.log.Infof("copying buckets done")

	return nil

}

func (z *boltInstance) copyBucketsEncodeDecode(srcDb, dstDb Database) error {

	z.log.Infof("copying buckets...")
	err := srcDb.Select(func(srcTx Transaction) error {
		return dstDb.Update(func(dstTx Transaction) error {
			for entityName := range allEntities {
				z.log.Infof("  %s ...", entityName)
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
							z.log.Infof("    error encoding entity (%s): %s", k, errEncode.Error())
						}
					} else {
						z.log.Infof("    error decoding entity (%s): %s", k, errDecode.Error())
					}
				} // cursor
			} // entities
			return nil
		}) // update
	}) // select

	if err != nil {
		return err
	}

	z.log.Infof("copying buckets done")

	return nil

}

func (z *boltInstance) indexFeedsByURLLowercase(bFeeds, bTmp Bucket) error {

	feed := &Feed{}
	c := bFeeds.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {

		if err := feed.decode(v); err != nil {
			return err
		}

		key := []byte(strings.ToLower(feed.URL))

		feeds := Feeds{}
		if value := bTmp.Get(key); value != nil {
			if err := feeds.decode(value); err != nil {
				return err
			}
		}
		feeds = append(feeds, feed)

		if value, err := feeds.encode(); err == nil {
			if err := bTmp.Put(key, value); err != nil {
				return err
			}
		} else {
			return err
		}

	} // cursor

	return nil

}

func (z *boltInstance) indexGroupsByUserIDName(bGroups, bTmp Bucket) error {

	group := &Group{}
	c := bGroups.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {

		if err := group.decode(v); err != nil {
			return err
		}

		key := []byte(keyEncode(group.UserID, group.Name))

		groups := Groups{}
		if value := bTmp.Get(key); value != nil {
			if err := groups.decode(value); err != nil {
				return err
			}
		}
		groups = append(groups, group)

		if value, err := groups.encode(); err == nil {
			if err := bTmp.Put(key, value); err != nil {
				return err
			}
		} else {
			return err
		}

	} // cursor

	return nil

}

func (z *boltInstance) indexItemsByFeedIDGUID(bItems, bTmp Bucket) error {

	item := &Item{}
	c := bItems.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {

		if err := item.decode(v); err != nil {
			return err
		}

		key := []byte(keyEncode(item.FeedID, item.GUID))

		items := Items{}
		if value := bTmp.Get(key); value != nil {
			if err := items.decode(value); err != nil {
				return err
			}
		}
		items = append(items, item)

		if value, err := items.encode(); err == nil {
			if err := bTmp.Put(key, value); err != nil {
				return err
			}
		} else {
			return err
		}

	} // cursor

	return nil

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

func (z *boltInstance) indexUsersByUsernameLowecase(bUsers, bTmp Bucket) error {

	user := &User{}
	c := bUsers.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {

		if err := user.decode(v); err != nil {
			return err
		}

		key := []byte(strings.ToLower(user.Username))

		users := Users{}
		if value := bTmp.Get(key); value != nil {
			if err := users.decode(value); err != nil {
				return err
			}
		}
		users = append(users, user)

		if value, err := users.encode(); err == nil {
			if err := bTmp.Put(key, value); err != nil {
				return err
			}
		} else {
			return err
		}

	} // cursor

	return nil

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

func (z *boltInstance) makeLookupSubscriptionByFeedID(tx Transaction) lookupFunc {

	bSubscriptions := tx.Bucket(bucketData, entitySubscription)
	bTmp := z.createTempBucket(tx, entitySubscription)

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

	z.log.Infof("  remove feed duplicates...")

	return db.Update(func(tx Transaction) error {

		bFeeds := tx.Bucket(bucketData, entityFeed)
		bTmp := z.createTempBucket(tx, entityFeed)

		if err := z.indexFeedsByURLLowercase(bFeeds, bTmp); err != nil {
			return err
		}

		c := bTmp.Cursor()
		feeds := Feeds{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := feeds.decode(v); err != nil {
				return err
			}
			if len(feeds) > 1 {
				url := string(k)
				z.log.Infof("migrating subscriptions of duplicate feed: %s", url)
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

	z.log.Infof("rebuild indexes...")

	err := db.Update(func(tx Transaction) error {
		for entityName := range allEntities {
			z.log.Infof("  %s ...", entityName)
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
					z.log.Infof("    error decoding entity (%s): %s", k, err.Error())
					// Note: just reporting the error means that there is an entity in data without an index
				}
			} // cursor
		} // entities
		return nil
	}) // update

	if err != nil {
		return err
	}

	z.log.Infof("rebuild indexes done")

	return nil

}

func (z *boltInstance) removeBogusEntries(db Database) error {

	z.log.Infof("  remove bogus entries...")

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
					z.log.Infof("entry without user: %s (%s)", entry.UserID, entry.GetID())
					badIDs = append(badIDs, entry.GetID())
				} else if !feedExists(entry.FeedID) {
					z.log.Infof("entry without feed: %s (%s)", entry.FeedID, entry.GetID())
					badIDs = append(badIDs, entry.GetID())
				} else if !itemExists(entry.ItemID) {
					z.log.Infof("entry without item: %s (%s)", entry.ItemID, entry.GetID())
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

	z.log.Infof("  remove bogus groups...")

	return db.Update(func(tx Transaction) error {

		userExists := z.makeLookupUser(tx)
		badIDs := []string{}

		c := tx.Bucket(bucketData, entityGroup).Cursor()

		group := &Group{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			group.clear()
			if err := group.decode(v); err == nil {
				if !userExists(group.UserID) {
					z.log.Infof("group without user: %s (%s)", group.UserID, group.GetID())
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

	z.log.Infof("  remove bogus groups from subscriptions...")

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
						z.log.Infof("subscription with invalid group: %s (%s %s)", groupID, subscription.GetID(), subscription.Title)
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
					z.log.Infof("subscription without groups: %s %s", subscription.GetID(), subscription.Title)
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

	z.log.Infof("  remove bogus items...")

	return db.Update(func(tx Transaction) error {

		feedExists := z.makeLookupFeed(tx)
		badIDs := []string{}

		c := tx.Bucket(bucketData, entityItem).Cursor()

		item := &Item{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			item.clear()
			if err := item.decode(v); err == nil {
				if !feedExists(item.FeedID) {
					z.log.Infof("item without feed: %s (%s %s)", item.FeedID, item.GetID(), item.GUID)
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

	z.log.Infof("  remove bogus subscriptions...")

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
					z.log.Infof("subscription without user: %s (%s %s)", subscription.UserID, subscription.GetID(), subscription.Title)
					badIDs = append(badIDs, subscription.GetID())
				} else if !feedExists(subscription.FeedID) {
					z.log.Infof("subscription without feed: %s (%s %s)", subscription.FeedID, subscription.GetID(), subscription.Title)
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

	z.log.Infof("  remove bogus transmissions...")

	return db.Update(func(tx Transaction) error {

		feedExists := z.makeLookupFeed(tx)
		badIDs := []string{}

		c := tx.Bucket(bucketData, entityTransmission).Cursor()

		transmission := &Transmission{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			transmission.clear()
			if err := transmission.decode(v); err == nil {
				if !feedExists(transmission.FeedID) {
					z.log.Infof("transmission without feed: %s (%s)", transmission.FeedID, transmission.GetID())
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

	z.log.Infof("  remove feeds without subscriptions...")

	return db.Update(func(tx Transaction) error {

		subscriptionExists := z.makeLookupSubscriptionByFeedID(tx)

		badIDs := []string{}
		c := tx.Bucket(bucketData, entityFeed).Cursor()

		feed := &Feed{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			feed.clear()
			if err := feed.decode(v); err == nil {
				if !subscriptionExists(feed.ID) {
					z.log.Infof("feed without subscription: %s %s", feed.ID, feed.Title)
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

	z.log.Infof("  remove subscription duplicates...")

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

func (z *boltInstance) warnGroupsWithSameName(db Database) error {

	z.log.Infof("  warn groups with same name...")

	return db.Update(func(tx Transaction) error {

		bGroups := tx.Bucket(bucketData, entityGroup)
		bTmp := z.createTempBucket(tx, entityGroup)

		if err := z.indexGroupsByUserIDName(bGroups, bTmp); err != nil {
			return err
		}

		c := bTmp.Cursor()
		groups := Groups{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := groups.decode(v); err != nil {
				return err
			}
			if len(groups) > 1 {
				z.log.Infof("multiple groups with same name: %s", k)
			}
		} // loop

		return nil

	})

}

func (z *boltInstance) warnItemsWithSameGUID(db Database) error {

	z.log.Infof("  warn item with same guid...")

	return db.Update(func(tx Transaction) error {

		bItems := tx.Bucket(bucketData, entityItem)
		bTmp := z.createTempBucket(tx, entityItem)

		if err := z.indexItemsByFeedIDGUID(bItems, bTmp); err != nil {
			return err
		}

		c := bTmp.Cursor()
		items := Items{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := items.decode(v); err != nil {
				return err
			}
			if len(items) > 1 {
				z.log.Infof("multiple items with same GUID: %02d %s - %s", len(items), items[0].Created.Format(time.RFC3339), k)
			}
		} // loop

		return nil

	})

}

func (z *boltInstance) warnUsersWithSameUsername(db Database) error {

	z.log.Infof("  warn users with same username...")

	return db.Update(func(tx Transaction) error {

		bUsers := tx.Bucket(bucketData, entityUser)
		bTmp := z.createTempBucket(tx, entityUser)

		if err := z.indexUsersByUsernameLowecase(bUsers, bTmp); err != nil {
			return err
		}

		c := bTmp.Cursor()
		users := Users{}
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if err := users.decode(v); err != nil {
				return err
			}
			if len(users) > 1 {
				z.log.Infof("multiple users with same name: %s", k)
			}
		} // loop

		return nil

	})

}
