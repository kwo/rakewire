package model

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

// ByID maps items to their ID
func (z Items) ByID() map[string]*Item {
	result := make(map[string]*Item)
	for _, item := range z {
		result[item.ID] = item
	}
	return result
}

// ByURL maps items to their URL
func (z Items) ByURL() map[string]*Item {
	result := make(map[string]*Item)
	for _, item := range z {
		result[item.URL] = item
	}
	return result
}

// GroupByFeedID groups collections of elements in Items by FeedID
func (z Items) GroupByFeedID() map[string]Items {
	result := make(map[string]Items)
	for _, item := range z {
		a := result[item.FeedID]
		a = append(a, item)
		result[item.FeedID] = a
	}
	return result
}

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
