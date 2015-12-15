package bolt

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"strconv"
	"time"
)

// GetFeeds list feeds
func (z *Service) GetFeeds() ([]*m.Feed, error) {

	// define index keys
	minKeys := []string{chMin}
	nxtKeys := []string{chMax}

	var result []*m.Feed

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.FeedEntity)).Bucket([]byte(m.FeedIndexURL))
		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.FeedEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			id, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			if data, ok := kvGet(id, b); ok {
				f := &m.Feed{}
				if err := f.Deserialize(data); err != nil {
					return err
				}
				result = append(result, f)
			}

		}

		return nil

	})

	return result, err

}

// GetFetchFeeds get feeds to be fetched within the given max time parameter.
func (z *Service) GetFetchFeeds(maxTime time.Time) ([]*m.Feed, error) {

	// define index keys
	if maxTime.IsZero() {
		maxTime = time.Now()
	}
	f := &m.Feed{}
	f.NextFetch = maxTime
	minKeys := []string{chMin}
	nxtKeys := f.IndexKeys()[m.FeedIndexNextFetch]

	var result []*m.Feed

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.FeedEntity)).Bucket([]byte(m.FeedIndexNextFetch))
		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.FeedEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			id, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			if data, ok := kvGet(id, b); ok {
				f := &m.Feed{}
				if err := f.Deserialize(data); err != nil {
					return err
				}
				result = append(result, f)
			}

		}

		return nil

	})

	return result, err

}

// GetFeedByID return feed given id
func (z *Service) GetFeedByID(id uint64) (*m.Feed, error) {

	found := false
	result := &m.Feed{}

	err := z.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.FeedEntity))
		if data, ok := kvGet(id, b); ok {
			found = true
			return result.Deserialize(data)
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else if !found {
		return nil, nil
	}
	return result, nil

}

// GetFeedByURL return feed given url
func (z *Service) GetFeedByURL(url string) (*m.Feed, error) {

	found := false
	result := &m.Feed{}

	err := z.db.View(func(tx *bolt.Tx) error {
		if data, ok := kvGetFromIndex(m.FeedEntity, m.FeedIndexURL, []string{url}, tx); ok {
			found = true
			return result.Deserialize(data)
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else if !found {
		return nil, nil
	}
	return result, nil

}

// SaveFeed save feeds
func (z *Service) SaveFeed(feed *m.Feed) error {

	if feed == nil {
		return fmt.Errorf("Nil feed")
	}

	z.Lock()
	err := z.db.Update(func(tx *bolt.Tx) error {

		// save feed log if available
		if feed.Attempt != nil {
			if err := kvSave(feed.Attempt, tx); err != nil {
				return err
			}
		}

		// save entries
		if feed.Entries != nil {
			for _, entry := range feed.Entries {
				if err := kvSave(entry, tx); err != nil {
					return err
				}
			}
		}

		// save feed itself
		if err := kvSave(feed, tx); err != nil {
			return err
		}

		return nil

	})
	z.Unlock()

	return err

}
