package bolt

import (
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"time"
)

// GetFeedLog retrieves the past fetch attempts for the feed in reverse chronological order.
// If since is equal to 0, return all.
func (z *Service) GetFeedLog(feedID string, since time.Duration) ([]*m.FeedLog, error) {

	fl := &m.FeedLog{}
	fl.FeedID = feedID
	fl.StartTime = time.Now().Truncate(time.Second).Add(-since)
	minKeys := fl.IndexKeys()[m.FeedLogIndexFeedTime]
	nxtKeys := []string{chMax}

	var result []*m.FeedLog

	err := z.db.View(func(tx *bolt.Tx) error {

		maps, err := kvQuery(m.FeedLogEntity, m.FeedLogIndexFeedTime, minKeys, nxtKeys, tx)
		if err != nil {
			return err
		}

		for _, data := range maps {
			fl := &m.FeedLog{}
			err = fl.Deserialize(data)
			if err != nil {
				return err
			}
			result = append(result, fl)
		}

		return nil

	})

	// reverse order of result
	for left, right := 0, len(result)-1; left < right; left, right = left+1, right-1 {
		result[left], result[right] = result[right], result[left]
	}

	return result, err

}
