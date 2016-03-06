package modelng

import (
	"time"
)

const (
	entityEntry       = "Entry"
	indexEntryRead    = "Read"
	indexEntryStar    = "Star"
	indexEntryUpdated = "Updated"
)

var (
	indexesEntry = []string{
		indexEntryRead, indexEntryStar, indexEntryUpdated,
	}
)

// Entries is a collection of Entry elements
type Entries []*Entry

// Entry defines an item status for a user
type Entry struct {
	UserID  string    `json:"userId"`
	ItemID  string    `json:"itemId"`
	Updated time.Time `json:"updated,omitempty"`
	Read    bool      `json:"read,omitempty"`
	Star    bool      `json:"star,omitempty"`
}
