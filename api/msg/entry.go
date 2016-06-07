package msg

import (
	"time"
)

// Entries is a list of Entry structs
type Entries []*Entry

// Entry defines an entry in a subscription
type Entry struct {
	Subscription string    `json:"subscription,omitempty"`
	GUID         string    `json:"guid,omitempty"`
	Title        string    `json:"title,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
	Read         bool      `json:"read,omitempty"`
	Star         bool      `json:"star,omitempty"`
}

// EntryListRequest defines the request to list entries
type EntryListRequest struct {
	Subscription string `json:"subscription,omitempty"`
}

// EntryListResponse returns a list of entries
type EntryListResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message,omitempty"`
	Entries Entries `json:"entries,omitempty"`
}
