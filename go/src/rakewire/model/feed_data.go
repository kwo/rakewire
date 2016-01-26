package model

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// FeedsAll list feeds
func FeedsAll(tx Transaction) ([]*Feed, error) {

	var result []*Feed

	bIndex := tx.Bucket(bucketIndex).Bucket(FeedEntity).Bucket(FeedIndexURL)
	b := tx.Bucket(bucketData).Bucket(FeedEntity)

	c := bIndex.Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, b); ok {
			f := &Feed{}
			if err := f.Deserialize(data); err != nil {
				return nil, err
			}
			result = append(result, f)
		}

	}

	return result, nil

}

// FeedsFetch get feeds to be fetched within the given max time parameter.
func FeedsFetch(maxTime time.Time, tx Transaction) ([]*Feed, error) {

	// define index keys
	if maxTime.IsZero() {
		maxTime = time.Now()
	}
	f := &Feed{}
	f.NextFetch = maxTime
	nxtKeys := f.IndexKeys()[FeedIndexNextFetch]

	var result []*Feed

	bIndex := tx.Bucket(bucketIndex).Bucket(FeedEntity).Bucket(FeedIndexNextFetch)
	b := tx.Bucket(bucketData).Bucket(FeedEntity)

	c := bIndex.Cursor()
	nxt := []byte(kvKeys(nxtKeys))
	for k, v := c.First(); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, b); ok {
			f := &Feed{}
			if err := f.Deserialize(data); err != nil {
				return nil, err
			}
			result = append(result, f)
		}

	}

	return result, nil

}

// FeedByID return feed given id
func FeedByID(id uint64, tx Transaction) (feed *Feed, err error) {
	b := tx.Bucket(bucketData).Bucket(FeedEntity)
	if data, ok := kvGet(id, b); ok {
		feed = &Feed{}
		err = feed.Deserialize(data)
	} else {
		err = fmt.Errorf("Feed not found: %d", id)
	}
	return
}

// FeedByURL return feed given url
func FeedByURL(url string, tx Transaction) (feed *Feed, err error) {
	if data, ok := kvGetFromIndex(FeedEntity, FeedIndexURL, []string{strings.ToLower(url)}, tx); ok {
		feed = &Feed{}
		err = feed.Deserialize(data)
	} else {
		err = fmt.Errorf("Feed not found: %s", url)
	}
	return
}

// Save save feeds
func (feed *Feed) Save(tx Transaction) ([]*Item, error) {

	if feed == nil {
		return nil, fmt.Errorf("Nil feed")
	}

	newItems := []*Item{}

	// save feed log if available
	if feed.Attempt != nil {
		if err := kvSave(FeedLogEntity, feed.Attempt, tx); err != nil {
			return nil, err
		}
	}

	// save items
	if feed.Items != nil {
		for _, item := range feed.Items {
			if item.ID == 0 {
				newItems = append(newItems, item)
			}
			if err := kvSave(ItemEntity, item, tx); err != nil {
				return nil, err
			}
		}
	}

	if err := EntriesAddNew(newItems, tx); err != nil {
		return nil, err
	}

	// save feed itself
	if err := kvSave(FeedEntity, feed, tx); err != nil {
		return nil, err
	}

	return newItems, nil

}

// Delete removes a feed and associated items from the database.
func (feed *Feed) Delete(tx Transaction) error {

	// remove items
	items, err := ItemsByFeed(feed.ID, tx)
	if err != nil {
		return err
	}
	for _, item := range items {
		if err := kvDelete(ItemEntity, item, tx); err != nil {
			return err
		}
	}

	// remove feedlogs
	feedlogs, err := FeedLogsByFeed(feed.ID, time.Now().Sub(time.Time{}), tx)
	if err != nil {
		return err
	}
	for _, feedlog := range feedlogs {
		if err := kvDelete(FeedLogEntity, feedlog, tx); err != nil {
			return err
		}
	}

	// remove feed itself
	return kvDelete(FeedEntity, feed, tx)

}

// FeedDuplicates finds duplicate feeds keyed by original feed.
func FeedDuplicates(tx Transaction) (map[string][]uint64, error) {

	result := make(map[string][]uint64)

	b := tx.Bucket(bucketData).Bucket(FeedEntity)
	err := b.ForEach(func(k, v []byte) error {
		if fieldName := kvKeyElement(k, 1); fieldName == "URL" {
			id, err := kvKeyElementID(k, 0)
			if err != nil {
				return err
			}
			url := string(v)
			result[url] = append(result[url], id)
		}
		return nil
	})

	return result, err

}
