package fever

import (
	"encoding/xml"
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
	Version       int      `json:"api_version" xml:"api_version"`
	Authorized    int      `json:"auth" xml:"auth"`
	LastRefreshed int64    `json:"last_refreshed_on_time,string,omitempty" xml:"last_refreshed_on_time,omitempty"`
	Groups        []*Group `json:"groups,omitempty" xml:"groups,omitempty"`
	XMLName       xml.Name `json:"-" xml:"response"`
}

// Group is the fever group construct
type Group struct {
	ID      uint64   `json:"id" xml:"id"`
	Title   string   `json:"title" xml:"title"`
	XMLName xml.Name `json:"-" xml:"group"`
}
