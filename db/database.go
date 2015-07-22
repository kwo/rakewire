package db

import (
	m "rakewire.com/model"
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
	Repair() error
}
