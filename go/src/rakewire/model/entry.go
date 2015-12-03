package model

import (
	"github.com/pborman/uuid"
	"time"
)

// Entry from a feed
type Entry struct {
	ID      string    `json:"id" kv:"key"`
	EntryID string    `json:"-" kv:"EntryID:2"`
	FeedID  string    `json:"feedId" kv:"Date:1,EntryID:1"`
	Created time.Time `json:"created" kv:"Date:2"`
	Updated time.Time `json:"updated"`
	URL     string    `json:"url"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Content string    `json:"contents"`
}

// NewEntry instantiate a new Entry object with a new UUID
func NewEntry(feedID string, entryID string) *Entry {
	return &Entry{
		EntryID: entryID,
		FeedID:  feedID,
	}
}

// GenerateNewID generates a new UUID for Entry.ID
func (z *Entry) GenerateNewID() {
	z.ID = uuid.NewUUID().String()
}
