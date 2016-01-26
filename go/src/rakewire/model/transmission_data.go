package model

import (
	"bytes"
	"strconv"
	"time"
)

// TransmissionsByFeed retrieves the past fetch attempts for the feed in reverse chronological order.
// If since is equal to 0, return all.
func TransmissionsByFeed(feedID uint64, since time.Duration, tx Transaction) ([]*Transmission, error) {

	transmissions := []*Transmission{}

	// define index keys
	now := time.Now().Truncate(time.Second)
	fl := &Transmission{}
	fl.FeedID = feedID
	fl.StartTime = now.Add(-since)
	minKeys := fl.IndexKeys()[TransmissionIndexFeedTime]
	fl.StartTime = now.Add(1 * time.Minute) // max later than now
	nxtKeys := fl.IndexKeys()[TransmissionIndexFeedTime]

	bIndex := tx.Bucket(bucketIndex).Bucket(TransmissionEntity).Bucket(TransmissionIndexFeedTime)
	b := tx.Bucket(bucketData).Bucket(TransmissionEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, b); ok {
			transmission := &Transmission{}
			if err := transmission.Deserialize(data); err != nil {
				return nil, err
			}
			transmissions = append(transmissions, transmission)
		}

	}

	// reverse order of result
	for left, right := 0, len(transmissions)-1; left < right; left, right = left+1, right-1 {
		transmissions[left], transmissions[right] = transmissions[right], transmissions[left]
	}

	return transmissions, nil

}

// LastFetchTime retrieves the most recent fetch activity
func LastFetchTime(tx Transaction) (lastFetchTime time.Time, err error) {

	lastFetchTime = time.Now().Truncate(time.Second)

	bIndex := tx.Bucket(bucketIndex).Bucket(TransmissionEntity).Bucket(TransmissionIndexTime)
	c := bIndex.Cursor()
	k, _ := c.Last()
	if k != nil {
		startTimeStr := kvKeyElement(k, 0)
		t, err := time.Parse(time.RFC3339, startTimeStr)
		if err == nil {
			lastFetchTime = t
		}
	}

	return

}
