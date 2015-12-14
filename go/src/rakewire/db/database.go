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
	UserGetByUsername(username string) (*m.User, error)
	UserGetByFeverHash(feverhash string) (*m.User, error)
	UserSave(user *m.User) error

	GetFeedByID(feedID string) (*m.Feed, error)
	GetFeedByURL(url string) (*m.Feed, error)
	GetFeeds() ([]*m.Feed, error)
	// GetFetchFeeds get feeds to be fetched within the given max time parameter.
	GetFetchFeeds(max time.Time) ([]*m.Feed, error)
	// GetFeedLog retrieves the past fetch attempts for the feed in reverse chronological order.
	// If since is equal to 0, return all.
	GetFeedLog(feedID string, since time.Duration) ([]*m.FeedLog, error)
	GetFeedEntries(feedID string, since time.Duration) ([]*m.Entry, error)
	GetFeedEntriesFromIDs(feedID string, guIDs []string) (map[string]*m.Entry, error)
	SaveFeed(*m.Feed) error

	Repair() error
}

// DataObject defines the functions necessary for objects to be persisted to the database
type DataObject interface {
	GetName() string
	GetID() uint64
	SetID(id uint64)
	Clear()
	Serialize() map[string]string
	Deserialize(map[string]string) error
	IndexKeys() map[string][]string
}
