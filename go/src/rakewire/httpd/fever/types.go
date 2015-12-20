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
	Version       int      `json:"api_version"`
	Authorized    int      `json:"auth"`
	LastRefreshed int64    `json:"last_refreshed_on_time,string,omitempty"`
	Groups        []*Group `json:"groups,omitempty"`
}

// Group is the fever group construct
type Group struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}
