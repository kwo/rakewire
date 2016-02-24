package model

import (
	"time"
)

// TransmissionsByFeed retrieves the past fetch attempts for the feed in reverse chronological order.
// If since is equal to 0, return all.
func TransmissionsByFeed(feedID string, since time.Duration, tx Transaction) (Transmissions, error) {

	transmissions := Transmissions{}

	// transmission index FeedTime = FeedID|StartTime : TransmissionID
	now := time.Now().Truncate(time.Second)
	min := kvKeyEncode(feedID, kvKeyTimeEncode(now.Add(-since)))
	max := kvKeyEncode(feedID, kvKeyTimeEncode(now))
	bIndex := tx.Bucket(bucketIndex, transmissionEntity, transmissionIndexFeedTime)
	bTransmission := tx.Bucket(bucketData, transmissionEntity)

	err := bIndex.IterateIndex(bTransmission, min, max, func(record Record) error {
		transmission := &Transmission{}
		if err := transmission.deserialize(record); err != nil {
			return err
		}
		transmissions = append(transmissions, transmission)
		return nil
	})

	// reverse order of result
	transmissions.Reverse()

	return transmissions, err

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
