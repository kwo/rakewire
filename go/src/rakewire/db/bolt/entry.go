package bolt

import (
	"github.com/boltdb/bolt"
	m "rakewire/model"
)

// GetFeedEntriesFromIDs retrieves entries for specific GUIDs
func (z *Service) GetFeedEntriesFromIDs(feedID uint64, guIDs []string) (map[string]*m.Entry, error) {

	result := make(map[string]*m.Entry)

	e := &m.Entry{}
	e.FeedID = feedID

	for _, guID := range guIDs {

		e.GUID = guID
		indexKeys := e.IndexKeys()[m.EntryIndexGUID]

		err := z.db.View(func(tx *bolt.Tx) error {
			if data, ok := kvGetFromIndex(m.EntryEntity, m.EntryIndexGUID, indexKeys, tx); ok {
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

	} // loop

	return result, nil

}
