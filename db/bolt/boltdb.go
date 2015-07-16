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
	bucketIndex          = "Index"
	bucketIndexFeedByURL = "idxFeedByURL"
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
		_, err := tx.CreateBucketIfNotExists([]byte(bucketFeed))
		if err != nil {
			return err
		}
		b, err := tx.CreateBucketIfNotExists([]byte(bucketIndex))
		if err != nil {
			return err
		}
		_, err = b.CreateBucketIfNotExists([]byte(bucketIndexFeedByURL))
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
		b := tx.Bucket([]byte(bucketFeed))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			f := db.Feed{}
			if err := f.Decode(v); err != nil {
				return err
			}
			if f.ID != string(k) {
				logger.Printf("ID/key mismatch: %s/%s\n", k, f.ID)
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

	var data []byte

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketFeed))
		data = b.Get([]byte(id))
		return nil
	})

	if err != nil {
		return nil, err
	} else if data == nil {
		return nil, nil
	}

	result := db.Feed{}
	err = result.Decode(data)
	return &result, err

}

// GetFeedByURL return feed given url
func (z *Database) GetFeedByURL(url string) (*db.Feed, error) {

	var data []byte

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketFeed))
		i := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(bucketIndexFeedByURL))
		data = i.Get([]byte(url))
		if data != nil {
			data = b.Get(data)
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else if data == nil {
		return nil, nil
	}

	result := db.Feed{}
	err = result.Decode(data)
	return &result, err

}

// SaveFeeds save feeds
func (z *Database) SaveFeeds(feeds *db.Feeds) (int, error) {

	var counter int
	err := z.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucketFeed))
		idxFeedByURL := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(bucketIndexFeedByURL))

		for _, f := range feeds.Values {

			data, err := f.Encode()
			if err != nil {
				return err
			}
			if err = b.Put([]byte(f.ID), data); err != nil {
				return err
			}
			if err = idxFeedByURL.Put([]byte(f.URL), []byte(f.ID)); err != nil {
				return err
			}
			counter++

		} // for

		return nil

	}) // update

	return counter, err

}

// Repair the database
func (z *Database) Repair() error {

	return z.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucketFeed))
		indexes := tx.Bucket([]byte(bucketIndex))
		logger.Printf("dropping index %s\n", bucketIndexFeedByURL)
		err := indexes.DeleteBucket([]byte(bucketIndexFeedByURL))
		if err != nil {
			return err
		}
		logger.Printf("readding index %s\n", bucketIndexFeedByURL)
		idxFeedByURL, err := indexes.CreateBucket([]byte(bucketIndexFeedByURL))
		if err != nil {
			return err
		}

		c := b.Cursor()

		logger.Printf("populating index %s\n", bucketIndexFeedByURL)
		for k, v := c.First(); k != nil; k, v = c.Next() {

			f := db.Feed{}
			if err := f.Decode(v); err != nil {
				return err
			}
			if f.ID != string(k) {
				logger.Printf("ID/key mismatch: %s/%s\n", k, f.ID)
			} else {
				err := idxFeedByURL.Put([]byte(f.URL), []byte(f.ID))
				if err != nil {
					return err
				}
				logger.Printf("added url to index %s: %s\n", bucketIndexFeedByURL, f.URL)
			}

		} // for

		return nil

	}) // update

}
