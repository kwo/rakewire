package model

//go:generate gokv $GOFILE

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Item from a feed
type Item struct {
	ID      uint64    `json:"id"`
	GUID    string    `json:"-" kv:"GUID:2"`
	FeedID  uint64    `json:"feedId" kv:"+required,GUID:1"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	URL     string    `json:"url"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Content string    `json:"contents"`
}

// NewItem instantiate a new Item object
func NewItem(feedID uint64, guID string) *Item {
	return &Item{
		FeedID: feedID,
		GUID:   guID,
	}
}

// Hash generated a fingerprint for the item to test if it has been updated or not.
func (z *Item) Hash() string {
	hash := sha256.New()
	hash.Write([]byte(z.Author))
	hash.Write([]byte(z.Content))
	hash.Write([]byte(z.Title))
	hash.Write([]byte(z.URL))
	return hex.EncodeToString(hash.Sum(nil))
}