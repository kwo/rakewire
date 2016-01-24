package model

import (
	"github.com/boltdb/bolt"
	"log"
	"sync"
	"time"
)

// top level buckets
const (
	bucketConfig = "Config"
	bucketData   = "Data"
	bucketIndex  = "Index"
)

func NewBoltDatabase(boltDB *bolt.DB) Database {
	return &boltDatabase{db: boltDB}
}

type boltDatabase struct {
	sync.Mutex
	db *bolt.DB
}

func (z *boltDatabase) Location() string {
	return z.db.Path()
}

func (z *boltDatabase) Select(fn func(tx Transaction) error) error {
	return z.db.View(func(tx *bolt.Tx) error {
		bt := &boltTransaction{tx: tx}
		return fn(bt)
	})
}

func (z *boltDatabase) Update(fn func(transaction Transaction) error) error {
	z.Lock()
	defer z.Unlock()
	return z.db.Update(func(tx *bolt.Tx) error {
		bt := &boltTransaction{tx: tx}
		return fn(bt)
	})
}

// Repair ensures the integrity of the database.
func (z *boltDatabase) Repair() error {
	return z.Update(func(tx Transaction) error {

		if err := removeInvalidEntries(tx); err != nil {
			return err
		}

		if err := rebuildIndexes(tx.(*boltTransaction)); err != nil {
			return err
		}

		return nil

	})
}

type boltTransaction struct {
	tx *bolt.Tx
}

func (z *boltTransaction) Bucket(name string) Bucket {
	bucket := z.tx.Bucket([]byte(name))
	return &boltBucket{bucket: bucket}
}

type boltBucket struct {
	bucket *bolt.Bucket
}

func (z *boltBucket) Bucket(name string) Bucket {
	bucket := z.bucket.Bucket([]byte(name))
	return &boltBucket{bucket: bucket}
}

func (z *boltBucket) Cursor() Cursor {
	cursor := z.bucket.Cursor()
	return &boltCursor{cursor: cursor}
}

func (z *boltBucket) Delete(key []byte) error {
	return z.bucket.Delete(key)
}

func (z *boltBucket) ForEach(fn func(key, value []byte) error) error {
	return z.bucket.ForEach(fn)
}

func (z *boltBucket) Get(key []byte) []byte {
	return z.bucket.Get(key)
}

func (z *boltBucket) NextSequence() (uint64, error) {
	return z.bucket.NextSequence()
}

func (z *boltBucket) Put(key, value []byte) error {
	return z.bucket.Put(key, value)
}

type boltCursor struct {
	cursor *bolt.Cursor
}

func (z *boltCursor) Delete() error {
	return z.cursor.Delete()
}

func (z *boltCursor) First() ([]byte, []byte) {
	return z.cursor.First()
}

func (z *boltCursor) Last() ([]byte, []byte) {
	return z.cursor.Last()
}

func (z *boltCursor) Next() ([]byte, []byte) {
	return z.cursor.Next()
}

func (z *boltCursor) Prev() ([]byte, []byte) {
	return z.cursor.Prev()
}

func (z *boltCursor) Seek(seek []byte) ([]byte, []byte) {
	return z.cursor.Seek(seek)
}

// OpenDatabase opens the database at the specified location
func OpenDatabase(location string) (Database, error) {

	boltDB, err := bolt.Open(location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	err = boltDB.Update(func(tx *bolt.Tx) error {
		return checkSchema(tx)
	})
	if err != nil {
		boltDB.Close()
		return nil, err
	}

	return NewBoltDatabase(boltDB), nil

}

// CloseDatabase properly closes database resource
func CloseDatabase(d Database) error {

	boltDB := d.(*boltDatabase).db

	if err := boltDB.Close(); err != nil {
		return err
	}

	return nil

}

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

	b, err = bucketData.CreateBucketIfNotExists([]byte(EntryEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(UserEntryEntity))
	if err != nil {
		return err
	}

	b, err = bucketData.CreateBucketIfNotExists([]byte(UserFeedEntity))
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

	entry := NewEntry(feed.ID, "")
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(EntryEntity))
	if err != nil {
		return err
	}
	for k := range entry.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	ue := UserEntry{}
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(UserEntryEntity))
	if err != nil {
		return err
	}
	for k := range ue.IndexKeys() {
		_, err = b.CreateBucketIfNotExists([]byte(k))
		if err != nil {
			return err
		}
	}

	uf := NewUserFeed(user.ID, feed.ID)
	b, err = bucketIndex.CreateBucketIfNotExists([]byte(UserFeedEntity))
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

func removeInvalidEntries(tx Transaction) error {

	log.Printf("%-7s %-7s remove invalid entries...", logInfo, logName)

	if err := removeInvalidEntriesForEntity(UserEntity, &User{}, tx); err != nil {
		return err
	}

	if err := removeInvalidEntriesForEntity(GroupEntity, &Group{}, tx); err != nil {
		return err
	}

	if err := removeInvalidEntriesForEntity(FeedEntity, &Feed{}, tx); err != nil {
		return err
	}

	if err := removeInvalidEntriesForEntity(FeedLogEntity, &FeedLog{}, tx); err != nil {
		return err
	}

	if err := removeInvalidEntriesForEntity(EntryEntity, &Entry{}, tx); err != nil {
		return err
	}

	if err := removeInvalidEntriesForEntity(UserEntryEntity, &UserEntry{}, tx); err != nil {
		return err
	}

	if err := removeInvalidEntriesForEntity(UserFeedEntity, &UserFeed{}, tx); err != nil {
		return err
	}

	log.Printf("%-7s %-7s remove invalid entries complete", logInfo, logName)

	return nil
}

func removeInvalidEntriesForEntity(entityName string, dao DataObject, tx Transaction) error {

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
						// invalid entry, remove complete record, all keys
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
		log.Printf("%-7s %-7s removing invalid entry %s: %s", logDebug, logName, entityName, key)
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

	if err := rebuildIndexesForEntity(EntryEntity, &Entry{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(UserEntryEntity, &UserEntry{}, tx); err != nil {
		return err
	}

	if err := rebuildIndexesForEntity(UserFeedEntity, &UserFeed{}, tx); err != nil {
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
