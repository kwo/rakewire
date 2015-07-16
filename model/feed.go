package model

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"io"
	"time"
)

// Feeds collection of Feed
type Feeds struct {
	Values []*Feed
	Index  map[string]*Feed
}

// Feed feed descriptior
type Feed struct {
	// Etag from HTTP Request - used for conditional GETs
	ETag string `json:"etag,omitempty"`
	// Type of feed: Atom, RSS2, etc.
	Flavor string `json:"flavor,omitempty"`
	// Not yet in use: how often to poll the feed
	Frequency int `json:"frequency,omitempty"`
	// Feed generator
	Generator string `json:"generator,omitempty"`
	// Hub URL
	Hub string `json:"hub,omitempty"`
	// Feed icon
	Icon string `json:"icon,omitempty"`
	// UUID
	ID string `json:"id"`
	// Time of last fetch attempt (LastStatus is true) or completion (LastStatus is false)
	LastFetch *time.Time `json:"lastFetch,omitempty"`
	// Last-Modified time from HTTP Request - used for conditional GETs
	LastModified *time.Time `json:"lastModified,omitempty"`
	// The last status, true if error, false if successfully completed
	LastStatus bool `json:"lastStatus,omitempty"`
	// Time the feed was last updated (from feed)
	LastUpdated *time.Time `json:"lastUpdated,omitempty"`
	// Feed title
	Title string `json:"title"`
	// URL updated if feed is permenently redirected
	URL string `json:"url"`
}

// ========== Feed ==========

// NewFeed instantiate a new Feed object with a new UUID
func NewFeed(url string) *Feed {
	x := Feed{
		ID:  uuid.NewUUID().String(),
		URL: url,
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
