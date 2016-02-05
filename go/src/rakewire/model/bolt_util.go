package model

import (
	"fmt"
	"github.com/boltdb/bolt"
	semver "github.com/hashicorp/go-version"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// top level buckets
const (
	bucketConfig    = "Config"
	bucketData      = "Data"
	bucketIndex     = "Index"
	databaseVersion = "db.version"
)

func updateDatabaseVersion(db Database, version string) error {

	return db.Update(func(tx Transaction) error {

		cfg := NewConfiguration()
		if err := cfg.Load(tx); err != nil {
			return err
		}

		cfg.Set(databaseVersion, version)
		return cfg.Save(tx)

	})

}

func checkSchema(tx *bolt.Tx) error {

	var b *bolt.Bucket

	// top level
	bConfig, err := tx.CreateBucketIfNotExists([]byte(bucketConfig))
	if err != nil {
		return err
	}
	bData, err := tx.CreateBucketIfNotExists([]byte(bucketData))
	if err != nil {
		return err
	}
	bIndex, err := tx.CreateBucketIfNotExists([]byte(bucketIndex))
	if err != nil {
		return err
	}

	// data & indexes
	for entityName, entityIndexes := range allEntities {
		if err := addEntryIfNotExists(bConfig, "sequence."+strings.ToLower(entityName), "0"); err != nil {
			return err
		}
		if _, err = bData.CreateBucketIfNotExists([]byte(entityName)); err != nil {
			return err
		}
		if b, err = bIndex.CreateBucketIfNotExists([]byte(entityName)); err == nil {
			for _, indexName := range entityIndexes {
				if _, err = b.CreateBucketIfNotExists([]byte(indexName)); err != nil {
					return err
				}
			} // entityIndexes
		} else {
			return err
		}
	} // allEntities

	return nil

}

// Returns 0.0.0 if an error occurs
func getDatabaseVersion(location string) string {

	version := "0.0.0"

	boltDB, err := bolt.Open(location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if boltDB != nil {
		defer boltDB.Close()
	}
	if err == nil {
		boltDB.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucketConfig))
			if b != nil {
				ver := b.Get([]byte(databaseVersion))
				if ver != nil && len(ver) > 0 {
					versionStr := string(ver)
					if _, err := semver.NewVersion(versionStr); err == nil {
						version = versionStr
					}
				}
			}
			return nil
		})
	}

	return version

}

