package model

import (
	"encoding/json"
	"sort"
	"strings"
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

// AdjustFetchTime sets the FetchTime to interval units in the future.
func (z *Feed) AdjustFetchTime(interval time.Duration) {
	z.NextFetch = time.Now().Add(interval).Truncate(time.Second)
}

// GetID returns the unique ID for the object
func (z *Feed) GetID() string {
	return z.ID
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

func (z *Feed) clear() {
	z.ID = empty
	z.URL = empty
	z.SiteURL = empty
	z.ETag = empty
	z.LastModified = time.Time{}
	z.LastUpdated = time.Time{}
	z.NextFetch = time.Time{}
	z.Notes = empty
	z.Title = empty
	z.Status = empty
	z.StatusMessage = empty
	z.StatusSince = time.Time{}
}

func (z *Feed) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Feed) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Feed) hasIncrementingID() bool {
	return true
}

func (z *Feed) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexFeedNextFetch] = []string{keyEncodeTime(z.NextFetch), z.ID}
	result[indexFeedURL] = []string{strings.ToLower(z.URL)}
	return result
}

func (z *Feed) setID(tx Transaction) error {
	id, err := tx.NextID(entityFeed)
	if err != nil {
		return err
	}
	z.ID = keyEncodeUint(id)
	return nil
}

// Feeds is a collection of Feed elements
type Feeds []*Feed

func (z Feeds) Len() int      { return len(z) }
func (z Feeds) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Feeds) Less(i, j int) bool {
	return z[i].ID < z[j].ID
}

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

// ByURLAll groups elements in the Feeds collection by lowercase URL
func (z Feeds) ByURLAll() map[string]Feeds {
	result := make(map[string]Feeds)
	for _, feed := range z {
		url := strings.ToLower(feed.URL)
		feeds := result[url]
		feeds = append(feeds, feed)
		result[url] = feeds
	}
	return result
}

// SortByID sort collection by ID
func (z Feeds) SortByID() {
	sort.Stable(z)
}

func (z *Feeds) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Feeds) encode() ([]byte, error) {
	return json.Marshal(z)
}
