package db

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"time"
)

// FeedInfo feed descriptior
type FeedInfo struct {
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

// NewFeedInfo instantiate a new FeedInfo object with a new UUID
func NewFeedInfo(url string) (*FeedInfo, error) {
	x := FeedInfo{
		ID:  uuid.NewUUID().String(),
		URL: url,
	}
	return &x, nil
}

// Marshal serialize FeedInfo object to bytes
func (z *FeedInfo) Marshal() ([]byte, error) {
	return json.Marshal(z)
}

// Unmarshal serialize FeedInfo object to bytes
func (z *FeedInfo) Unmarshal(data []byte) error {
	return json.Unmarshal(data, z)
}
