package bolt

import (
	m "rakewire/model"
	"time"
)

// GetFeedEntries retrieves entries for a feed since a specific time
func (z *Database) GetFeedEntries(feedID string, since time.Duration) ([]*m.Entry, error) {
	// TODO: get feed entries
	return nil, nil
}
