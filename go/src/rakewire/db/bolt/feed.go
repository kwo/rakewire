package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"time"
)

// GetFeeds list feeds
func (z *Database) GetFeeds() ([]*m.Feed, error) {

	var result []*m.Feed

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketFeed))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {

			f := m.Feed{}
			if err := f.Decode(v); err != nil {
				return err
			}
			if f.ID != string(k) {
				logger.Warnf("ID/key mismatch: %s/%s\n", k, f.ID)
			} else {
				result = append(result, &f)
			}

		} // for

		return nil

	})

	return result, err

}

// GetFetchFeeds get feeds to be fetched within the given max time parameter.
func (z *Database) GetFetchFeeds(maxTime *time.Time) ([]*m.Feed, error) {

	var max []byte
	if maxTime == nil {
		max = []byte(formatMaxTime(time.Now()))
	} else {
		max = []byte(formatMaxTime(*maxTime))
	}

	var result []*m.Feed

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketFeed))
		idxNextFetch := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(bucketIndexNextFetch))
		c := idxNextFetch.Cursor()

		//logger.Debugf("max: %s\n", string(max))
		for k, uuid := c.First(); k != nil && bytes.Compare(k, max) <= 0; k, uuid = c.Next() {

			//logger.Debugf("key: %s: %s", k, uuid)

			v := b.Get(uuid)
			f := &m.Feed{}
			if err := f.Decode(v); err != nil {
				return err
			}
			result = append(result, f)

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

// SaveFeed save feeds
func (z *Database) SaveFeed(fNew *m.Feed) error {

	// get old record
	fOld, err := z.GetFeedByID(fNew.ID)
	if err != nil {
		return err
	}

	// save new record
	err = z.saveFeed(fNew, fOld)
	if err != nil {
		return err
	}

	return nil

}

func (z *Database) saveFeed(fNew *m.Feed, fOld *m.Feed) error {

	err := z.db.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucketFeed))
		indexes := tx.Bucket([]byte(bucketIndex))
		idxFeedByURL := indexes.Bucket([]byte(bucketIndexFeedByURL))
		idxNextFetch := indexes.Bucket([]byte(bucketIndexNextFetch))

		// log record
		if fNew.Attempt != nil {
			fNew.Last = fNew.Attempt
			if err := z.addFeedLog(tx, fNew.ID, fNew.Last); err != nil {
				return err
			}
			if fNew.Last.HTTP.StatusCode == 200 {
				fNew.Last200 = fNew.Last
			}
		}

		// encode record (must be after feed log)
		data, err := fNew.Encode()
		if err != nil {
			return err
		}

		// save record
		if err = b.Put([]byte(fNew.ID), data); err != nil {
			return err
		}

		// remove old index entries
		if fOld != nil {

			if err := idxFeedByURL.Delete([]byte(fOld.URL)); err != nil {
				return err
			}
			if err := idxNextFetch.Delete([]byte(fetchKey(fOld))); err != nil {
				return err
			}

		}

		// add index entries
		if err := idxFeedByURL.Put([]byte(fNew.URL), []byte(fNew.ID)); err != nil {
			return err
		}
		if err := idxNextFetch.Put([]byte(fetchKey(fNew)), []byte(fNew.ID)); err != nil {
			return err
		}

		return nil

	})

	if err == nil {
		z.checkIndexForEntries(bucketIndexNextFetch, fNew.ID, 1)
	} else {
		logger.Warnf("Cannot check for duplicates: %s", err.Error())
	}

	return err

}
