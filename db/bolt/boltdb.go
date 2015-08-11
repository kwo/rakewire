package bolt

import (
	"github.com/boltdb/bolt"
	"rakewire.com/db"
	"rakewire.com/logging"
	m "rakewire.com/model"
	"time"
)

const (
	bucketFeed           = "Feed"
	bucketFeedLog        = "FeedLog"
	bucketIndex          = "Index"
	bucketIndexFeedByURL = "idxFeedByURL"
	bucketIndexNextFetch = "idxNextFetch"
)

// Database implementation of Database
type Database struct {
	db *bolt.DB
}

var (
	logger = logging.New("db")
)

// Open the database
func (z *Database) Open(cfg *db.Configuration) error {

	db, err := bolt.Open(cfg.Location, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Printf("Cannot open database: %s", err.Error())
		return err
	}
	z.db = db

	// check that buckets exist
	err = z.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketFeed))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(bucketFeedLog))
		if err != nil {
			return err
		}
		b, err := tx.CreateBucketIfNotExists([]byte(bucketIndex))
		if err != nil {
			return err
		}
		if _, err = b.CreateBucketIfNotExists([]byte(bucketIndexFeedByURL)); err != nil {
			return err
		}
		if _, err = b.CreateBucketIfNotExists([]byte(bucketIndexNextFetch)); err != nil {
			return err
		}
		return nil
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

// Repair the database
func (z *Database) Repair() error {

	return z.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucketFeed))
		indexes := tx.Bucket([]byte(bucketIndex))

		logger.Printf("dropping index %s\n", bucketIndexFeedByURL)
		if err := indexes.DeleteBucket([]byte(bucketIndexFeedByURL)); err != nil {
			return err
		}
		logger.Printf("dropping index %s\n", bucketIndexNextFetch)
		if err := indexes.DeleteBucket([]byte(bucketIndexNextFetch)); err != nil {
			return err
		}

		logger.Printf("creating index %s\n", bucketIndexFeedByURL)
		idxFeedByURL, err := indexes.CreateBucket([]byte(bucketIndexFeedByURL))
		if err != nil {
			return err
		}
		logger.Printf("creating index %s\n", bucketIndexNextFetch)
		idxNextFetch, err := indexes.CreateBucket([]byte(bucketIndexNextFetch))
		if err != nil {
			return err
		}

		c := b.Cursor()

		logger.Printf("populating indexes")
		for k, v := c.First(); k != nil; k, v = c.Next() {

			f := &m.Feed{}
			if err := f.Decode(v); err != nil {
				return err
			}

			if err := idxFeedByURL.Put([]byte(f.URL), []byte(f.ID)); err != nil {
				return err
			}

			if err := idxNextFetch.Put([]byte(fetchKey(f)), []byte(f.ID)); err != nil {
				return err
			}

		} // for

		return nil

	}) // update

}
