package bolt

import (
	"github.com/boltdb/bolt"
	m "rakewire/model"
)

// GetFeedEntriesFromIDs retrieves entries for specific GUIDs
func (z *Service) GetFeedEntriesFromIDs(feedID string, guIDs []string) (map[string]*m.Entry, error) {

	result := make(map[string]*m.Entry)

	for _, guID := range guIDs {

		err := z.db.View(func(tx *bolt.Tx) error {
			data, err := kvGetIndex(m.EntityEntry, m.IndexEntryGUID, []string{feedID, guID}, tx)
			if err != nil {
				return err
			}
			if data != nil {
				e := &m.Entry{}
				if err := e.Deserialize(data); err != nil {
					return err
				}
				result[e.GUID] = e
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

	}

	return result, nil

}
