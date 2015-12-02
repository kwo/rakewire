package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"time"
)

// GetFeeds list feeds
func (z *Database) GetFeeds() ([]*m.Feed, error) {

	result := []*m.Feed{}
	add := func() interface{} {
		f := m.NewFeed("")
		result = append(result, f)
		return f
	}

	err := z.db.View(func(tx *bolt.Tx) error {
		return Query("Feed", empty, nil, []interface{}{chMax}, add, tx)
	})

	return result, err

}

// GetFetchFeeds get feeds to be fetched within the given max time parameter.
func (z *Database) GetFetchFeeds(maxTime time.Time) ([]*m.Feed, error) {

	if maxTime.IsZero() {
		maxTime = time.Now()
	}

	result := []*m.Feed{}
	add := func() interface{} {
		f := m.NewFeed("")
		result = append(result, f)
		return f
	}

	err := z.db.View(func(tx *bolt.Tx) error {
		return Query("Feed", "NextFetch", nil, []interface{}{maxTime}, add, tx)
	})

	return result, err

}

// GetFeedByID return feed given UUID
func (z *Database) GetFeedByID(id string) (*m.Feed, error) {

	result := m.NewFeed("")
	result.ID = id

	err := z.db.View(func(tx *bolt.Tx) error {
		return Get(result, tx)
	})
	if err != nil {
		return nil, err
	}

	if result != nil && result.ID != id {
		return nil, nil
	}

	return result, nil

}

// GetFeedByURL return feed given url
func (z *Database) GetFeedByURL(url string) (*m.Feed, error) {

	feeds := []*m.Feed{}
	add := func() interface{} {
		f := m.NewFeed(url)
		feeds = append(feeds, f)
		return f
	}

	err := z.db.View(func(tx *bolt.Tx) error {
		return Query("Feed", "URL", []interface{}{url}, []interface{}{url}, add, tx)
	})
	if err != nil {
		return nil, err
	}

	if len(feeds) == 0 {
		return nil, nil
	} else if len(feeds) > 1 {
		return nil, fmt.Errorf("Unique index returned multiple results: %s, URL: %s", "Feed/URL", url)
	}

	return feeds[0], nil

}

// SaveFeed save feeds
func (z *Database) SaveFeed(feed *m.Feed) error {

	if feed == nil {
		return fmt.Errorf("Nil feed")
	}

	z.Lock()
	err := z.db.Update(func(tx *bolt.Tx) error {

		// save feed log if available
		if feed.Attempt != nil {
			if _, err := Put(feed.Attempt, tx); err != nil {
				return err
			}
		}

		// save entries // TODO: recalculate update times
		if feed.Entries != nil {
			for _, entry := range feed.Entries {
				if _, err := Put(entry, tx); err != nil {
					return err
				}
			}
		}

		// save feed itself
		if _, err := Put(feed, tx); err != nil {
			return err
		}

		return nil

	})
	z.Unlock()

	return err

}
