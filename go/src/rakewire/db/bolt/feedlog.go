package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"strconv"
	"time"
)

// GetFeedLog retrieves the past fetch attempts for the feed in reverse chronological order.
// If since is equal to 0, return all.
func (z *Service) GetFeedLog(feedID uint64, since time.Duration) ([]*m.FeedLog, error) {

	var result []*m.FeedLog

	// define index keys
	now := time.Now().Truncate(time.Second)
	fl := &m.FeedLog{}
	fl.FeedID = feedID
	fl.StartTime = now.Add(-since)
	minKeys := fl.IndexKeys()[m.FeedLogIndexFeedTime]
	fl.StartTime = now.Add(1 * time.Minute) // max later than now
	nxtKeys := fl.IndexKeys()[m.FeedLogIndexFeedTime]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.FeedLogEntity)).Bucket([]byte(m.FeedLogIndexFeedTime))
		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.FeedLogEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			id, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			if data, ok := kvGet(id, b); ok {
				fl := &m.FeedLog{}
				if err := fl.Deserialize(data); err != nil {
					return err
				}
				result = append(result, fl)
			}

		}

		return nil

	})

	// reverse order of result
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}

	return result, err

}

// FeedLogGetLastFetchTime retrieves the most recent fetch activity
func (z *Service) FeedLogGetLastFetchTime() (time.Time, error) {

	startTime := time.Now().Truncate(time.Second)
	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.FeedLogEntity)).Bucket([]byte(m.FeedLogIndexTime))
		c := bIndex.Cursor()
		k, _ := c.Last()
		if k != nil {
			startTimeStr := kvKeyElement(k, 0)
			t, err := time.Parse(time.RFC3339, startTimeStr)
			if err != nil {
				return err
			}
			startTime = t
		}

		return nil

	})

	return startTime, err

}
