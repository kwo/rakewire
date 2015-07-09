package db

import (
	m "rakewire.com/model"
)

// Database interface
type Database interface {
	Open(cfg *m.DatabaseConfiguration) error
	Close() error
	// return feeds keyed by ID
	GetFeeds() (map[string]*FeedInfo, error)
	SaveFeeds([]*FeedInfo) (int, error)
}
