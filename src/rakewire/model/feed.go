package model

import (
	"bufio"
	"encoding/json"
	"github.com/pborman/uuid"
	"io"
	"rakewire/feedparser"
	"strings"
	"time"
)

// #TODO: eliminate Feeds struct

// Feeds collection of Feed
type Feeds struct {
	Values []*Feed
	Index  map[string]*Feed
}

// Feed feed descriptior
type Feed struct {
	// Current fetch attempt for feed
	Attempt *FeedLog `json:"-"`
	// Feed object parsed from Body
	Feed *feedparser.Feed `json:"-"`
	// Type of feed: Atom, RSS2, etc.
	Flavor string `json:"flavor,omitempty"`
	// Feed generator
	Generator string `json:"generator,omitempty"`
	// Hub URL
	Hub string `json:"hub,omitempty"`
	// Feed icon
	Icon string `json:"icon,omitempty"`
	// UUID
	ID string `json:"id"`
	// Time the feed was last updated
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	// Last fetch
	Last *FeedLog `json:"last"`
	// Last successful fetch with status code 200
	Last200 *FeedLog `json:"last200"`
	// Past fetch attempts for feed
	Log []*FeedLog `json:"-"`
	// Time of next scheduled fetch
	NextFetch time.Time `json:"nextFetch"`
	// Feed title
	Title string `json:"title"`
	// URL updated if feed is permenently redirected
	URL string `json:"url"`
}

// ========== Feed ==========

// NewFeed instantiate a new Feed object with a new UUID
func NewFeed(url string) *Feed {
	nextFetch := time.Now().UTC().Truncate(time.Second)
	id := uuid.NewUUID().String()
	x := Feed{
		ID:        id,
		URL:       url,
		Last:      &FeedLog{},
		Last200:   &FeedLog{},
		NextFetch: nextFetch,
	}
	return &x
}

// Decode Feed object from bytes
func (z *Feed) Decode(data []byte) error {
	return json.Unmarshal(data, z)
}

// Encode Feed object to bytes
func (z *Feed) Encode() ([]byte, error) {
	return json.MarshalIndent(z, "", " ")
}

// UpdateFetchTime increases the fetch interval
func (z *Feed) UpdateFetchTime(lastUpdated time.Time) {

	if !lastUpdated.IsZero() {
		z.LastUpdated = lastUpdated
	}

	now := time.Now().UTC().Truncate(time.Second)
	if z.LastUpdated.IsZero() {
		z.LastUpdated = now
	}

	d := now.Sub(z.LastUpdated) // how long since the last update?

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
	now := time.Now().UTC().Truncate(time.Second)
	nextFetch := now.Add(interval)
	z.NextFetch = nextFetch
}

// ========== Feeds ==========

// NewFeeds instantiate a new Feeds collection
func NewFeeds() *Feeds {
	x := Feeds{}
	x.reindex()
	return &x
}

// Add add a Feed to the collection
func (z *Feeds) Add(fd *Feed) {
	z.Values = append(z.Values, fd)
	z.Index[fd.ID] = fd
}

// Get a Feed by id
func (z *Feeds) Get(id string) *Feed {
	return z.Index[id]
}

// Size numbers of feeds in collection
func (z *Feeds) Size() int {
	return len(z.Values)
}

// Serialize serialize Feed object to bytes
func (z *Feeds) Serialize(w io.Writer) error {
	return json.NewEncoder(w).Encode(&z.Values)
}

// Deserialize serialize Feed object to bytes
func (z *Feeds) Deserialize(r io.Reader) error {
	err := json.NewDecoder(r).Decode(&z.Values)
	if err != nil {
		return err
	}
	z.reindex()
	return nil
}

func (z *Feeds) reindex() {
	z.Index = make(map[string]*Feed)
	for _, d := range z.Values {
		z.Index[d.ID] = d
	}
}

// ParseListToFeeds parse url list to feeds
func ParseListToFeeds(r io.Reader) *Feeds {
	feeds := NewFeeds()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var url = strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			feeds.Add(NewFeed(url))
		}
	}
	return feeds
}
