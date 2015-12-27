package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
	"log"
	m "rakewire/model"
	"strconv"
)

// UserEntrySave saves userentries to the database.
func (z *Service) UserEntrySave(userentries []*m.UserEntry) error {

	z.Lock()
	defer z.Unlock()

	err := z.db.Update(func(tx *bolt.Tx) error {
		for _, userentry := range userentries {
			if err := kvSave(m.UserEntryEntity, userentry, tx); err != nil {
				return err
			}
		}
		return nil
	})

	return err

}

// UserEntryAddNew saves userentries to the database.
func (z *Service) UserEntryAddNew(allEntries []*m.Entry, tx *bolt.Tx) error {

	keyedEntries := sortEntriesByFeed(allEntries)

	for feedID, entries := range keyedEntries {

		// define index keys
		uf := &m.UserFeed{}
		uf.FeedID = feedID
		minKeys := uf.IndexKeys()[m.UserFeedIndexFeed]
		uf.FeedID++
		nxtKeys := uf.IndexKeys()[m.UserFeedIndexFeed]

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserFeedEntity)).Bucket([]byte(m.UserFeedIndexFeed))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {

			log.Printf("%-7s %-7s UserEntryAddNew cursor: %s", logTrace, logName, k)

			userID, err := kvKeyElementID(k, 1)
			if err != nil {
				return err
			}

			log.Printf("%-7s %-7s UserEntryAddNew userID: %d", logTrace, logName, userID)

			for _, entry := range entries {
				userentry := &m.UserEntry{
					UserID:  userID,
					EntryID: entry.ID,
					Updated: entry.Updated,
				}
				log.Printf("%-7s %-7s UserEntryAddNew add: %d", logTrace, logName, entry.ID)
				if err := kvSave(m.UserEntryEntity, userentry, tx); err != nil {
					return err
				}
			}

		}

	}

	return nil

}

// UserEntryGetUnreadForUser retrieves user entries for a user and feed.
func (z *Service) UserEntryGetUnreadForUser(userID uint64) ([]*m.UserEntry, error) {

	var result []*m.UserEntry

	// define index keys
	ue := &m.UserEntry{}
	ue.UserID = userID
	minKeys := ue.IndexKeys()[m.UserEntryIndexRead]
	ue.UserID = userID + 1
	ue.Read = true
	nxtKeys := ue.IndexKeys()[m.UserEntryIndexRead]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexRead))
		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))
		bEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.EntryEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			log.Printf("%-7s %-7s user entries get unread: %s: %s", logDebug, logName, k, v)

			id, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			if data, ok := kvGet(id, bUserEntry); ok {
				ue := &m.UserEntry{}
				if err := ue.Deserialize(data); err != nil {
					return err
				}
				if data, ok := kvGet(ue.EntryID, bEntry); ok {
					e := &m.Entry{}
					if err := e.Deserialize(data); err != nil {
						return err
					}
					ue.Entry = e
					result = append(result, ue)
				}
			}

		}

		return nil

	})

	return result, err

}

func sortEntriesByFeed(entries []*m.Entry) map[uint64][]*m.Entry {

	result := make(map[uint64][]*m.Entry)

	for _, entry := range entries {
		a := result[entry.FeedID]
		a = append(a, entry)
		result[entry.FeedID] = a
	}

	return result

}
