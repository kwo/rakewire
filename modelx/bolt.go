package model

import (
	"bytes"
	"github.com/boltdb/bolt"
	semver "github.com/hashicorp/go-version"
	"sync"
	"time"
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

	return newBoltDatabase(boltDB), nil

}

// CloseDatabase properly closes database resource
func CloseDatabase(d Database) error {

	if d == nil {
		return nil
	}

	boltDB := d.(*boltDatabase).db

	if err := boltDB.Close(); err != nil {
		return err
	}

	return nil

}

// CheckDatabaseIntegrity checks the database integrity.
func CheckDatabaseIntegrity(location string) error {
	return checkIntegrity(location)
}

// CheckDatabaseIntegrityIf checks the database integrity only if the database version differs from the app version.
func CheckDatabaseIntegrityIf(location string) error {

	dbVersion, err := semver.NewVersion(getDatabaseVersion(location))
	if err != nil {
		return err
	}

	appVersion, err := semver.NewVersion(Version)
	if err != nil {
		return err
	}

	if dbVersion.LessThan(appVersion) {
		return checkIntegrity(location)
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

func (z *boltBucket) Delete(key string) error {
	return z.bucket.Delete([]byte(key))
}

func (z *boltBucket) DeleteRecord(id string) error {

	keys := [][]byte{}

	c := z.bucket.Cursor()
	min, max := kvKeyMinMax(id)
	for k, _ := c.Seek([]byte(min)); k != nil && bytes.Compare(k, []byte(max)) <= 0; k, _ = c.Next() {
		keys = append(keys, k)
		// do not delete in a cursor, it is buggy, sometimes advancing the position
	} // for loop

	for _, k := range keys {
		if err := z.bucket.Delete(k); err != nil {
			return err
		}
	}

	return nil

}

func (z *boltBucket) Get(key string) string {
	return string(z.bucket.Get([]byte(key)))
}

// GetIndex retrieves a Record from the given bucket looking up its ID in the current index bucket.
func (z *boltBucket) GetIndex(b Bucket, id string) Record {

	if value := z.bucket.Get([]byte(id)); value != nil {
		return b.GetRecord(string(value))
	}

	return nil

}

func (z *boltBucket) GetRecord(id string) Record {

	found := false
	record := make(Record)

	c := z.bucket.Cursor()
	min, max := kvKeyMinMax(id)
	for k, v := c.Seek([]byte(min)); k != nil && bytes.Compare(k, []byte(max)) <= 0; k, v = c.Next() {
		// assume proper key format of ID/fieldname
		fieldname := kvKeyDecode(k)[1]
		record[fieldname] = string(v)
		found = true
	} // for loop

	if !found {
		return nil
	}

	return record

}

func (z *boltBucket) Iterate(onRecord OnRecord) error {

	firstRow := false
	var id, lastID string
	record := make(Record)

	err := z.bucket.ForEach(func(key, value []byte) error {

		elements := kvKeyDecode(key)
		id = elements[0]
		fieldname := elements[1]

		if !firstRow {
			lastID = id
			firstRow = true
		}

		if id != lastID {
			if err := onRecord(lastID, record); err != nil {
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
		if err := onRecord(id, record); err != nil {
			return err
		}
	}

	return nil

}

func (z *boltBucket) IterateIndex(b Bucket, minID, maxID string, onRecord OnRecord) error {

	min := []byte(minID)
	max := []byte(maxID)
	c := z.bucket.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		id := string(v)
		if record := b.GetRecord(id); record != nil {
			if err := onRecord(id, record); err != nil {
				return err
			}
		}
	} // for loop

	return nil

}

func (z *boltBucket) Put(key, value string) error {
	return z.bucket.Put([]byte(key), []byte(value))
}

func (z *boltBucket) PutRecord(id string, record Record) error {

	if err := z.DeleteRecord(id); err != nil {
		return err
	}

	for fieldname, v := range record {
		key := []byte(kvKeyEncode(id, fieldname))
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