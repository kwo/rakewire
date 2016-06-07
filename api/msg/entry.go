package msg

import (
	"time"
)

// Entries is a list of Entry structs
type Entries []*Entry

// Entry defines an entry in a subscription
type Entry struct {
	Subscription string    `json:"subscription,omitempty"`
	GUID         string    `json:"guid,omitempty"`
	Title        string    `json:"title,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
	Read         bool      `json:"read,omitempty"`
	Star         bool      `json:"star,omitempty"`
}

// EntryListRequest defines the request to list entries
type EntryListRequest struct {
	Subscription string `json:"subscription,omitempty"`
	// TODO: group
	// TODO: min/max times
	// TODO: unread only
	// TODO: starred only
}

// EntryListResponse returns a list of entries
type EntryListResponse struct {
	Status  int     `json:"status"`
	Message string  `json:"message,omitempty"`
	Entries Entries `json:"entries,omitempty"`
}

// EntryUpdateRequest defines the request to update entries
type EntryUpdateRequest struct {
	Entries Entries `json:"entries,omitempty"`
}

// EntryUpdateResponse returns the status of an update request
type EntryUpdateResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

// BySubscription groups elements in the Entries collection by Subscription
func (z Entries) BySubscription() map[string]Entries {
	result := make(map[string]Entries)
	for _, entry := range z {
		entries := result[entry.Subscription]
		entries = append(entries, entry)
		result[entry.Subscription] = entries
	}
	return result
}