func checkIntegrity(location string) error {

	log.Println("checking database integrity...")

	// rename database file to backup name, create new file, open both files
	backupFilename, err := backupDatabase(location)
	if err != nil {
		return err
	}
	log.Printf("original database saved to %s\n", backupFilename)

	oldBoltDB, err := bolt.Open(backupFilename, 0400, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer oldBoltDB.Close()

	newBoltDB, err := bolt.Open(location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer newBoltDB.Close()

	// ensure correct buckets exist in new file
	log.Print("ensuring database structure...")
	if err := oldBoltDB.Update(func(tx *bolt.Tx) error {
		return checkSchema(tx)
	}); err != nil {
		oldBoltDB.Close()
		return err
	}
	if err := newBoltDB.Update(func(tx *bolt.Tx) error {
		return checkSchema(tx)
	}); err != nil {
		newBoltDB.Close()
		return err
	}
	log.Println("ensuring database structure...finished.")

	oldDB := newBoltDatabase(oldBoltDB)
	newDB := newBoltDatabase(newBoltDB)

	// copy all kv pairs in data buckets to new file, only if they are valid, report invalid records
	log.Println("migrating data...")
	if err := copyBuckets(oldDB, newDB); err != nil {
		return err
	}
	if err := copyContainers(oldDB, newDB); err != nil {
		return err
	}
	log.Println("migrating data...complete")

	log.Println("validating data...")
	if err := checkSubscriptions(newDB); err != nil {
		return err
	}
	if err := checkFeeds(newDB); err != nil {
		return err
	}
	if err := checkFeedDuplicates(newDB); err != nil {
		return err
	}
	if err := checkGroups(newDB); err != nil {
		return err
	}
	if err := checkSubscriptions(newDB); err != nil {
		return err
	}
	if err := checkFeeds(newDB); err != nil {
		return err
	}
	if err := checkItems(newDB); err != nil {
		return err
	}
	if err := checkTransmissions(newDB); err != nil {
		return err
	}
	if err := checkEntries(newDB); err != nil {
		return err
	}
	log.Println("validating data...complete")

	log.Println("rebuilding indexes...")
	if err := rebuildIndexes(newDB); err != nil {
		return err
	}
	log.Println("rebuilding indexes...complete")

	log.Println("update schema version...")
	if err := updateDatabaseVersion(newDB, Version); err != nil {
		return err
	}

	log.Println("checking database integrity...complete")

	return nil

}

func addEntryIfNotExists(b *bolt.Bucket, key, value string) error {

	v := b.Get([]byte(key))
	if v == nil || len(v) == 0 {
		if err := b.Put([]byte(key), []byte(value)); err != nil {
			return err
		}
	}

	return nil

}

func copyBuckets(src, dst Database) error {

	bucketNames := []string{"Config"}

	for _, bucketName := range bucketNames {
		log.Printf("  %s...", bucketName)
		err := src.Select(func(srcTx Transaction) error {
			srcBucket := srcTx.Bucket(bucketName)
			return dst.Update(func(dstTx Transaction) error {
				dstBucket := dstTx.Bucket(bucketName)
				srcCursor := srcBucket.Cursor()
				for k, v := srcCursor.First(); k != nil; k, v = srcCursor.Next() {
					if err := dstBucket.Put(k, v); err != nil {
						return err
					}
				}
				return nil
			})
		})
		if err != nil {
			return err
		}
	}

	return nil

}

func copyContainers(src, dst Database) error {

	containers := make(map[string]Object)
	containers[entryEntity] = &Entry{}
	containers[feedEntity] = &Feed{}
	containers[groupEntity] = &Group{}
	containers[itemEntity] = &Item{}
	containers[subscriptionEntity] = &Subscription{}
	containers[transmissionEntity] = &Transmission{}
	containers[userEntity] = &User{}

	for entityName, entity := range containers {
		var maxID uint64
		log.Printf("  %s...", entityName)
		err := src.Select(func(srcTx Transaction) error {
			srcContainer, err := srcTx.Container(bucketData, entityName)
			if err != nil {
				return err
			}
			return dst.Update(func(dstTx Transaction) error {
				dstContainer, err := dstTx.Container(bucketData, entityName)
				if err != nil {
					return err
				}
				return srcContainer.Iterate(func(record Record) error {
					entity.clear()
					if err := entity.deserialize(record, true); err != nil {
						log.Printf("  Error in record (%d): %s", record.GetID(), err.Error())
						return nil
					}
					if record.GetID() > maxID {
						maxID = record.GetID()
					}
					return dstContainer.Put(record)
				}, true)
			})
		})
		if err != nil {
			return err
		}
	}

	return nil

}

func checkEntries(db Database) error {

	log.Printf("  %s...\n", entryEntity)

	return db.Update(func(tx Transaction) error {

		userCache := make(map[uint64]Record)
		users, err := tx.Container(bucketData, userEntity)
		if err != nil {
			return err
		}
		userExists := func(userID uint64) bool {
			if _, ok := userCache[userID]; !ok {
				if u, err := users.Get(userID); err == nil {
					userCache[userID] = u
				}
			}
			u, _ := userCache[userID]
			return u != nil
		}

		itemCache := make(map[uint64]Record)
		items, err := tx.Container(bucketData, itemEntity)
		if err != nil {
			return err
		}
		itemExists := func(itemID uint64) bool {
			if _, ok := itemCache[itemID]; !ok {
				if i, err := items.Get(itemID); err == nil {
					itemCache[itemID] = i
				}
			}
			i, _ := itemCache[itemID]
			return i != nil
		}

		subscriptionCache := make(map[uint64]Record)
		subscriptions, err := tx.Container(bucketData, subscriptionEntity)
		if err != nil {
			return err
		}
		subscriptionExists := func(subscriptionID uint64) bool {
			if _, ok := subscriptionCache[subscriptionID]; !ok {
				if s, err := subscriptions.Get(subscriptionID); err == nil {
					subscriptionCache[subscriptionID] = s
				}
			}
			s, _ := subscriptionCache[subscriptionID]
			return s != nil
		}

		badIDs := []uint64{}

		entries, err := tx.Container(bucketData, entryEntity)
		if err != nil {
			return err
		}

		entry := &Entry{}
		entries.Iterate(func(record Record) error {
			entry.clear()
			if err := entry.deserialize(record); err == nil {
				if !userExists(entry.UserID) {
					log.Printf("    entry without user: %d", entry.ID)
					badIDs = append(badIDs, entry.ID)
				} else if !itemExists(entry.ItemID) {
					log.Printf("    entry without item: %d", entry.ID)
					badIDs = append(badIDs, entry.ID)
				} else if !subscriptionExists(entry.SubscriptionID) {
					log.Printf("    entry without subscription: %d", entry.ID)
					badIDs = append(badIDs, entry.ID)
				}
			}
			return nil
		})

		// remove bad entries
		for _, id := range badIDs {
			// use container, because operating on database without indexes yet
			if err := entries.Delete(id); err != nil {
				return err
			}
		}

		return nil

	})

}

func checkFeedDuplicates(db Database) error {

	log.Printf("  %s duplicates...\n", feedEntity)

	return db.Update(func(tx Transaction) error {

		subscriptions, err := tx.Container(bucketData, subscriptionEntity)
		if err != nil {
			return err
		}
		findSubscriptions := func(feedID uint64) Subscriptions {
			result := Subscriptions{}
			subscriptions.Iterate(func(record Record) error {
				subscription := &Subscription{}
				if err := subscription.deserialize(record); err == nil {
					if subscription.FeedID == feedID {
						result = append(result, subscription)
					}
				}
				return nil
			})
			return result
		}

		feeds, err := FeedsAll(tx)
		if err != nil {
			return err
		}

		feedsByURL := feeds.GroupAllByURL()
		for url, feeds := range feedsByURL {
			if feeds.Len() > 1 {
				log.Printf("    multiple feeds for url: %s", url)

				feeds.SortByID()
				masterFeedID := feeds.First().ID

				// find subscriptions for remaining feeds, update with master feedID
				for _, feed := range feeds[1:] {
					subs := findSubscriptions(feed.ID)
					for _, s := range subs {
						log.Printf("      updating subscription for url: %d", s.ID)
						s.FeedID = masterFeedID
						if err := subscriptions.Put(s.serialize()); err != nil {
							return err
						}
					} // loop thru subscriptions
				} // loop thru remaining feeds

			} // multiple feeds for url
		} // loop feeds by url

		return nil
	})

}

func checkFeeds(db Database) error {

	log.Printf("  %s...\n", feedEntity)

	return db.Update(func(tx Transaction) error {

		subscriptions, err := tx.Container(bucketData, subscriptionEntity)
		if err != nil {
			return err
		}
		subscriptionExists := func(feedID uint64) bool {
			result := false
			subscription := &Subscription{}
			subscriptions.Iterate(func(record Record) error {
				subscription.clear()
				if err := subscription.deserialize(record); err == nil {
					if subscription.FeedID == feedID {
						result = true
						return fmt.Errorf("    found matching subscription: %d", subscription.ID)
					}
				}
				return nil
			})
			return result
		}

		badIDs := []uint64{}
		feeds, err := tx.Container(bucketData, feedEntity)
		if err != nil {
			return err
		}

		feed := &Feed{}
		feeds.Iterate(func(record Record) error {
			feed.clear()
			if err := feed.deserialize(record); err == nil {
				if !subscriptionExists(feed.ID) {
					log.Printf("    feed without subscription: %d %s", feed.ID, feed.Title)
					badIDs = append(badIDs, feed.ID)
				}
			}
			return nil
		})

		// remove bad feeds
		for _, id := range badIDs {
			// use container, because operating on database without indexes yet
			if err := feeds.Delete(id); err != nil {
				return err
			}
		}

		return nil

	})

}

func checkGroups(db Database) error {

	log.Printf("  %s...\n", groupEntity)

	return db.Update(func(tx Transaction) error {

		userCache := make(map[uint64]Record)
		users, err := tx.Container(bucketData, userEntity)
		if err != nil {
			return err
		}
		userExists := func(userID uint64) bool {
			if _, ok := userCache[userID]; !ok {
				if u, err := users.Get(userID); err == nil {
					userCache[userID] = u
				}
			}
			u, _ := userCache[userID]
			return u != nil
		}

		badIDs := []uint64{}

		groups, err := tx.Container(bucketData, groupEntity)
		if err != nil {
			return err
		}

		group := &Group{}
		groups.Iterate(func(record Record) error {
			group.clear()
			if err := group.deserialize(record); err == nil {
				if !userExists(group.UserID) {
					log.Printf("    group without user: %d %d %s", group.ID, group.UserID, group.Name)
					badIDs = append(badIDs, group.ID)
				}
			}
			return nil
		})

		// remove bad groups
		for _, id := range badIDs {
			if err := groups.Delete(id); err != nil {
				return err
			}
		}

		return nil

	})

}

func checkItems(db Database) error {

	log.Printf("  %s...\n", itemEntity)

	return db.Update(func(tx Transaction) error {

		feedCache := make(map[uint64]Record)
		feeds, err := tx.Container(bucketData, feedEntity)
		if err != nil {
			return err
		}
		feedExists := func(feedID uint64) bool {
			if _, ok := feedCache[feedID]; !ok {
				if u, err := feeds.Get(feedID); err == nil {
					feedCache[feedID] = u
				}
			}
			f, _ := feedCache[feedID]
			return f != nil
		}

		badIDs := []uint64{}
		items, err := tx.Container(bucketData, itemEntity)
		if err != nil {
			return err
		}

		item := &Item{}
		items.Iterate(func(record Record) error {
			item.clear()
			if err := item.deserialize(record); err == nil {
				if !feedExists(item.FeedID) {
					log.Printf("    item without feed: %d (%d %s)", item.FeedID, item.ID, item.Title)
					badIDs = append(badIDs, item.ID)
				}
			}
			return nil
		})

		// remove bad items
		for _, id := range badIDs {
			// use container, because operating on database without indexes yet
			if err := items.Delete(id); err != nil {
				return err
			}
		}

		return nil

	})

}

func checkSubscriptions(db Database) error {

	log.Printf("  %s...\n", subscriptionEntity)

	return db.Update(func(tx Transaction) error {

		userCache := make(map[uint64]Record)
		users, err := tx.Container(bucketData, userEntity)
		if err != nil {
			return err
		}
		userExists := func(userID uint64) bool {
			if _, ok := userCache[userID]; !ok {
				if u, err := users.Get(userID); err == nil {
					userCache[userID] = u
				}
			}
			u, _ := userCache[userID]
			return u != nil
		}

		feedCache := make(map[uint64]Record)
		feeds, err := tx.Container(bucketData, feedEntity)
		if err != nil {
			return err
		}
		feedExists := func(feedID uint64) bool {
			if _, ok := feedCache[feedID]; !ok {
				if u, err := feeds.Get(feedID); err == nil {
					feedCache[feedID] = u
				}
			}
			f, _ := feedCache[feedID]
			return f != nil
		}

		groupCache := make(map[uint64]*Group)
		groups, err := tx.Container(bucketData, groupEntity)
		if err != nil {
			return err
		}
		groupExists := func(groupID, userID uint64) bool {
			if _, ok := groupCache[groupID]; !ok {
				if record, err := groups.Get(groupID); err == nil {
					if record != nil {
						g := &Group{}
						if err := g.deserialize(record); err == nil {
							groupCache[groupID] = g
						}
					} else {
						groupCache[groupID] = nil
					}
				}
			}
			g, _ := groupCache[groupID]
			return g != nil && g.UserID == userID
		}

		badIDs := []uint64{}
		cleanedSubscriptions := Subscriptions{}

		subscriptions, err := tx.Container(bucketData, subscriptionEntity)
		if err != nil {
			return err
		}

		subscription := &Subscription{}
		subscriptions.Iterate(func(record Record) error {
			subscription.clear()
			if err := subscription.deserialize(record); err == nil {
				if !userExists(subscription.UserID) {
					log.Printf("    subscription without user: %d (%d %s)", subscription.UserID, subscription.ID, subscription.Title)
					badIDs = append(badIDs, subscription.ID)
				} else if !feedExists(subscription.FeedID) {
					log.Printf("    subscription without feed: %d (%d %s)", subscription.FeedID, subscription.ID, subscription.Title)
					badIDs = append(badIDs, subscription.ID)
				} else {

					// remove invalid groups
					invalidGroupIDs := []uint64{}
					for _, groupID := range subscription.GroupIDs {
						if !groupExists(groupID, subscription.UserID) {
							log.Printf("    subscription with invalid group: %d (%d %s)", groupID, subscription.ID, subscription.Title)
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
						log.Printf("    subscription without groups: %d %s", subscription.ID, subscription.Title)
					}

				}
			}
			return nil
		})

		// remove bad subscriptions
		for _, id := range badIDs {
			// use container, because operating on database without indexes yet
			if err := subscriptions.Delete(id); err != nil {
				return err
			}
		}

		// resave cleaned subscriptions
		for _, subscription := range cleanedSubscriptions {
			// use container, because operating on database without indexes yet
			if err := subscriptions.Put(subscription.serialize()); err != nil {
				return err
			}
		}

		return nil

	})

}

func checkTransmissions(db Database) error {

	log.Printf("  %s...\n", transmissionEntity)

	return db.Update(func(tx Transaction) error {

		feedCache := make(map[uint64]Record)
		feeds, err := tx.Container(bucketData, feedEntity)
		if err != nil {
			return err
		}
		feedExists := func(feedID uint64) bool {
			if _, ok := feedCache[feedID]; !ok {
				if u, err := feeds.Get(feedID); err == nil {
					feedCache[feedID] = u
				}
			}
			f, _ := feedCache[feedID]
			return f != nil
		}

		badIDs := []uint64{}
		transmissions, err := tx.Container(bucketData, transmissionEntity)
		if err != nil {
			return err
		}

		transmission := &Transmission{}
		transmissions.Iterate(func(record Record) error {
			transmission.clear()
			if err := transmission.deserialize(record); err == nil {
				if !feedExists(transmission.FeedID) {
					log.Printf("    transmission without feed: %d (%d %s)", transmission.FeedID, transmission.ID, transmission.URL)
					badIDs = append(badIDs, transmission.ID)
				}
			}
			return nil
		})

		// remove bad transmissions
		for _, id := range badIDs {
			// use container, because operating on database without indexes yet
			if err := transmissions.Delete(id); err != nil {
				return err
			}
		}

		return nil

	})

}

func rebuildIndexes(db Database) error {

	containers := make(map[string]Object)
	containers[entryEntity] = &Entry{}
	containers[feedEntity] = &Feed{}
	containers[groupEntity] = &Group{}
	containers[itemEntity] = &Item{}
	containers[subscriptionEntity] = &Subscription{}
	containers[transmissionEntity] = &Transmission{}
	containers[userEntity] = &User{}

	for entityName, entity := range containers {
		log.Printf("  %s...", entityName)

		err := db.Update(func(tx Transaction) error {
			container, err := tx.Container(bucketData, entityName)
			if err != nil {
				return err
			}
			return container.Iterate(func(record Record) error {
				entity.clear()
				if err := entity.deserialize(record); err != nil {
					return err
				}
				if err := kvSaveIndexes(entityName, entity.getID(), entity.indexKeys(), nil, tx); err != nil {
					return err
				}
				return nil
			})
		})

		if err != nil {
			return err
		}

	} // loop entities

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
