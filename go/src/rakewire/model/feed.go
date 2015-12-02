package model

import (
	"github.com/pborman/uuid"
	"time"
)

// Feed feed descriptor
type Feed struct {
	Attempt *FeedLog `kv:"-" json:"-"`
	Entries []*Entry `kv:"-" json:"-"`

	ID  string `kv:"key" json:"id"`
	URL string `kv:"URL:1" json:"url"`

	ETag         string    `json:"etag"`
	LastModified time.Time `json:"lastModified"`

	LastUpdated time.Time `json:"lastUpdated"`
	NextFetch   time.Time `kv:"NextFetch:1" json:"nextFetch"`

	Notes string `json:"notes,omitempty"`
	Title string `json:"title"`

	Status        string    `json:"status"`
	StatusMessage string    `json:"statusMessage"`
	StatusSince   time.Time `json:"statusSince"` // time of last status
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
