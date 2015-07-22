package bolt

import (
	"bytes"
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

// GetFeeds list feeds
func (z *Database) GetFeeds() (*m.Feeds, error) {

	result := m.NewFeeds()

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketFeed))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			f := m.Feed{}
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

// GetFetchFeeds get feeds to be fetched
func (z *Database) GetFetchFeeds(maxTime *time.Time) (*m.Feeds, error) {

	var max []byte
	if maxTime == nil {
		max = []byte(formatMaxTime(time.Now()))
	} else {
		max = []byte(formatMaxTime(*maxTime))
	}

	result := m.NewFeeds()

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketFeed))
		idxNextFetch := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(bucketIndexNextFetch))
		c := idxNextFetch.Cursor()

		//logger.Printf("max: %s\n", string(max))
		for k, uuid := c.First(); k != nil && bytes.Compare(k, max) <= 0; k, uuid = c.Next() {

			//logger.Printf("key: %s: %s", k, uuid)

			v := b.Get(uuid)
			f := &m.Feed{}
			if err := f.Decode(v); err != nil {
				return err
			}
			result.Add(f)

		} // for

		return nil

	})

	return result, err

}

// GetFeedByID return feed given UUID
func (z *Database) GetFeedByID(id string) (*m.Feed, error) {

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

	result := m.Feed{}
	err = result.Decode(data)
	return &result, err

}

// GetFeedByURL return feed given url
func (z *Database) GetFeedByURL(url string) (*m.Feed, error) {

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

	result := m.Feed{}
	err = result.Decode(data)
	return &result, err

}

// SaveFeeds save feeds
func (z *Database) SaveFeeds(feeds *m.Feeds) error {

	for _, f := range feeds.Values {

		// get old record
		f0, err := z.GetFeedByID(f.ID)
		if err != nil {
			return err
		}

		// save new record
		err = z.saveFeed(f, f0)
		if err != nil {
			return err
		}

	} // loop

	return nil

}

func (z *Database) saveFeed(f *m.Feed, f0 *m.Feed) error {

	err := z.db.Update(func(tx *bolt.Tx) error {

		data, err := f.Encode()
		if err != nil {
			return err
		}

		b := tx.Bucket([]byte(bucketFeed))
		indexes := tx.Bucket([]byte(bucketIndex))
		idxFeedByURL := indexes.Bucket([]byte(bucketIndexFeedByURL))
		idxNextFetch := indexes.Bucket([]byte(bucketIndexNextFetch))

		// save record
		if err = b.Put([]byte(f.ID), data); err != nil {
			return err
		}

		// remove old index entries
		if f0 != nil {

			if err := idxFeedByURL.Delete([]byte(f0.URL)); err != nil {
				return err
			}
			if err := idxNextFetch.Delete([]byte(fetchKey(f0))); err != nil {
				return err
			}

		}

		// add index entries
		if err := idxFeedByURL.Put([]byte(f.URL), []byte(f.ID)); err != nil {
			return err
		}
		if err := idxNextFetch.Put([]byte(fetchKey(f)), []byte(f.ID)); err != nil {
			return err
		}

		return nil

	})

	if err == nil {
		z.checkIndexForEntries(bucketIndexNextFetch, f.ID, 1)
	} else {
		logger.Println("Cannot check for duplicates, error")
	}

	return err

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
