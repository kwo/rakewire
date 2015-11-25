package bolt

import (
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"time"
)

// GetFeedLog retrieves the past fetch attempts for the feed in reverse chronological order.
// If since is equal to 0, return all.
func (z *Database) GetFeedLog(feedID string, since time.Duration) ([]*m.FeedLog, error) {

	maxDate := time.Now()
	minDate := maxDate.Add(-since)

	result := []*m.FeedLog{}
	add := func() interface{} {
		fl := &m.FeedLog{}
		result = append(result, fl)
		return fl
	}

	err := z.db.View(func(tx *bolt.Tx) error {
		return Query("FeedLog", "FeedTime", []interface{}{feedID, minDate}, []interface{}{feedID, maxDate}, add, tx)
	})

	// reverse order of result
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}

	return result, err

}
