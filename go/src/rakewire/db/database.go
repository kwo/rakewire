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
	GroupGetAllByUser(userID uint64) ([]*m.Group, error)

	UserGetByUsername(username string) (*m.User, error)
	UserGetByFeverHash(feverhash string) (*m.User, error)
	UserSave(user *m.User) error

	UserEntryGetTotalForUser(userID uint64) (uint, error)
	UserEntryGet(userID uint64, ids []uint64) ([]*m.UserEntry, error)
	UserEntryGetNext(userID uint64, minID uint64, count int) ([]*m.UserEntry, error)
	UserEntryGetPrev(userID uint64, maxID uint64, count int) ([]*m.UserEntry, error)
	//UserEntryGetUnreadForUser(userID uint64) ([]*m.UserEntry, error)
	//UserEntrySave(userentries []*m.UserEntry) error

	UserFeedGetAllByUser(userID uint64) ([]*m.UserFeed, error)

	GetFeedByID(feedID uint64) (*m.Feed, error)
	GetFeedByURL(url string) (*m.Feed, error)
	GetFeeds() ([]*m.Feed, error)
	// GetFetchFeeds get feeds to be fetched within the given max time parameter.
	GetFetchFeeds(max time.Time) ([]*m.Feed, error)
	// GetFeedLog retrieves the past fetch attempts for the feed in reverse chronological order.
	// If since is equal to 0, return all.
	GetFeedLog(feedID uint64, since time.Duration) ([]*m.FeedLog, error)
	GetFeedEntriesFromIDs(feedID uint64, guIDs []string) (map[string]*m.Entry, error)
	FeedSave(*m.Feed) ([]*m.Entry, error)
}

// DataObject defines the functions necessary for objects to be persisted to the database
type DataObject interface {
	GetID() uint64
	SetID(id uint64)
	Clear()
	Serialize() map[string]string
	Deserialize(map[string]string) error
	IndexKeys() map[string][]string
}
