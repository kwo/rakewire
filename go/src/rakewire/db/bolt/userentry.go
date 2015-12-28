package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
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

	keyedEntries := groupEntriesByFeed(allEntries)

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
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			userID, err := kvKeyElementID(k, 1)
			if err != nil {
				return err
			}

			// TODO: get whole userfeed, add option to mark entries as read when new

			userFeedID, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			for _, entry := range entries {
				userentry := &m.UserEntry{
					UserID:     userID,
					UserFeedID: userFeedID,
					EntryID:    entry.ID,
					Updated:    entry.Updated,
				}
				if err := kvSave(m.UserEntryEntity, userentry, tx); err != nil {
					return err
				}
			}

		}

	}

	return nil

}

// UserEntryGetUnreadForUser retrieves unread user entries for a user.
func (z *Service) UserEntryGetUnreadForUser(userID uint64) ([]*m.UserEntry, error) {

	var result []*m.UserEntry

	// define index keys
	ue := &m.UserEntry{}
	ue.UserID = userID
	minKeys := ue.IndexKeys()[m.UserEntryIndexRead]
	ue.UserID = userID + 1
	ue.IsRead = true
	nxtKeys := ue.IndexKeys()[m.UserEntryIndexRead]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexRead))
		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))
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

// UserEntryGetTotalForUser retrieves the total count of entries for the user.
func (z *Service) UserEntryGetTotalForUser(userID uint64) (uint, error) {

	var result uint

	// define index keys
	ue := &m.UserEntry{}
	ue.UserID = userID
	minKeys := ue.IndexKeys()[m.UserEntryIndexUser]
	ue.UserID = userID + 1
	nxtKeys := ue.IndexKeys()[m.UserEntryIndexUser]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexUser))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
			result++
		} // loop

		return nil

	})

	return result, err

}

// UserEntryGet retrieves the specific user entries for a user.
func (z *Service) UserEntryGet(userID uint64, ids []uint64) ([]*m.UserEntry, error) {

	var result []*m.UserEntry

	err := z.db.View(func(tx *bolt.Tx) error {

		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))
		bEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.EntryEntity))

		for _, id := range ids {

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

// UserEntryGetNext retrieves the next X user entries for a user.
func (z *Service) UserEntryGetNext(userID uint64, minID uint64, count int) ([]*m.UserEntry, error) {

	var result []*m.UserEntry

	// define index keys
	ue := &m.UserEntry{}
	ue.UserID = userID
	ue.ID = minID
	minKeys := ue.IndexKeys()[m.UserEntryIndexUser]
	ue.UserID = userID + 1
	nxtKeys := ue.IndexKeys()[m.UserEntryIndexUser]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexUser))
		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))
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

			if count > 0 && len(result) >= count {
				break
			}

		} // loop

		return nil

	})

	return result, err

}

// UserEntryGetPrev retrieves the previous X user entries for a user.
func (z *Service) UserEntryGetPrev(userID uint64, maxID uint64, count int) ([]*m.UserEntry, error) {

	var result []*m.UserEntry

	// define index keys
	ue := &m.UserEntry{}
	ue.UserID = userID
	ue.ID = maxID
	maxKeys := ue.IndexKeys()[m.UserEntryIndexUser]
	ue.ID = 0
	minKeys := ue.IndexKeys()[m.UserEntryIndexUser]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexUser))
		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))
		bEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.EntryEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		max := []byte(kvKeys(maxKeys))
		for k, v := c.Seek(max); k != nil && bytes.Compare(k, min) >= 0; k, v = c.Prev() {

			//log.Printf("%-7s %-7s user entries get unread: %s: %s", logDebug, logName, k, v)

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

			if count > 0 && len(result) >= count {
				break
			}

		} // loop

		return nil

	})

	return result, err

}

func groupEntriesByFeed(entries []*m.Entry) map[uint64][]*m.Entry {

	result := make(map[uint64][]*m.Entry)

	for _, entry := range entries {
		a := result[entry.FeedID]
		a = append(a, entry)
		result[entry.FeedID] = a
	}

	return result

}
