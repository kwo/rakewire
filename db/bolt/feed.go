package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
	m "rakewire.com/model"
	"time"
)

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
			// #TODO:0 remove hack that feed has .Last element
			if f.Last == nil {
				f.Last = &m.FeedLog{
					FeedID: f.ID,
				}
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
	// #TODO:0 remove hack that feed has .Last element
	if result.Last == nil {
		result.Last = &m.FeedLog{
			FeedID: result.ID,
		}
	}
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
	// #TODO:0 remove hack that feed has .Last element
	if result.Last == nil {
		result.Last = &m.FeedLog{
			FeedID: result.ID,
		}
	}
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

		// log record
		if f.Attempt != nil {
			f.Last = f.Attempt
			f.Attempt = nil
			if err := z.addFeedLog(tx, f.Last); err != nil {
				return err
			}
		}

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
