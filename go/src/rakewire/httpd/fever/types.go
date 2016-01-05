package fever

import (
	m "rakewire/model"
	"time"
)

// API top level struct
type API struct {
	prefix string
	db     Database
}

// Database defines the interface to the database
type Database interface {
	GroupGetAllByUser(userID uint64) ([]*m.Group, error)
	UserGetByFeverHash(feverhash string) (*m.User, error)
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
	UserFeedGetAllByUser(userID uint64) ([]*m.UserFeed, error)
	FeedLogGetLastFetchTime() (time.Time, error)
}

// Response defines the json/xml response return by requests.
type Response struct {
	Version       int          `json:"api_version"`
	Authorized    uint8        `json:"auth"`
	LastRefreshed int64        `json:"last_refreshed_on_time,string,omitempty"`
	Groups        []*Group     `json:"groups,omitempty"`
	FeedGroups    []*FeedGroup `json:"feeds_groups,omitempty"`
	Feeds         []*Feed      `json:"feeds,omitempty"`
	Items         []*Item      `json:"items,omitempty"`
	ItemCount     uint         `json:"total_items,omitempty"`
	UnreadItemIDs string       `json:"unread_item_ids,omitempty"`
	SavedItemIDs  string       `json:"saved_item_ids,omitempty"`
}

// Group is the fever group construct
type Group struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

// FeedGroup is the fever feed-group construct
type FeedGroup struct {
	GroupID uint64 `json:"group_id"`
	FeedIDs string `json:"feed_ids"` // comma separated
}

// Feed is a fever feed construct
type Feed struct {
	ID          uint64 `json:"id"`
	FaviconID   uint64 `json:"favicon_id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	SiteURL     string `json:"site_url"`
	IsSpark     uint8  `json:"is_spark"`
	LastUpdated int64  `json:"last_updated_on_time,string"`
}

// Item is a fever item construct
type Item struct {
	ID         uint64 `json:"id"`
	UserFeedID uint64 `json:"feed_id"`
	Title      string `json:"title"`
	Author     string `json:"author"`
	HTML       string `json:"html"`
	URL        string `json:"url"`
	IsSaved    uint8  `json:"is_saved"`
	IsRead     uint8  `json:"is_read"`
	Created    int64  `json:"created_on_time"`
}
