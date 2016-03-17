package model

import (
	"bytes"
	"strings"
	"time"
)

// F groups all feed database methods
var F = &feedStore{}

type feedStore struct{}

func (z *feedStore) Delete(tx Transaction, id string) error {
	return delete(tx, entityFeed, id)
}

func (z *feedStore) Get(tx Transaction, id string) *Feed {
	bData := tx.Bucket(bucketData, entityFeed)
	if data := bData.Get([]byte(id)); data != nil {
		feed := &Feed{}
		if err := feed.decode(data); err == nil {
			return feed
		}
	}
	return nil
}

func (z *feedStore) GetByURL(tx Transaction, url string) *Feed {
	// index Feed URL = URL (lowercase) : FeedID
	b := tx.Bucket(bucketIndex, entityFeed, indexFeedURL)
	if id := b.Get([]byte(strings.ToLower(url))); id != nil {
		return z.Get(tx, string(id))
	}
	return nil
}

// GetNext returns all feeds which are due to be fetched within the given max time.
func (z *feedStore) GetNext(tx Transaction, maxTime time.Time) Feeds {
	// index Feed NextFetch = FetchTime|FeedID : FeedID
	feeds := Feeds{}
	nxtTime := maxTime.Add(1 * time.Second).Truncate(time.Second)
	nxt := []byte(keyEncodeTime(nxtTime))
	b := tx.Bucket(bucketIndex, entityFeed, indexFeedNextFetch)
	c := b.Cursor()
	for k, v := c.First(); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
		feedID := string(v)
		if feed := z.Get(tx, feedID); feed != nil {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func (z *feedStore) New(url string) *Feed {
	return &Feed{
		URL:       url,
		NextFetch: time.Now().Truncate(time.Second),
	}
}

func (z *feedStore) Save(tx Transaction, feed *Feed) error {
	return save(tx, entityFeed, feed)
}
