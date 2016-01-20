package db

import (
	"rakewire/model"
	"time"
)

// Configuration configuration
type Configuration struct {
	Location string
}

// Database interface
type Database interface {
	Location() string
	Select(fn func(tx model.Transaction) error) error
	Update(fn func(tx model.Transaction) error) error
	Repair() error
	ModelDatabase() model.Database

	// Feed

	FeedDelete(feed *model.Feed) error
	FeedSave(*model.Feed) ([]*model.Entry, error)
	GetFeedByID(feedID uint64) (*model.Feed, error)
	GetFeedByURL(url string) (*model.Feed, error)
	GetFeeds() ([]*model.Feed, error)
	// GetFetchFeeds get feeds to be fetched within the given max time parameter.
	GetFetchFeeds(max time.Time) ([]*model.Feed, error)
	FeedDuplicates() (map[string][]uint64, error)

	// FeedLog

	// GetFeedLog retrieves the past fetch attempts for the feed in reverse chronological order.
	// If since is equal to 0, return all.
	GetFeedLog(feedID uint64, since time.Duration) ([]*model.FeedLog, error)
	FeedLogGetLastFetchTime() (time.Time, error)

	// Entry

	EntriesGet(feedID uint64) ([]*model.Entry, error)
	GetFeedEntriesFromIDs(feedID uint64, guIDs []string) (map[string]*model.Entry, error)

	// Group

	GroupDelete(group *model.Group) error
	GroupSave(group *model.Group) error
	GroupGetAllByUser(userID uint64) ([]*model.Group, error)

	// UserFeed

	UserFeedDelete(userfeed *model.UserFeed) error
	UserFeedSave(userfeed *model.UserFeed) error
	UserFeedGetAllByUser(userID uint64) ([]*model.UserFeed, error)
	UserFeedGetByFeed(feedID uint64) ([]*model.UserFeed, error)

	// UserEntry

	UserEntryAdd(allEntries []*model.Entry) error
	UserEntryGetTotalForUser(userID uint64) (uint, error)

	UserEntryGetByID(userID uint64, ids []uint64) ([]*model.UserEntry, error)
	UserEntryGetNext(userID uint64, minID uint64, count int) ([]*model.UserEntry, error)
	UserEntryGetPrev(userID uint64, maxID uint64, count int) ([]*model.UserEntry, error)
	UserEntryGetUnreadForUser(userID uint64) ([]*model.UserEntry, error)
	UserEntryGetStarredForUser(userID uint64) ([]*model.UserEntry, error)

	UserEntrySave(userentries []*model.UserEntry) error
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
