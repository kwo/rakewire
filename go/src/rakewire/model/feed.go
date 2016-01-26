package model

//go:generate gokv $GOFILE

import (
	"time"
)

// Feed feed descriptor
type Feed struct {
	Transmission *Transmission `json:"-" kv:"-"`
	Items        []*Item       `json:"-" kv:"-"`

	ID      uint64 `json:"id"  kv:"NextFetch:2"`
	URL     string `json:"url" kv:"+required,URL:1:lower"`
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

// AddItem to the feed
func (z *Feed) AddItem(guID string) *Item {
	item := NewItem(z.ID, guID)
	z.Items = append(z.Items, item)
	return item
}

// UpdateFetchTime increases the fetch interval
func (z *Feed) UpdateFetchTime(lastUpdated time.Time) {

	now := time.Now()

	bumpFetchTime :=
		func(interval time.Duration) {
			min := now.Add(5 * time.Minute)
			result := lastUpdated
			for result.Before(min) {
				result = result.Add(interval)
			}
			z.NextFetch = result.Truncate(time.Second)
		}

	d := now.Sub(lastUpdated) // how long since the last update?

	switch {
	case d < 2*time.Hour:
		bumpFetchTime(15 * time.Minute)
	case true:
		bumpFetchTime(1 * time.Hour)
	}

}

// AdjustFetchTime sets the FetchTime to interval units in the future.
func (z *Feed) AdjustFetchTime(interval time.Duration) {
	z.NextFetch = time.Now().Add(interval).Truncate(time.Second)
}
