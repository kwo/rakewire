package model

//go:generate gokv $GOFILE

import (
	"time"
)

// Feed feed descriptor
type Feed struct {
	Attempt *FeedLog `json:"-" kv:"-"`
	Entries []*Entry `json:"-" kv:"-"`

	ID      uint64 `json:"id"  kv:"NextFetch:2"`
	URL     string `json:"url" kv:"URL:1:lower"`
	SiteURL string `json:"siteURL"`

	ETag         string    `json:"etag"`
	LastModified time.Time `json:"lastModified"`

	LastUpdated time.Time `json:"lastUpdated"`
	NextFetch   time.Time `json:"nextFetch" kv:"NextFetch:1"`

	Notes string `json:"notes,omitempty"`
	Title string `json:"title"`

	Status        string    `json:"status"`
	StatusMessage string    `json:"statusMessage"`
	StatusSince   time.Time `json:"statusSince"` // time of last status
}

// NewFeed instantiate a new Feed object
func NewFeed(url string) *Feed {
	return &Feed{
		URL:       url,
		NextFetch: time.Now().Truncate(time.Second),
	}
}

// AddEntry to the feed
func (z *Feed) AddEntry(guID string) *Entry {
	entry := NewEntry(z.ID, guID)
	z.Entries = append(z.Entries, entry)
	return entry
}

// UpdateFetchTime increases the fetch interval
func (z *Feed) UpdateFetchTime(lastUpdated time.Time) {

	now := time.Now()

	bumpFetchTime :=
		func(interval time.Duration) {
			min := now.Add(1 * time.Minute)
			result := lastUpdated
			for result.Before(min) {
				result = result.Add(interval)
			}
			z.NextFetch = result.Truncate(time.Second)
		}

	d := now.Sub(lastUpdated) // how long since the last update?

	switch {
	case d < 30*time.Minute:
		bumpFetchTime(10 * time.Minute)
	case d > 72*time.Hour:
		bumpFetchTime(24 * time.Hour)
	case true:
		bumpFetchTime(1 * time.Hour)
	}

}

// AdjustFetchTime sets the FetchTime to interval units in the future.
func (z *Feed) AdjustFetchTime(interval time.Duration) {
	z.NextFetch = time.Now().Add(interval).Truncate(time.Second)
}
