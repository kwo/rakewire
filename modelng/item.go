package modelng

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const (
	entityItem    = "Item"
	indexItemGUID = "GUID"
)

var (
	indexesItem = []string{
		indexItemGUID,
	}
)

// Items is a collection Item objects
type Items []*Item

// Item from a feed
type Item struct {
	ID      string    `json:"id"`
	GUID    string    `json:"guid" `
	FeedID  string    `json:"feedId"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
	URL     string    `json:"url,omitempty"`
	Author  string    `json:"author,omitempty"`
	Title   string    `json:"title,omitempty"`
	Content string    `json:"content,omitempty"`
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
