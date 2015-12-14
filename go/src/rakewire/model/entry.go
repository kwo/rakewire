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
	EntityEntry    = "Entry"
	IndexEntryGUID = "GUID"
	IndexEntryDate = "Date"
)

const (
	fGUID    = "GUID"
	fFeedID  = "FeedID"
	fCreated = "Created"
	fUpdated = "Updated"
	fURL     = "URL"
	fAuthor  = "Author"
	fTitle   = "Title"
	fContent = "Content"
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
	return EntityEntry
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

	if z.ID != 0 {
		result[fID] = strconv.FormatUint(z.ID, 10)
	}

	if z.GUID != empty {
		result[fGUID] = z.GUID
	}

	if z.FeedID != empty {
		result[fFeedID] = z.FeedID
	}

	if !z.Created.IsZero() {
		result[fCreated] = z.Created.Format(timeFormat)
	}

	if !z.Updated.IsZero() {
		result[fUpdated] = z.Updated.Format(timeFormat)
	}

	if z.URL != empty {
		result[fURL] = z.URL
	}

	if z.Author != empty {
		result[fAuthor] = z.Author
	}

	if z.Title != empty {
		result[fTitle] = z.Title
	}

	if z.Content != empty {
		result[fContent] = z.Content
	}

	return result

}

// Deserialize serializes an object to a list of key-values.
func (z *Entry) Deserialize(values map[string]string) error {

	for k, v := range values {
		switch k {
		case fID:
			id, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			z.ID = id
		case fGUID:
			z.GUID = v
		case fFeedID:
			z.FeedID = v
		case fCreated:
			t, err := time.Parse(timeFormat, v)
			if err != nil {
				return err
			}
			z.Created = t
		case fUpdated:
			t, err := time.Parse(timeFormat, v)
			if err != nil {
				return err
			}
			z.Updated = t
		case fURL:
			z.URL = v
		case fAuthor:
			z.Author = v
		case fTitle:
			z.Title = v
		case fContent:
			z.Content = v
		}
	}

	return nil

}

// IndexKeys returns the keys of all indexes for this object.
func (z *Entry) IndexKeys() map[string][]string {
	data := z.Serialize()
	result := make(map[string][]string)
	result[IndexEntryDate] = []string{data[fFeedID], data[fCreated]}
	result[IndexEntryGUID] = []string{data[fFeedID], data[fGUID]}
	return result
}
