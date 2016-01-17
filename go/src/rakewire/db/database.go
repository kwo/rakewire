package db

import (
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"time"
)

// Configuration configuration
type Configuration struct {
	Location string
}

// Database interface
type Database interface {
	BoltDB() *bolt.DB

	// Feed

	FeedDelete(feed *m.Feed) error
	FeedSave(*m.Feed) ([]*m.Entry, error)
	GetFeedByID(feedID uint64) (*m.Feed, error)
	GetFeedByURL(url string) (*m.Feed, error)
	GetFeeds() ([]*m.Feed, error)
	// GetFetchFeeds get feeds to be fetched within the given max time parameter.
	GetFetchFeeds(max time.Time) ([]*m.Feed, error)
	FeedDuplicates() (map[string][]uint64, error)

	// FeedLog

	// GetFeedLog retrieves the past fetch attempts for the feed in reverse chronological order.
	// If since is equal to 0, return all.
	GetFeedLog(feedID uint64, since time.Duration) ([]*m.FeedLog, error)
	FeedLogGetLastFetchTime() (time.Time, error)

	// Entry

	EntriesGet(feedID uint64) ([]*m.Entry, error)
	GetFeedEntriesFromIDs(feedID uint64, guIDs []string) (map[string]*m.Entry, error)

	// User

	UserSave(user *m.User) error
	UserGetByUsername(username string) (*m.User, error)
	UserGetByFeverHash(feverhash string) (*m.User, error)

	// Group

	GroupDelete(group *m.Group) error
	GroupSave(group *m.Group) error
	GroupGetAllByUser(userID uint64) ([]*m.Group, error)

	// UserFeed

	UserFeedDelete(userfeed *m.UserFeed) error
	UserFeedSave(userfeed *m.UserFeed) error
	UserFeedGetAllByUser(userID uint64) ([]*m.UserFeed, error)
	UserFeedGetByFeed(feedID uint64) ([]*m.UserFeed, error)

	// UserEntry

	UserEntryAdd(allEntries []*m.Entry) error
	UserEntryGetTotalForUser(userID uint64) (uint, error)

	UserEntryGetByID(userID uint64, ids []uint64) ([]*m.UserEntry, error)
	UserEntryGetNext(userID uint64, minID uint64, count int) ([]*m.UserEntry, error)
	UserEntryGetPrev(userID uint64, maxID uint64, count int) ([]*m.UserEntry, error)
	UserEntryGetUnreadForUser(userID uint64) ([]*m.UserEntry, error)
	UserEntryGetStarredForUser(userID uint64) ([]*m.UserEntry, error)

	UserEntrySave(userentries []*m.UserEntry) error
	UserEntryUpdateReadByFeed(userID, userFeedID uint64, maxTime time.Time, read bool) error
	UserEntryUpdateStarByFeed(userID, userFeedID uint64, maxTime time.Time, star bool) error
	UserEntryUpdateReadByGroup(userID, groupID uint64, maxTime time.Time, read bool) error
	UserEntryUpdateStarByGroup(userID, groupID uint64, maxTime time.Time, star bool) error
}

// DataObject defines the functions necessary for objects to be persisted to the database
type DataObject interface {
	GetID() uint64
	SetID(id uint64)
	Clear()
	Serialize(flags ...bool) map[string]string
	Deserialize(map[string]string) error
	IndexKeys() map[string][]string
}
