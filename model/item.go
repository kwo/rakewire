package model

//go:generate gokv $GOFILE

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Item from a feed
type Item struct {
	ID      string    `json:"id"`
	GUID    string    `json:"guid" kv:"+groupby,GUID:2"`
	FeedID  string    `json:"feedID" kv:"+required,+groupall,GUID:1"`
	Created time.Time `json:"created,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
	URL     string    `json:"url,omitempty"`
	Author  string    `json:"author,omitempty"`
	Title   string    `json:"title,omitempty"`
	Content string    `json:"content,omitempty"`
}

// NewItem instantiate a new Item object
func NewItem(feedID string, guID string) *Item {
	return &Item{
		FeedID: feedID,
		GUID:   guID,
	}
}

func (z *Item) setID(fn fnUniqueID) error {
	if z.ID == empty {
		if id, err := fn(); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
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
