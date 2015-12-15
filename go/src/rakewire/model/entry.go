package model

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

// Entry from a feed
type Entry struct {
	ID      uint64    `json:"id"`
	GUID    string    `json:"-"`
	FeedID  string    `json:"feedId"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
	URL     string    `json:"url"`
	Author  string    `json:"author"`
	Title   string    `json:"title"`
	Content string    `json:"contents"`
}

// index constants
const (
	EntryEntity    = "Entry"
	EntryIndexGUID = "GUID"
	EntryIndexDate = "Date"
)

const (
	eID      = "ID"
	eGUID    = "GUID"
	eFeedID  = "FeedID"
	eCreated = "Created"
	eUpdated = "Updated"
	eURL     = "URL"
	eAuthor  = "Author"
	eTitle   = "Title"
	eContent = "Content"
)

// NewEntry instantiate a new Entry object with a new UUID
func NewEntry(feedID string, guID string) *Entry {
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

// GetName return the name of the entity.
func (z *Entry) GetName() string {
	return EntryEntity
}

// GetID return the primary key of the object.
func (z *Entry) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *Entry) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Entry) Clear() {
	z.ID = 0
	z.GUID = empty
	z.FeedID = empty
	z.Created = time.Time{}
	z.Updated = time.Time{}
	z.URL = empty
	z.Author = empty
	z.Title = empty
	z.Content = empty
}

// Serialize serializes an object to a list of key-values.
func (z *Entry) Serialize() map[string]string {
	result := make(map[string]string)
	setUint(z.ID, eID, result)
	setString(z.GUID, eGUID, result)
	setString(z.FeedID, eFeedID, result)
	setTime(z.Created, eCreated, result)
	setTime(z.Updated, eUpdated, result)
	setString(z.URL, eURL, result)
	setString(z.Author, eAuthor, result)
	setString(z.Title, eTitle, result)
	setString(z.Content, eContent, result)
	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *Entry) Deserialize(values map[string]string) error {

	for k, v := range values {
		switch k {
		case eID:
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			z.ID = id
		case eGUID:
			z.GUID = v
		case eFeedID:
			z.FeedID = v
		case eCreated:
			t, err := time.Parse(timeFormat, v)
			if err != nil {
				return err
			}
			z.Created = t
		case eUpdated:
			t, err := time.Parse(timeFormat, v)
			if err != nil {
				return err
			}
			z.Updated = t
		case eURL:
			z.URL = v
		case eAuthor:
			z.Author = v
		case eTitle:
			z.Title = v
		case eContent:
			z.Content = v
		}
	}

	return nil

}

// IndexKeys returns the keys of all indexes for this object.
func (z *Entry) IndexKeys() map[string][]string {
	data := z.Serialize()
	result := make(map[string][]string)
	result[EntryIndexDate] = []string{data[eFeedID], data[eCreated]}
	result[EntryIndexGUID] = []string{data[eFeedID], data[eGUID]}
	return result
}
