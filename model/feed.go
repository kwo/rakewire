package model

import (
	"bufio"
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"io"
	"strings"
	"time"
)

const (
	// FeedIntervalMin is the minimum feed fetch interval (2m48s750ms) x 2^9 = 1 day
	FeedIntervalMin time.Duration = time.Millisecond * 168750
	// FeedIntervalMax is the maximum feed fetch interval
	FeedIntervalMax time.Duration = time.Hour * 24
)

// Feeds collection of Feed
type Feeds struct {
	Values []*Feed
	Index  map[string]*Feed
}

// Feed feed descriptior
// Also a super type of fetch.Request and fetch.Response
type Feed struct {
	// Body is the HTTP payload
	Body []byte `json:"-"`
	// Checksum of HTTP payload (independent of etag)
	Checksum string `json:"checksum,omitempty"`
	// Etag from HTTP Request - used for conditional GETs
	ETag string `json:"etag,omitempty"`
	// Type of feed: Atom, RSS2, etc.
	Flavor string `json:"flavor,omitempty"`
	// how often to poll the feed in minutes
	Interval time.Duration `json:"interval"`
	// Feed generator
	Generator string `json:"generator,omitempty"`
	// Hub URL
	Hub string `json:"hub,omitempty"`
	// Feed icon
	Icon string `json:"icon,omitempty"`
	// UUID
	ID string `json:"id"`
	// Time of last fetch attempt
	LastAttempt *time.Time `json:"lastAttempt,omitempty"`
	// Time of last successful fetch completion
	LastFetch *time.Time `json:"lastFetch"`
	// Last-Modified time from HTTP Request - used for conditional GETs
	LastModified *time.Time `json:"lastModified,omitempty"`
	// Time the feed was last updated (from feed)
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Last HTTP status code
	StatusCode int `json:"statusCode,omitempty"`
	// Feed title
	Title string `json:"title"`
	// URL updated if feed is permenently redirected
	URL string `json:"url"`
}

// ========== Feed ==========

// NewFeed instantiate a new Feed object with a new UUID
func NewFeed(url string) *Feed {
	lastFetch := time.Now().Add(-24 * time.Hour).Truncate(time.Second)
	x := Feed{
		Interval:  FeedIntervalMin,
		ID:        uuid.NewUUID().String(),
		LastFetch: &lastFetch,
		URL:       url,
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

// GetNextFetchTime get the next time to poll feed
func (z *Feed) GetNextFetchTime() *time.Time {
	result := z.LastFetch.Add(z.Interval).Truncate(time.Second)
	return &result
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
