package bolt

import (
	"github.com/boltdb/bolt"
	"log"
	"rakewire.com/db"
	m "rakewire.com/model"
	"time"
)

// Database implementation of Database
type Database struct {
	db *bolt.DB
}

// Open the database
func (z *Database) Open(cfg *m.DatabaseConfiguration) error {

	db, err := bolt.Open(cfg.Location, 0600, &bolt.Options{Timeout: 1 * time.Second})
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

// Close the database
func (z *Database) Close() error {
	var db = z.db
	z.db = nil
	return db.Close()
}

// GetFeeds list feeds
func (z *Database) GetFeeds() (map[string]*db.FeedInfo, error) {

	result := make(map[string]*db.FeedInfo)

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("FeedInfo"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			f := db.FeedInfo{}
			if err := f.Unmarshal(v); err != nil {
				return err
			}
			if f.ID != string(k) {
				log.Printf("ID/key mismatch: %s/%s\n", k, f.ID)
			} else {
				result[f.ID] = &f
			}

		} // for

		return nil

	})

	if err != nil {
		return nil, err
	}

	return result, nil

}

// SaveFeeds save feeds
func (z *Database) SaveFeeds(feeds []*db.FeedInfo) (int, error) {

	var counter int
	err := z.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("FeedInfo"))

		for _, f := range feeds {

			data, err := f.Marshal()
			if err != nil {
				return err
			}

			if err := b.Put([]byte(f.ID), data); err != nil {
				return err
			}
			counter++

		} // for

		return nil

	}) // update

	return counter, err

}
