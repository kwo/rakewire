package model

import (
	"bytes"
	"strconv"
)

// EntriesByGUIDs retrieves entries for specific GUIDs
func EntriesByGUIDs(feedID uint64, guIDs []string, tx Transaction) (map[string]*Entry, error) {

	entries := make(map[string]*Entry)

	e := &Entry{}
	e.FeedID = feedID

	for _, guID := range guIDs {

		e.GUID = guID
		indexKeys := e.IndexKeys()[EntryIndexGUID]

		if data, ok := kvGetFromIndex(EntryEntity, EntryIndexGUID, indexKeys, tx); ok {
			entry := &Entry{}
			if err := entry.Deserialize(data); err != nil {
				return nil, err
			}
			entries[entry.GUID] = entry
		}

	} // loop

	return entries, nil

}

// EntriesByFeed retrieves entries for the given feed
func EntriesByFeed(feedID uint64, tx Transaction) ([]*Entry, error) {

	entries := []*Entry{}

	// define index keys
	e := &Entry{}
	e.FeedID = feedID
	minKeys := e.IndexKeys()[EntryIndexGUID]
	e.FeedID = feedID + 1
	nxtKeys := e.IndexKeys()[EntryIndexGUID]

	bIndex := tx.Bucket(bucketIndex).Bucket(EntryEntity).Bucket(EntryIndexGUID)
	bEntry := tx.Bucket(bucketData).Bucket(EntryEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		//log.Printf("%-7s %-7s user entries get unread: %s: %s", logDebug, logName, k, v)

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bEntry); ok {
			entry := &Entry{}
			if err := entry.Deserialize(data); err != nil {
				return nil, err
			}
			entries = append(entries, entry)
		}

	} // loop

	return entries, nil

}

// Delete removes an entry from the database.
func (entry *Entry) Delete(tx Transaction) error {
	return kvSave(EntryEntity, entry, tx)
}
