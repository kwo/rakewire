package model

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// top level buckets
const (
	bucketConfig = "Config"
	bucketData   = "Data"
	bucketIndex  = "Index"
)

func checkSchema(tx *bolt.Tx) error {

	var b *bolt.Bucket

	// top level
	_, err := tx.CreateBucketIfNotExists([]byte(bucketConfig))
	if err != nil {
		return err
	}
	bucketData, err := tx.CreateBucketIfNotExists([]byte(bucketData))
	if err != nil {
		return err
	}
	bucketIndex, err := tx.CreateBucketIfNotExists([]byte(bucketIndex))
	if err != nil {
		return err
	}

	// data
	b, err = bucketData.CreateBucketIfNotExists([]byte(userEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(groupEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(feedEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(transmissionEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(itemEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(entryEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(subscriptionEntity))
	if err != nil {
		return err
	}

	// indexes

	user := NewUser("")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(userEntity))
	if err != nil {
		return err
	}
	for k := range user.indexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	group := NewGroup(0, "")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(groupEntity))
	if err != nil {
		return err
	}
	for k := range group.indexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	feed := NewFeed("")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(feedEntity))
	if err != nil {
		return err
	}
	for k := range feed.indexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	transmission := NewTransmission(feed.ID)
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(transmissionEntity))
	if err != nil {
		return err
	}
	for k := range transmission.indexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	item := NewItem(feed.ID, "")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(itemEntity))
	if err != nil {
		return err
	}
	for k := range item.indexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	ue := Entry{}
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(entryEntity))
	if err != nil {
		return err
	}
	for k := range ue.indexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	uf := NewSubscription(user.ID, feed.ID)
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(subscriptionEntity))
	if err != nil {
		return err
	}
	for k := range uf.indexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	return nil

}

func upgradeSchema(tx *bolt.Tx) error {

	return nil

}

func checkIntegrity(location string) error {

	log.Println("checking database integrity...")

	// rename database file to backup name, create new file, open both files
	backupFilename, err := renameWithTimestamp(location)
	if err != nil {
		return err
	}
	log.Printf("original database saved to %s\n", backupFilename)

	oldBoltDB, err := bolt.Open(backupFilename, 0600, &bolt.Options{Timeout: 1 * time.Second})
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

	log.Println("interrogating groups...")
	if err := checkGroups(newDB); err != nil {
		return err
	}
	log.Println("interrogating groups...complete")

	log.Println("interrogating subscriptions...")
	if err := checkSubscriptions(newDB); err != nil {
		return err
	}
	log.Println("interrogating subscriptions...complete")

	log.Println("checking database integrity...complete")

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
				return srcBucket.ForEach(func(k, v []byte) error {
					return dstBucket.Put(k, v)
				})
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
						log.Printf("Error in record (%d): %s", record.GetID(), err.Error())
						return nil
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

func checkGroups(db Database) error {

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
			if err := group.deserialize(record); err != nil {
				return err
			}
			if !userExists(group.UserID) {
				log.Printf("  group without user: %d %d %s", group.ID, group.UserID, group.Name)
				badIDs = append(badIDs, group.ID)
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

func checkSubscriptions(db Database) error {

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
			if err := subscription.deserialize(record); err != nil {
				return err
			}
			if !userExists(subscription.UserID) {
				log.Printf("  subscription without user: %d (%d %s)", subscription.UserID, subscription.ID, subscription.Title)
				badIDs = append(badIDs, subscription.ID)
			} else if !feedExists(subscription.FeedID) {
				log.Printf("  subscription without feed: %d (%d %s)", subscription.FeedID, subscription.ID, subscription.Title)
				badIDs = append(badIDs, subscription.ID)
			} else {

				// remove invalid groups
				invalidGroupIDs := []uint64{}
				for _, groupID := range subscription.GroupIDs {
					if !groupExists(groupID, subscription.UserID) {
						log.Printf("  subscription with invalid group: %d (%d %s)", groupID, subscription.ID, subscription.Title)
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
					log.Printf("  subscription without groups: %d %s", subscription.ID, subscription.Title)
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

func rebuildIndexes(tx *boltTransaction) error {

	log.Printf("%-7s %-7s rebuilding indexes...", logInfo, logName)

	if err := tx.tx.DeleteBucket([]byte(bucketIndex)); err != nil {
		return err
	}

	if err := checkSchema(tx.tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(userEntity, &User{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(groupEntity, &Group{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(feedEntity, &Feed{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(transmissionEntity, &Transmission{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(itemEntity, &Item{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(entryEntity, &Entry{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(subscriptionEntity, &Subscription{}, tx); err != nil {
		return err
	}

	log.Printf("%-7s %-7s rebuilding indexes complete", logInfo, logName)

	return nil

}

func rebuildIndexesForEntity(entityName string, dao Object, tx Transaction) error {

	bEntity := tx.Bucket(bucketData).Bucket(entityName)
	ids, err := kvGetUniqueIDs(bEntity) // TODO: what about really large buckets
	if err != nil {
		return err
	}

	for _, id := range ids {
		dao.clear()
		if data, ok := kvGet(id, bEntity); ok {
			if err := dao.deserialize(data); err != nil {
				return err
			}
			if err := kvSaveIndexes(entityName, id, dao.indexKeys(), nil, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func renameWithTimestamp(location string) (string, error) {

	now := time.Now().Truncate(time.Second)
	timestamp := now.Format("20060102150405")

	dir := filepath.Dir(location)
	ext := filepath.Ext(location)
	filename := strings.TrimSuffix(filepath.Base(location), ext)

	newFilename := fmt.Sprintf("%s%s%s-%s%s", dir, string(os.PathSeparator), filename, timestamp, ext)
	err := os.Rename(location, newFilename)

	return newFilename, err

}
