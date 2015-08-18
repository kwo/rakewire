package db

import (
	m "rakewire/model"
	"time"
)

// Configuration configuration
type Configuration struct {
	Location string
}

// Database interface
type Database interface {
	GetFeedByID(id string) (*m.Feed, error)
	GetFeedByURL(url string) (*m.Feed, error)
	GetFeeds() (*m.Feeds, error)
	GetFetchFeeds(max *time.Time) (*m.Feeds, error)
	SaveFeeds(*m.Feeds) error
	// GetFeedLog retrieves the past fetch attempts for the feed in reverse chronological order.
	// If since is equal to 0, return all.
	GetFeedLog(id string, since time.Duration) ([]*m.FeedLog, error)
	Repair() error
}
