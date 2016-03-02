package store

import (
	"github.com/boltdb/bolt"
	"sync"
	"time"
)

type boltInstance struct{}

// Open opens the store at the specified location
func (z *boltInstance) Open(location string) (Store, error) {

	boltDB, err := bolt.Open(location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &boltStore{db: boltDB}, nil

}

// Close properly closes store resource
func (z *boltInstance) Close(d Store) error {

	if d == nil {
		return nil
	}

	boltStore := d.(*boltStore).db

	if err := boltStore.Close(); err != nil {
		return err
	}

	return nil

}

type boltStore struct {
	sync.Mutex
	db *bolt.DB
}

func (z *boltStore) Location() string {
	return z.db.Path()
}

func (z *boltStore) Select(fn func(tx Transaction) error) error {
	return z.db.View(func(tx *bolt.Tx) error {
		bt := &boltTransaction{tx: tx}
		return fn(bt)
	})
}

func (z *boltStore) Update(fn func(transaction Transaction) error) error {
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
