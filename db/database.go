package db

import (
	"time"
)

// Database interface
type Database interface {
	GetFeedByID(id string) (*Feed, error)
	GetFeedByURL(url string) (*Feed, error)
	GetFeeds() (*Feeds, error)
	GetFetchFeeds(max *time.Time) (*Feeds, error)
	SaveFeeds(*Feeds) error
	Repair() error
}

// Configuration configuration
type Configuration struct {
	Location string
}
