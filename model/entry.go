package model

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
		indexEntryFeedReadUpdated, indexEntryFeedStarUpdated, indexEntryFeedUpdated,
		indexEntryReadUpdated, indexEntryStarUpdated, indexEntryUpdated,
	}
)

// Entries is a collection of Entry elements
type Entries []*Entry

// Reverse reverses the order of the collection
func (z Entries) Reverse() Entries {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
	return z
}

// Unique returns an array of unique Entry elements
func (z Entries) Unique() Entries {

	uniques := make(map[string]*Entry)
	for _, entry := range z {
		uniques[entry.GetID()] = entry
	}

	entries := Entries{}
	for _, entry := range uniques {
		entries = append(entries, entry)
	}

	return entries

}

// Limit truncate the size of the collection to the given number of items.
func (z Entries) Limit(limit uint) Entries {
	if len(z) > int(limit) {
		return z[:limit]
	}
	return z
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
