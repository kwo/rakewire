package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"strconv"
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

// EntriesGet retrieves entries for the given feed
func (z *Service) EntriesGet(feedID uint64) ([]*m.Entry, error) {

	result := []*m.Entry{}

	// define index keys
	e := &m.Entry{}
	e.FeedID = feedID
	minKeys := e.IndexKeys()[m.EntryIndexGUID]
	e.FeedID = feedID + 1
	nxtKeys := e.IndexKeys()[m.EntryIndexGUID]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.EntryEntity)).Bucket([]byte(m.EntryIndexGUID))
		bEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.EntryEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			//log.Printf("%-7s %-7s user entries get unread: %s: %s", logDebug, logName, k, v)

			id, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			if data, ok := kvGet(id, bEntry); ok {
				e := &m.Entry{}
				if err := e.Deserialize(data); err != nil {
					return err
				}
				result = append(result, e)
			}

		} // loop

		return nil

	})

	return result, err

}

// EntryDelete removes an entry from the database.
func (z *Service) EntryDelete(entry *m.Entry) error {
	z.Lock()
	defer z.Unlock()
	return z.db.Update(func(tx *bolt.Tx) error {
		return kvSave(m.EntryEntity, entry, tx)
	})
}
