package model

import (
	"github.com/boltdb/bolt"
)

// NewBoltDatabase creates a new Database backed by a bolt DB
func NewBoltDatabase(db *bolt.DB) Database {
	return &BoltDatabase{DB: db}
}

// BoltDatabase implements a Database for boltDB
type BoltDatabase struct {
	DB *bolt.DB
}

// Select retrieve a transaction for read-only operations
func (z *BoltDatabase) Select(fn func(tx Transaction) error) error {
	return z.DB.View(func(tx *bolt.Tx) error {
		bt := &boltTransaction{tx: tx}
		return fn(bt)
	})
}

// Update retrieve a transaction for read-write operations
func (z *BoltDatabase) Update(fn func(transaction Transaction) error) error {
	return z.DB.Update(func(tx *bolt.Tx) error {
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
