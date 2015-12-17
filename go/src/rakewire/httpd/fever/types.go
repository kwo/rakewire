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
	UserGetByFeverHash(feverhash string) (*m.User, error)
}

// Response defines the json/xml response return by requests.
type Response struct {
	Version    int `json:"api_version" xml:"api_version"`
	Authorized int `json:"auth" xml:"auth"`
	// LastRefreshed is the time of the last refesh on the server expressed as seconds since the unix epoch.
	LastRefreshed int64    `json:"last_refreshed_on_time,string,omitempty" xml:"last_refreshed_on_time,omitempty"`
	XMLName       xml.Name `json:"-" xml:"response"`
}
