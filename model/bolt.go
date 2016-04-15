package model

import (
	"github.com/boltdb/bolt"
	"rakewire/logger"
	"sync"
	"time"
)

const (
	bucketTmp = "tmp"
)

// Instance allows opening and closing new Databases
var Instance = &boltInstance{
	log: logger.New("db"),
}

type boltInstance struct {
	log *logger.Logger
}

// Open opens the store at the specified location
func (z *boltInstance) Open(location string) (Database, error) {

	boltDB, err := bolt.Open(location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	err = boltDB.Update(func(tx *bolt.Tx) error {
		return z.checkSchema(tx)
	})
	if err != nil {
		boltDB.Close()
		return nil, err
	}

	return &boltDatabase{db: boltDB}, nil

}

// Close properly closes store resource
func (z *boltInstance) Close(db Database) error {

	if db == nil {
		return nil
	}

	boltDB := db.(*boltDatabase).db

	if err := boltDB.Close(); err != nil {
		return err
	}

	return nil

}

func (z *boltInstance) checkSchema(tx *bolt.Tx) error {

	var b *bolt.Bucket

	// top level
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

func (z *boltInstance) createTempTopLevelBucket(db Database) error {
	boltDb := db.(*boltDatabase).db
	return boltDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketTmp))
		return err
	})
}

func (z *boltInstance) removeTempTopLevelBucket(db Database) error {
	boltDb := db.(*boltDatabase).db
	return boltDb.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucketTmp))
	})
}

func (z *boltInstance) createTempBucket(tx Transaction, bucketName string) Bucket {

	boltTx := tx.(*boltTransaction).tx
	bTmp := boltTx.Bucket([]byte(bucketTmp))

	tempBucketName := []byte(bucketName)
	bTmp.DeleteBucket(tempBucketName) // ignore error
	b, err := bTmp.CreateBucket(tempBucketName)
	if err != nil {
		return nil
	}

	return &boltBucket{bucket: b}

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

func (z *boltTransaction) Bucket(names ...string) Bucket {
	var b *bolt.Bucket
	for i, name := range names {
		if i == 0 {
			b = z.tx.Bucket([]byte(name))
		} else {
			b = b.Bucket([]byte(name))
		}
		if b == nil {
			return nil
		}
	}
	return &boltBucket{bucket: b}
}

type boltBucket struct {
	bucket *bolt.Bucket
}

func (z *boltBucket) Bucket(names ...string) Bucket {
	b := z.bucket
	for _, name := range names {
		b = b.Bucket([]byte(name))
		if b == nil {
			return nil
		}
	}
	return &boltBucket{bucket: b}
}

func (z *boltBucket) Cursor() Cursor {
	cursor := z.bucket.Cursor()
	return &boltCursor{cursor: cursor}
}

func (z *boltBucket) Delete(key []byte) error {
	return z.bucket.Delete(key)
}

func (z *boltBucket) Get(key []byte) []byte {
	return z.bucket.Get(key)
}

func (z *boltBucket) Put(key, value []byte) error {
	return z.bucket.Put(key, value)
}

type boltCursor struct {
	cursor *bolt.Cursor
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
