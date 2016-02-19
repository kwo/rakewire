package model

import (
	"bytes"
	"time"
)

// TransmissionsByFeed retrieves the past fetch attempts for the feed in reverse chronological order.
// If since is equal to 0, return all.
func TransmissionsByFeed(feedID string, since time.Duration, tx Transaction) (Transmissions, error) {

	transmissions := Transmissions{}

	// transmission index FeedTime = FeedID|StartTime : TransmissionID
	now := time.Now().Truncate(time.Second)
	min := kvKeyEncode(feedID, kvKeyTimeEncode(now.Add(-since)))
	max := kvKeyEncode(feedID, kvKeyTimeEncode(now.Add(1*time.Minute))) // max later than now
	bIndex := tx.Bucket(bucketIndex).Bucket(transmissionEntity).Bucket(transmissionIndexFeedTime)
	b := tx.Bucket(bucketData).Bucket(transmissionEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		transmissionID := string(v)

		if data, ok := kvGet(transmissionID, b); ok {
			transmission := &Transmission{}
			if err := transmission.deserialize(data); err != nil {
				return nil, err
			}
			transmissions = append(transmissions, transmission)
		}

	}

	// reverse order of result
	transmissions.Reverse()

	return transmissions, nil

}

// LastFetchTime retrieves the most recent fetch activity
func LastFetchTime(tx Transaction) (lastFetchTime time.Time, err error) {

	// transmission index Time = StartTime|TransmissionID : TransmissionID

	lastFetchTime = time.Now().Truncate(time.Second)

	bIndex := tx.Bucket(bucketIndex).Bucket(transmissionEntity).Bucket(transmissionIndexTime)
	c := bIndex.Cursor()
	k, _ := c.Last()
	if k != nil {
		if t, err := kvKeyTimeDecode(kvKeyDecode(k)[0]); err == nil {
			lastFetchTime = t
		}
	}

	return

}
