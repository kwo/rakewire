package bolt

import (
	"github.com/boltdb/bolt"
	"rakewire.com/db"
	"time"
)

// Database implementation of Database
type Database struct {
	db *bolt.DB
}

func (z *Database) init(dbFilename string) error {

	db, err := bolt.Open(dbFilename, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	z.db = db

	// check that buckets exist
	err = z.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("FeedInfo"))
		return err
	})

	return err

}

func (z *Database) destroy() error {
	var db = z.db
	z.db = nil
	return db.Close()
}

func (z *Database) getFeeds() (map[string]*db.FeedInfo, error) {

	z.db.View(func(tx *bolt.Tx) error {
		//b, _ := tx.CreateBucketIfNotExists([]byte("FeedInfo"))
		//b.Put(key []byte, value []byte)
		return nil
	})

	return nil, nil
}
