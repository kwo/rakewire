package modelng

import (
	"time"
)

const (
	entityEntry               = "Entry"
	indexEntryFeedReadUpdated = "FeedReadUpdated"
	indexEntryFeedStarUpdated = "FeedStarUpdated"
	indexEntryFeedUpdated     = "FeedUpdated"
	indexEntryReadUpdated     = "ReadUpdated"
	indexEntryStarUpdated     = "StarUpdated"
	indexEntryUpdated         = "Updated"
)

var (
	indexesEntry = []string{
		indexEntryFeedReadUpdated, indexEntryFeedStarUpdated, indexEntryFeedUpdated, indexEntryReadUpdated, indexEntryStarUpdated, indexEntryUpdated,
	}
)

// Entries is a collection of Entry elements
type Entries []*Entry

// Unique returns an array of unique Entry elements
func (z Entries) Unique() Entries {

	uniques := make(map[string]*Entry)
	for _, entry := range z {
		uniques[entry.getID()] = entry
	}

	entries := Entries{}
	for _, entry := range uniques {
		entries = append(entries, entry)
	}

	return entries

}

// Reverse reverses the order of the collection
func (z Entries) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

// Entry defines an item status for a user
type Entry struct {
	UserID  string    `json:"userId"`
	ItemID  string    `json:"itemId"`
	FeedID  string    `json:"feedId"`
	Updated time.Time `json:"updated,omitempty"`
	Read    bool      `json:"read,omitempty"`
	Star    bool      `json:"star,omitempty"`
}
