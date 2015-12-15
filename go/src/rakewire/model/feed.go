package model

import (
	"time"
)

// Feed feed descriptor
type Feed struct {
	Attempt *FeedLog `json:"-"`
	Entries []*Entry `json:"-"`

	ID  uint64 `json:"id"`
	URL string `json:"url"`

	ETag         string    `json:"etag"`
	LastModified time.Time `json:"lastModified"`

	LastUpdated time.Time `json:"lastUpdated"`
	NextFetch   time.Time `json:"nextFetch"`

	Notes string `json:"notes,omitempty"`
	Title string `json:"title"`

	Status        string    `json:"status"`
	StatusMessage string    `json:"statusMessage"`
	StatusSince   time.Time `json:"statusSince"` // time of last status
}

// index constants
const (
	FeedEntity         = "Feed"
	FeedIndexNextFetch = "NextFetch"
	FeedIndexURL       = "URL"
)

const (
	fID            = "ID"
	fURL           = "URL"
	fETag          = "ETag"
	fLastModified  = "LastModified"
	fLastUpdated   = "LastUpdated"
	fNextFetch     = "NextFetch"
	fNotes         = "Notes"
	fTitle         = "Title"
	fStatus        = "Status"
	fStatusMessage = "StatusMessage"
	fStatusSince   = "StatusSince"
)

// NewFeed instantiate a new Feed object
func NewFeed(url string) *Feed {
	return &Feed{
		URL:       url,
		NextFetch: time.Now(),
	}
}

// AddEntry to the feed
func (z *Feed) AddEntry(guID string) *Entry {
	entry := NewEntry(z.ID, guID)
	z.Entries = append(z.Entries, entry)
	return entry
}

// UpdateFetchTime increases the fetch interval
func (z *Feed) UpdateFetchTime(lastUpdated time.Time) {

	if lastUpdated.IsZero() {
		return
	}

	z.LastUpdated = lastUpdated

	d := time.Now().Sub(z.LastUpdated) // how long since the last update?

	switch {
	case d < 30*time.Minute:
		z.AdjustFetchTime(10 * time.Minute)
	case d > 72*time.Hour:
		z.AdjustFetchTime(24 * time.Hour)
	case true:
		z.AdjustFetchTime(1 * time.Hour)
	}

}

// AdjustFetchTime sets the FetchTime to interval units in the future.
func (z *Feed) AdjustFetchTime(interval time.Duration) {
	z.NextFetch = time.Now().Add(interval)
}

// GetName return the name of the entity.
func (z *Feed) GetName() string {
	return FeedEntity
}

// GetID return the primary key of the object.
func (z *Feed) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *Feed) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Feed) Clear() {
	z.ID = 0
	z.URL = empty
	z.ETag = empty
	z.LastModified = time.Time{}
	z.LastUpdated = time.Time{}
	z.NextFetch = time.Time{}
	z.Notes = empty
	z.Title = empty
	z.Status = empty
	z.StatusMessage = empty
	z.StatusSince = time.Time{}
}

// Serialize serializes an object to a list of key-values.
func (z *Feed) Serialize() map[string]string {
	result := make(map[string]string)
	setUint(z.ID, fID, result)
	setString(z.URL, flURL, result)
	setString(z.ETag, fETag, result)
	setTime(z.LastModified, fLastModified, result)
	setTime(z.LastUpdated, fLastUpdated, result)
	setTime(z.NextFetch, fNextFetch, result)
	setString(z.Notes, fNotes, result)
	setString(z.Title, fTitle, result)
	setString(z.Status, fStatus, result)
	setString(z.StatusMessage, fStatusMessage, result)
	setTime(z.StatusSince, fStatusSince, result)
	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *Feed) Deserialize(values map[string]string) error {
	var errors []error
	z.ID = getUint(fID, values, errors)
	z.URL = getString(fURL, values, errors)
	z.ETag = getString(fETag, values, errors)
	z.LastModified = getTime(fLastModified, values, errors)
	z.LastUpdated = getTime(fLastUpdated, values, errors)
	z.NextFetch = getTime(fNextFetch, values, errors)
	z.Notes = getString(fNotes, values, errors)
	z.Title = getString(fTitle, values, errors)
	z.Status = getString(fStatus, values, errors)
	z.StatusMessage = getString(fStatusMessage, values, errors)
	z.StatusSince = getTime(fStatusSince, values, errors)
	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Feed) IndexKeys() map[string][]string {
	data := z.Serialize()
	result := make(map[string][]string)
	result[FeedIndexNextFetch] = []string{data[fNextFetch]}
	result[FeedIndexURL] = []string{data[fURL]}
	return result
}
