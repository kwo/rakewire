package model

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

// GetID returns the unique ID for the object
func (z *Item) GetID() string {
	return z.ID
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

func (z *Item) clear() {
	z.ID = empty
	z.GUID = empty
	z.FeedID = empty
	z.Created = time.Time{}
	z.Updated = time.Time{}
	z.URL = empty
	z.Author = empty
	z.Title = empty
	z.Content = empty
}

func (z *Item) decode(data []byte) error {
	z.clear()
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Item) encode() ([]byte, error) {
	return json.Marshal(z)
}

func (z *Item) indexes() map[string][]string {
	result := make(map[string][]string)
	result[indexItemGUID] = []string{z.FeedID, z.GUID}
	return result
}

func (z *Item) setID(tx Transaction) error {
	config := C.Get(tx)
	config.Sequences.Item = config.Sequences.Item + 1
	z.ID = keyEncodeUint(config.Sequences.Item)
	return C.Put(tx, config)
}

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

// GroupByGUID groups collections of elements in Items by GUID
func (z Items) GroupByGUID() map[string]*Item {
	result := make(map[string]*Item)
	for _, item := range z {
		result[item.GUID] = item
	}
	return result
}

func (z *Items) decode(data []byte) error {
	if err := json.Unmarshal(data, z); err != nil {
		return err
	}
	return nil
}

func (z *Items) encode() ([]byte, error) {
	return json.Marshal(z)
}
