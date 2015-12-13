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

// GetFeedEntriesFromIDs retrieves entries for specific GUIDs
func (z *Service) GetFeedEntriesFromIDs(feedID string, guIDs []string) (map[string]*m.Entry, error) {

	result := make(map[string]*m.Entry)

	for _, guID := range guIDs {

		entries := []*m.Entry{}
		add := func() interface{} {
			e := m.NewEntry(feedID, guID)
			entries = append(entries, e)
			return e
		}

		err := z.db.View(func(tx *bolt.Tx) error {
			return Query("Entry", "GUID", []interface{}{feedID, guID}, []interface{}{feedID, guID}, add, tx)
		})
		if err != nil {
			return nil, err
		}

		if len(entries) == 1 {
			result[entries[0].GUID] = entries[0]
		} else if len(entries) > 1 {
			return nil, fmt.Errorf("Unique index returned multiple results: %s, FeedID: %s, GUID: %s", "Entry/GUID", feedID, guID)
		}

	} // loop guIDs

	return result, nil

}
