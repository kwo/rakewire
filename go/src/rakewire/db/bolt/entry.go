package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"time"
)

// GetFeedEntries retrieves entries for a feed since a specific time
func (z *Service) GetFeedEntries(feedID string, since time.Duration) ([]*m.Entry, error) {
	// TODO: get feed entries
	return nil, nil
}

// GetFeedEntriesFromIDs retrieves entries for specific entryIDs
func (z *Service) GetFeedEntriesFromIDs(feedID string, entryIDs []string) ([]*m.Entry, error) {

	result := []*m.Entry{}

	for _, entryID := range entryIDs {

		entries := []*m.Entry{}
		add := func() interface{} {
			e := m.NewEntry(feedID, entryID)
			entries = append(entries, e)
			return e
		}

		err := z.db.View(func(tx *bolt.Tx) error {
			return Query("Entry", "EntryID", []interface{}{feedID, entryID}, []interface{}{feedID, entryID}, add, tx)
		})
		if err != nil {
			return nil, err
		}

		if len(entries) == 0 {
			result = append(result, nil)
		} else if len(entries) > 1 {
			return nil, fmt.Errorf("Unique index returned multiple results: %s, FeedID: %s, EntryID: %s", "Entry/EntryID", feedID, entryID)
		} else {
			result = append(result, entries[0])
		}

	} // loop entryIDs

	return result, nil

}
