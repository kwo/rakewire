package model

import (
	"github.com/boltdb/bolt"
	"log"
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
	b, err = bucketData.CreateBucketIfNotExists([]byte(UserEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(GroupEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(FeedEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(FeedLogEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(ItemEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(EntryEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(SubscriptionEntity))
	if err != nil {
		return err
	}

	// indexes

	user := NewUser("")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(UserEntity))
	if err != nil {
		return err
	}
	for k := range user.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	group := NewGroup(0, "")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(GroupEntity))
	if err != nil {
		return err
	}
	for k := range group.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	feed := NewFeed("")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(FeedEntity))
	if err != nil {
		return err
	}
	for k := range feed.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	feedlog := NewFeedLog(feed.ID)
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(FeedLogEntity))
	if err != nil {
		return err
	}
	for k := range feedlog.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	item := NewItem(feed.ID, "")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(ItemEntity))
	if err != nil {
		return err
	}
	for k := range item.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	ue := Entry{}
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(EntryEntity))
	if err != nil {
		return err
	}
	for k := range ue.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	uf := NewSubscription(user.ID, feed.ID)
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(SubscriptionEntity))
	if err != nil {
		return err
	}
	for k := range uf.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	return nil

}

func removeInvalidKeys(tx Transaction) error {

	log.Printf("%-7s %-7s remove invalid items...", logInfo, logName)

	if err := removeInvalidKeysForEntity(UserEntity, &User{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(GroupEntity, &Group{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(FeedEntity, &Feed{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(FeedLogEntity, &FeedLog{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(ItemEntity, &Item{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(EntryEntity, &Entry{}, tx); err != nil {
		return err
	}

	if err := removeInvalidKeysForEntity(SubscriptionEntity, &Subscription{}, tx); err != nil {
		return err
	}

	log.Printf("%-7s %-7s remove invalid items complete", logInfo, logName)

	return nil
}

func removeInvalidKeysForEntity(entityName string, dao DataObject, tx Transaction) error {

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
			if err := dao.Deserialize(data, true); err != nil {
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
			dao.Clear()

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

	if err := rebuildIndexesForEntity(UserEntity, &User{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(GroupEntity, &Group{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(FeedEntity, &Feed{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(FeedLogEntity, &FeedLog{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(ItemEntity, &Item{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(EntryEntity, &Entry{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(SubscriptionEntity, &Subscription{}, tx); err != nil {
		return err
	}

	log.Printf("%-7s %-7s rebuilding indexes complete", logInfo, logName)

	return nil

}

func rebuildIndexesForEntity(entityName string, dao DataObject, tx Transaction) error {

	bEntity := tx.Bucket(bucketData).Bucket(entityName)
	ids, err := kvGetUniqueIDs(bEntity) // TODO: what about really large buckets
	if err != nil {
		return err
	}

	for _, id := range ids {
		dao.Clear()
		if data, ok := kvGet(id, bEntity); ok {
			if err := dao.Deserialize(data); err != nil {
				return err
			}
			if err := kvSaveIndexes(entityName, id, dao.IndexKeys(), nil, tx); err != nil {
				return err
			}
		}
	}

	return nil
}
