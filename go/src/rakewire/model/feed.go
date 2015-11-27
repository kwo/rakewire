package model

import (
	"github.com/pborman/uuid"
	"rakewire/feedparser"
	"time"
)

// Feed feed descriptor
type Feed struct {
	// Current fetch attempt for feed
	Attempt *FeedLog `kv:"-" json:"-"`
	// ETag HTTP header on the last 200
	ETag string `json:"etag"`
	// Feed object parsed from Body
	Feed *feedparser.Feed `kv:"-" json:"-"`
	// UUID
	ID string `kv:"key" json:"id"`
	// LastModified HTTP header on the last 200
	LastModified time.Time `json:"lastModified"`
	// Time the feed was last updated
	LastUpdated time.Time `json:"lastUpdated"`
	// Past fetch attempts for feed
	Log []*FeedLog `kv:"-" json:"-"`
	// Time of next scheduled fetch
	NextFetch time.Time `kv:"NextFetch:1" json:"nextFetch"`
	// User notes of the feed
	Notes string `json:"notes,omitempty"`
	// User defined title of the feed
	Title string `json:"title"`
	// URL updated if feed is permenently redirected
	URL string `kv:"URL:1" json:"url"`
}

// AddLog adds a feedlog to the Feed returning its ID
func (z *Feed) AddLog(feedlog *FeedLog) string {
	z.Log = append(z.Log, feedlog)
	return feedlog.ID
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
