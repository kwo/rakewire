package bolt

import (
	"github.com/boltdb/bolt"
	"log"
	"rakewire.com/db"
	"rakewire.com/logging"
	m "rakewire.com/model"
	"time"
)

// Database implementation of Database
type Database struct {
	db *bolt.DB
}

var (
	logger = logging.New("db")
)

// Open the database
func (z *Database) Open(cfg *m.DatabaseConfiguration) error {

	db, err := bolt.Open(cfg.Location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Printf("Cannot open database: %s", err.Error())
		return err
	}
	z.db = db

	// check that buckets exist
	err = z.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Feed"))
		return err
	})

	if err != nil {
		logger.Printf("Cannot initialize database: %s", err.Error())
		return err
	}

	logger.Printf("Using database at %s\n", cfg.Location)

	return nil

}

// Close the database
func (z *Database) Close() error {
	if db := z.db; db != nil {
		z.db = nil
		if err := db.Close(); err != nil {
			logger.Printf("Error closing database: %s\n", err.Error())
			return err
		}
		logger.Println("Closed database")
	}
	return nil
}

// GetFeeds list feeds
func (z *Database) GetFeeds() (*db.Feeds, error) {

	result := db.NewFeeds()

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Feed"))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			f := db.Feed{}
			if err := f.Decode(v); err != nil {
				return err
			}
			if f.ID != string(k) {
				log.Printf("ID/key mismatch: %s/%s\n", k, f.ID)
			} else {
				result.Add(&f)
			}

		} // for

		return nil

	})

	return result, err

}

// GetFeedByID return feed given UUID
func (z *Database) GetFeedByID(id string) (*db.Feed, error) {

	var result db.Feed

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Feed"))
		data := b.Get([]byte(id))
		result.Decode(data)
		return nil
	})

	return &result, err

}

// SaveFeeds save feeds
func (z *Database) SaveFeeds(feeds *db.Feeds) (int, error) {

	var counter int
	err := z.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Feed"))

		for _, f := range feeds.Values {

			data, err := f.Encode()
			if err != nil {
				return err
			}
			if err = b.Put([]byte(f.ID), data); err != nil {
				return err
			}
			counter++

		} // for

		return nil

	}) // update

	return counter, err

}
