package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
	m "rakewire.com/model"
	"time"
)

// GetFeedLog retrieves the past fetch attempts for the feed in reverse chronological order.
// If since is equal to 0, return all.
func (z *Database) GetFeedLog(id string, since time.Duration) ([]*m.FeedLog, error) {

	var min []byte
	var max []byte
	maxDate := time.Now()
	if since == 0 {
		min = []byte(formatFeedLogKey(id, nil))
		max = []byte(formatFeedLogKey(id, &maxDate))
	} else {
		minDate := maxDate.Add(-since)
		min = []byte(formatFeedLogKey(id, &minDate))
		max = []byte(formatFeedLogKey(id, &maxDate))
	}

	var result []*m.FeedLog

	err := z.db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucketFeedLog))
		c := b.Cursor()

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			entry := &m.FeedLog{}
			if err := entry.Decode(v); err != nil {
				return err
			}
			result = append([]*m.FeedLog{entry}, result...)
		} // for

		return nil

	})

	return result, err

}

func (z *Database) addFeedLog(tx *bolt.Tx, entry *m.FeedLog) error {

	data, err := entry.Encode()
	if err != nil {
		return err
	}

	b := tx.Bucket([]byte(bucketFeedLog))

	// save record
	if err = b.Put([]byte(formatFeedLogKey(entry.FeedID, entry.StartTime)), data); err != nil {
		return err
	}

	return nil

}
