package model

import (
	"github.com/boltdb/bolt"
	"strings"
	"sync"
	"time"
)

const (
	containerSeparator = "/"
)

// OpenDatabase opens the database at the specified location
func OpenDatabase(location string, flags ...bool) (Database, error) {

	flagCheckIntegrity := len(flags) > 0 && flags[0]

	if flagCheckIntegrity {
		return nil, checkIntegrity(location)
	}

	boltDB, err := bolt.Open(location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	if err := boltDB.Update(func(tx *bolt.Tx) error {
		if err := checkSchema(tx); err != nil {
			return err
		}
		if err := upgradeSchema(tx); err != nil {
			return err
		}
		return nil
	}); err != nil {
		boltDB.Close()
		return nil, err
	}

	return newBoltDatabase(boltDB), nil

}

// CloseDatabase properly closes database resource
func CloseDatabase(d Database) error {

	boltDB := d.(*boltDatabase).db

	if err := boltDB.Close(); err != nil {
		return err
	}

	return nil

}

func newBoltDatabase(boltDB *bolt.DB) Database {
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

func (z *boltTransaction) Container(path string) Container {
	names := strings.Split(path, containerSeparator)
	var b *bolt.Bucket
	for _, name := range names {
		if b == nil {
			b = z.tx.Bucket([]byte(name))
		} else {
			b = b.Bucket([]byte(name))
		}
	}
	return &boltContainer{bucket: b}
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

type boltContainer struct {
	bucket *bolt.Bucket
}

func (z *boltContainer) Container(path string) Container {
	names := strings.Split(path, containerSeparator)
	b := z.bucket
	for _, name := range names {
		b = b.Bucket([]byte(name))
	}
	return &boltContainer{bucket: b}
}

func (z *boltContainer) Iterate(onRecord OnRecord, flags ...bool) error {

	flagIgnoreBaddies := len(flags) > 0 && flags[0]

	firstRow := false
	var lastID uint64
	record := make(Record)

	err := z.bucket.ForEach(func(key, value []byte) error {

		id, fieldname, err := kvBucketKeyDecode(key)
		if err != nil && !flagIgnoreBaddies {
			return err
		}

		if !firstRow {
			lastID = id
			firstRow = true
		}

		if id != lastID {
			if err := onRecord(record); err != nil {
				return err
			}
			// reset
			lastID = id
			record = make(Record)
		} // id switch

		record[fieldname] = string(value)
		return nil

	}) // for each

	if err != nil {
		return err
	}

	// fire last one
	if len(record) > 0 {
		if err := onRecord(record); err != nil {
			return err
		}
	}

	return nil

}

func (z *boltContainer) Put(id uint64, record Record) error {
	for fieldname, v := range record {
		key := kvBucketKeyEncode(id, fieldname)
		value := []byte(v)
		if err := z.bucket.Put(key, value); err != nil {
			return err
		}
	}
	return nil
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
