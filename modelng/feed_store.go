package modelng

import (
	"bytes"
	"strings"
	"time"
)

// F groups all feed database methods
var F = &feedStore{}

type feedStore struct{}

func (z *feedStore) Delete(id string, tx Transaction) error {
	return delete(entityFeed, id, tx)
}

func (z *feedStore) Get(id string, tx Transaction) *Feed {
	bData := tx.Bucket(bucketData, entityFeed)
	if data := bData.Get(keyEncode(id)); data != nil {
		feed := &Feed{}
		if err := feed.decode(data); err == nil {
			return feed
		}
	}
	return nil
}

func (z *feedStore) GetByURL(url string, tx Transaction) *Feed {
	// index Feed URL = URL (lowercase) : FeedID
	b := tx.Bucket(bucketIndex, entityFeed, indexFeedURL)
	if id := b.Get(keyEncode(strings.ToLower(url))); id != nil {
		return F.Get(string(id), tx)
	}
	return nil
}

// GetNext returns all feeds which are due to be fetched within the given max time.
func (z *feedStore) GetNext(maxTime time.Time, tx Transaction) Feeds {
	// index Feed NextFetch = FetchTime|FeedID : FeedID
	feeds := Feeds{}
	nxtTime := maxTime.Add(1 * time.Second).Truncate(time.Second)
	nxt := keyEncode(keyEncodeTime(nxtTime))
	b := tx.Bucket(bucketIndex, entityFeed, indexFeedNextFetch)
	c := b.Cursor()
	for k, v := c.First(); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
		feedID := string(v)
		if feed := F.Get(feedID, tx); feed != nil {
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

func (z *feedStore) Save(feed *Feed, tx Transaction) error {
	return save(entityFeed, feed, tx)
}
