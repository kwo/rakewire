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
	return deleteObject(tx, entityFeed, id)
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

func (z *feedStore) GetBySubscriptions(tx Transaction, subscriptions Subscriptions) Feeds {
	result := Feeds{}
	for _, subscription := range subscriptions {
		if feed := z.Get(tx, subscription.FeedID); feed != nil {
			result = append(result, feed)
		}
	}
	return result
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

func (z *feedStore) Range(tx Transaction) Feeds {
	feeds := Feeds{}
	c := tx.Bucket(bucketData, entityFeed).Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		feed := &Feed{}
		if err := feed.decode(v); err == nil {
			feeds = append(feeds, feed)
		}
	}
	return feeds
}

func (z *feedStore) Save(tx Transaction, feed *Feed) error {
	return saveObject(tx, entityFeed, feed)
}
