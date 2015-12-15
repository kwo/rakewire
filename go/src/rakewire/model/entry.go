package model

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Entry from a feed
type Entry struct {
	ID      uint64    `json:"id"`
	GUID    string    `json:"-"`
	FeedID  uint64    `json:"feedId"`
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
	z.FeedID = 0
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
	setUint(z.FeedID, eFeedID, result)
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
	var errors []error
	z.ID = getUint(eID, values, errors)
	z.GUID = getString(eGUID, values, errors)
	z.FeedID = getUint(eFeedID, values, errors)
	z.Created = getTime(eCreated, values, errors)
	z.Updated = getTime(eUpdated, values, errors)
	z.URL = getString(eURL, values, errors)
	z.Author = getString(eAuthor, values, errors)
	z.Title = getString(eTitle, values, errors)
	z.Content = getString(eContent, values, errors)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Entry) IndexKeys() map[string][]string {
	data := z.Serialize()
	result := make(map[string][]string)
	result[EntryIndexDate] = []string{data[eFeedID], data[eCreated]}
	result[EntryIndexGUID] = []string{data[eFeedID], data[eGUID]}
	// FIXME: index IDs must be zero-padded
	return result
}
