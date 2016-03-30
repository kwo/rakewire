package model

//go:generate gokv $GOFILE

import (
	"time"
)

const (
	entityFeed         = "Feed"
	indexFeedNextFetch = "NextFetch"
	indexFeedURL       = "URL"
)

var (
	indexesFeed = []string{
		indexFeedNextFetch, indexFeedURL,
	}
)

// Feeds is a collection of Feed elements
type Feeds []*Feed

// ByID groups elements in the Feeds collection by ID
func (z Feeds) ByID() map[string]*Feed {
	result := make(map[string]*Feed)
	for _, feed := range z {
		result[feed.ID] = feed
	}
	return result
}

// ByURL maps feeds keyed by URL.
// If multiple feeds exist with the same URL, the last feed will be keyed and the others ignored.
// See also ByURLAll
func (z Feeds) ByURL() map[string]*Feed {
	result := make(map[string]*Feed)
	for _, feed := range z {
		result[feed.URL] = feed
	}
	return result
}

// ByURLAll groups elements in the Feeds collection by URL
func (z Feeds) ByURLAll() map[string]Feeds {
	result := make(map[string]Feeds)
	for _, feed := range z {
		feeds := result[feed.URL]
		feeds = append(feeds, feed)
		result[feed.URL] = feeds
	}
	return result
}

// Feed feed descriptor
type Feed struct {
	ID            string    `json:"id"`
	URL           string    `json:"url"`
	SiteURL       string    `json:"siteURL,omitempty"`
	ETag          string    `json:"etag,omitempty"`
	LastModified  time.Time `json:"lastModified,omitempty"`
	LastUpdated   time.Time `json:"lastUpdated,omitempty"`
	NextFetch     time.Time `json:"nextFetch,omitempty"`
	Notes         string    `json:"notes,omitempty"`
	Title         string    `json:"title,omitempty"`
	Status        string    `json:"status,omitempty"`
	StatusMessage string    `json:"statusMessage,omitempty"`
	StatusSince   time.Time `json:"statusSince,omitempty"` // time of last status
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
