package model

//go:generate gokv $GOFILE

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Item from a feed
type Item struct {
	ID      uint64
	GUID    string `kv:"+groupby,GUID:2"`
	FeedID  uint64 `kv:"+required,+groupall,GUID:1"`
	Created time.Time
	Updated time.Time
	URL     string
	Author  string
	Title   string
	Content string
}

// NewItem instantiate a new Item object
func NewItem(feedID uint64, guID string) *Item {
	return &Item{
		FeedID: feedID,
		GUID:   guID,
	}
}

func (z *Item) setIDIfNecessary(fn fnNextID, tx Transaction) error {
	if z.ID == 0 {
		if id, _, err := fn(itemEntity, tx); err == nil {
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
