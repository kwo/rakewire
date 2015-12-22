package fever

import (
	m "rakewire/model"
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
	UserFeedGetAllByUser(userID uint64) ([]*m.UserFeed, error)
}

// Response defines the json/xml response return by requests.
type Response struct {
	Version       int          `json:"api_version"`
	Authorized    int          `json:"auth"`
	LastRefreshed int64        `json:"last_refreshed_on_time,string,omitempty"`
	Feeds         []*Feed      `json:"feeds,omitempty"`
	FeedGroups    []*FeedGroup `json:"feed_groups,omitempty"`
	Groups        []*Group     `json:"groups,omitempty"`
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
	IsSpark     bool   `json:"is_spark"`
	LastUpdated int64  `json:"last_updated_on_time,string"`
}
