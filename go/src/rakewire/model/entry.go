package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Entry from a feed
type Entry struct {
	ID      string    `json:"id" kv:"key"`
	GUID    string    `json:"-" kv:"GUID:2"`
	FeedID  string    `json:"feedId" kv:"Date:1,GUID:1"`
	Created time.Time `json:"created" kv:"Date:2"`
	Updated time.Time `json:"updated"`
	URL     string    `json:"url"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Content string    `json:"contents"`
}

// NewEntry instantiate a new Entry object with a new UUID
func NewEntry(feedID string, guID string) *Entry {
	return &Entry{
		GUID:   guID,
		FeedID: feedID,
	}
}

// GenerateNewID generates a new UUID for Entry.ID
func (z *Entry) GenerateNewID() {
	z.ID = getUUID()
}

// Hash generated a fingerprint for the entry to test if it has been updated or not.
func (z *Entry) Hash() string {
	hash := sha256.New()
	hash.Write([]byte(z.Author))
	hash.Write([]byte(z.Content))
	hash.Write([]byte(z.Title))
	hash.Write([]byte(z.URL))
	return hex.EncodeToString(hash.Sum(nil))
}
