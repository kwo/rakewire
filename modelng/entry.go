package modelng

import (
	"time"
)

const (
	entityEntry    = "Entry"
	indexEntryRead = "Read"
	indexEntryStar = "Star"
)

var (
	indexesEntry = []string{
		indexEntryRead, indexEntryStar,
	}
)

// Entries is a collection of Entry elements
type Entries []*Entry

// Entry defines an item status for a user
type Entry struct {
	UserID  string    `json:"userID"`
	ItemID  string    `json:"itemID"`
	Updated time.Time `json:"updated,omitempty"`
	Read    bool      `json:"read,omitempty"`
	Star    bool      `json:"star,omitempty"`
}
