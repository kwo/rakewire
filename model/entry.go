package model

import (
	"encoding/json"
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

// Entry defines an item status for a user
type Entry struct {
	UserID  string    `json:"userId"`
	ItemID  string    `json:"itemId"`
	FeedID  string    `json:"feedId"`
	Updated time.Time `json:"updated,omitempty"`
	Read    bool      `json:"read,omitempty"`
	Star    bool      `json:"star,omitempty"`
}

// GetID returns the unique ID for the object
func (z *Entry) GetID() string {
	return keyEncode(z.UserID, z.ItemID)
}

func (z *Entry) clear() {
	z.UserID = empty
	z.ItemID = empty
	z.FeedID = empty
	z.Updated = time.Time{}
	z.Read = false
	z.Star = false
}

func (z *Entry) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Entry) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Entry) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexEntryFeedReadUpdated] = []string{z.UserID, z.FeedID, keyEncodeBool(z.Read), keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryFeedStarUpdated] = []string{z.UserID, z.FeedID, keyEncodeBool(z.Star), keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryFeedUpdated] = []string{z.UserID, z.FeedID, keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryReadUpdated] = []string{z.UserID, keyEncodeBool(z.Read), keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryStarUpdated] = []string{z.UserID, keyEncodeBool(z.Star), keyEncodeTime(z.Updated), z.ItemID}
	result[indexEntryUpdated] = []string{z.UserID, keyEncodeTime(z.Updated), z.ItemID}
	return result
}

func (z *Entry) setID(tx Transaction) error {
	return nil
}

// Entries is a collection of Entry elements
type Entries []*Entry

// Limit truncate the size of the collection to the given number of items.
func (z Entries) Limit(limit uint) Entries {
	if len(z) > int(limit) {
		return z[:limit]
	}
	return z
}

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

func (z *Entries) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Entries) encode() ([]byte, error) {
	return json.Marshal(z)
}
