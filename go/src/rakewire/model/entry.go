package model

//go:generate gokv $GOFILE

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Entry from a feed
type Entry struct {
	ID      uint64    `json:"id"`
	GUID    string    `json:"-" kv:"GUID:2"`
	FeedID  uint64    `json:"feedId" kv:"GUID:1"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	URL     string    `json:"url"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Content string    `json:"contents"`
}

// NewEntry instantiate a new Entry object
func NewEntry(feedID uint64, guID string) *Entry {
	return &Entry{
		FeedID: feedID,
		GUID:   guID,
	}
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
