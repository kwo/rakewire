package model

import (
	"github.com/pborman/uuid"
	"rakewire/feedparser"
	"time"
)

// Feed feed descriptor
type Feed struct {
	// Current fetch attempt for feed
	Attempt *FeedLog `db:"-"`
	// Feed object parsed from Body
	Feed *feedparser.Feed `db:"-"`
	// UUID
	ID string `db:"primary-key"`
	// Time the feed was last updated
	LastUpdated time.Time
	// Last fetch
	Last string
	// Last successful fetch with status code 200
	Last200 string
	// Past fetch attempts for feed
	Log []*FeedLog `db:"-"`
	// Time of next scheduled fetch
	NextFetch time.Time `db:"NextFetch:1"`
	// User notes of the feed
	Notes string
	// User defined title of the feed
	Title string
	// URL updated if feed is permenently redirected
	URL string `db:"URL:1"`
}

// AddLog adds a feedlog to the Feed returning its ID
func (z *Feed) AddLog(feedlog *FeedLog) string {
	z.Log = append(z.Log, feedlog)
	return feedlog.ID
}

// GetLast returns the last feed fetch appempt.
func (z *Feed) GetLast() *FeedLog {
	return z.findLogByID(z.Last)
}

// GetLast200 returns the last successful feed fetch.
func (z *Feed) GetLast200() *FeedLog {
	return z.findLogByID(z.Last200)
}

func (z *Feed) findLogByID(id string) *FeedLog {
	for _, fl := range z.Log {
		if fl.ID == id {
			return fl
		}
	}
	return nil
}

// ========== Feed ==========

// NewFeed instantiate a new Feed object with a new UUID
func NewFeed(url string) *Feed {
	return &Feed{
		ID:        uuid.NewUUID().String(),
		URL:       url,
		NextFetch: time.Now(),
	}
}

// UpdateFetchTime increases the fetch interval
func (z *Feed) UpdateFetchTime(lastUpdated time.Time) {

	if lastUpdated.IsZero() {
		return
	}

	z.LastUpdated = lastUpdated

	d := time.Now().Sub(z.LastUpdated) // how long since the last update?

	switch {
	case d < 30*time.Minute:
		z.AdjustFetchTime(10 * time.Minute)
	case d > 72*time.Hour:
		z.AdjustFetchTime(24 * time.Hour)
	case true:
		z.AdjustFetchTime(1 * time.Hour)
	}

}

// AdjustFetchTime sets the FetchTime to interval units in the future.
func (z *Feed) AdjustFetchTime(interval time.Duration) {
	z.NextFetch = time.Now().Add(interval)
}
