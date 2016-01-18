package model

import (
	"github.com/boltdb/bolt"
	"sync"
	"time"
)

// top lebel buckets
const (
	bucketConfig = "Config"
	bucketData   = "Data"
	bucketIndex  = "Index"
)

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
