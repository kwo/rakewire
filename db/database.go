package db

import (
	"time"
)

// Database interface
type Database interface {
	init(cfg string) error
	destroy() error
	// return feeds keyed by ID
	getFeeds() (map[string]*FeedInfo, error)
}

// FeedInfo feed descriptior
type FeedInfo struct {
	// Etag from HTTP Request - used for conditional GETs
	ETag string
	// Type of feed: Atom, RSS2, etc.
	Flavor string
	// Not yet in use: how often to poll the feed
	Frequency int
	// Feed generator
	Generator string
	// Hub URL
	Hub string
	// Feed icon
	Icon string
	// UUID
	ID string
	// Time of last fetch attempt (LastStatus is true) or completion (LastStatus is false)
	LastFetch time.Time
	// Last-Modified time from HTTP Request - used for conditional GETs
	LastModified time.Time
	// The last status, true if error, false if successfully completed
	LastStatus bool
	// Time the feed was last updated (from feed)
	LastUpdated time.Time
	// Feed title
	Title string
	// URL updated if feed is permenently redirected
	URL string
}
