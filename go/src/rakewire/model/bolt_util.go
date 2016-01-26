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

	// rename database file to backup name, create new file, open both files
	newLocation, err := renameWithTimestamp(location)
	if err != nil {
		return err
	}
	log.Printf("original database saved to %s\n", newLocation)

	oldDB, err := bolt.Open(location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer oldDB.Close()

	newDB, err := bolt.Open(newLocation, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	defer newDB.Close()

	// ensure correct buckets exist in new file
	log.Print("ensuring database structure...")
	if err := oldDB.Update(func(tx *bolt.Tx) error {
		return checkSchema(tx)
	}); err != nil {
		oldDB.Close()
		return err
	}
	if err := newDB.Update(func(tx *bolt.Tx) error {
		return checkSchema(tx)
	}); err != nil {
		newDB.Close()
		return err
	}
	log.Println("ensuring database structure...finished.")

	// - copy all kv pairs in data buckets to new file, only if they are valid, report invalid records

	return nil

}

func removeInvalidKeys(tx Transaction) error {

	log.Printf("%-7s %-7s remove invalid items...", logInfo, logName)

	if err := removeInvalidKeysForEntity(userEntity, &User{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(groupEntity, &Group{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(feedEntity, &Feed{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(transmissionEntity, &Transmission{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(itemEntity, &Item{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(entryEntity, &Entry{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(subscriptionEntity, &Subscription{}, tx); err != nil {
		return err
	}

	log.Printf("%-7s %-7s remove invalid items complete", logInfo, logName)

	return nil
}

func removeInvalidKeysForEntity(entityName string, dao Object, tx Transaction) error {

	b := tx.Bucket(bucketData).Bucket(entityName)

	invalidKeys := [][]byte{}
	var lastID uint64
	data := make(map[string]string)
	keys := [][]byte{}
	firstRow := false
	err := b.ForEach(func(k, v []byte) error {

		id, _, err := kvBucketKeyDecode(k)
		if err != nil {
			invalidKeys = append(invalidKeys, k)
			return nil
		}

		if !firstRow {
			lastID = id
			firstRow = true
		}

		if id != lastID {

			// deserialize dao
			if err := dao.deserialize(data, true); err != nil {
				if derr, ok := err.(*DeserializationError); ok {
					if len(derr.Errors) > 0 || len(derr.MissingFieldnames) > 0 {
						// invalid item, remove complete record, all keys
						for _, key := range keys {
							invalidKeys = append(invalidKeys, key)
						}
					} else if len(derr.UnknownFieldnames) > 0 {
						// only remove unknown keys
						for _, fieldname := range derr.UnknownFieldnames {
							invalidKeys = append(invalidKeys, kvBucketKeyEncode(id, fieldname))
						}
					}
				} else {
					return err
				}
			}

			// reset
			lastID = id
			data = make(map[string]string)
			keys = [][]byte{}
			dao.clear()

		} // id switch

		data[kvKeyElement(k, 1)] = string(v)
		keys = append(keys, k)

		return nil

	})

	if err != nil {
		return err
	}

	// remove invalid keys
	for _, key := range invalidKeys {
		log.Printf("%-7s %-7s removing invalid item %s: %s", logDebug, logName, entityName, key)
		if err := b.Delete(key); err != nil {
			return err
		}
	}

	return nil

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
