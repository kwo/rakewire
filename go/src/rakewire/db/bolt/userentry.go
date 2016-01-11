package bolt

import (
	"bytes"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"strconv"
	"time"
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

// UserEntryAdd saves userentries to the database.
func (z *Service) UserEntryAdd(allEntries []*m.Entry) error {
	z.Lock()
	defer z.Unlock()
	return z.db.Update(func(tx *bolt.Tx) error {
		return z.UserEntryAddNew(allEntries, tx)
	})
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
		bUserFeed := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserFeedEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			userID, err := kvKeyElementID(k, 1)
			if err != nil {
				return err
			}

			userFeedID, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			userfeed := &m.UserFeed{}
			if data, ok := kvGet(userFeedID, bUserFeed); ok {
				if err := userfeed.Deserialize(data); err != nil {
					return err
				}
			}

			for _, entry := range entries {
				userentry := &m.UserEntry{
					UserID:     userID,
					UserFeedID: userFeedID,
					EntryID:    entry.ID,
					Updated:    entry.Updated,
					IsRead:     userfeed.AutoRead,
					IsStar:     userfeed.AutoStar,
				}
				if err := kvSave(m.UserEntryEntity, userentry, tx); err != nil {
					return err
				}
			}

		}

	}

	return nil

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

// UserEntryGetStarredForUser retrieves starred user entries for a user.
func (z *Service) UserEntryGetStarredForUser(userID uint64) ([]*m.UserEntry, error) {

	var result []*m.UserEntry

	// define index keys
	ue := &m.UserEntry{}
	ue.UserID = userID
	ue.IsStar = true
	minKeys := ue.IndexKeys()[m.UserEntryIndexStar]
	ue.UserID = userID + 1
	ue.IsStar = false
	nxtKeys := ue.IndexKeys()[m.UserEntryIndexStar]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexStar))
		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))
		bEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.EntryEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		// log.Printf("%-7s %-7s user entries get unread min: %s", logDebug, logName, min)
		// log.Printf("%-7s %-7s user entries get unread max: %s", logDebug, logName, nxt)

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

// UserEntryGetUnreadForUser retrieves unread user entries for a user.
func (z *Service) UserEntryGetUnreadForUser(userID uint64) ([]*m.UserEntry, error) {

	var result []*m.UserEntry

	// define index keys
	ue := &m.UserEntry{}
	ue.UserID = userID
	minKeys := ue.IndexKeys()[m.UserEntryIndexRead]
	ue.IsRead = true
	nxtKeys := ue.IndexKeys()[m.UserEntryIndexRead]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexRead))
		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))
		bEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.EntryEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		//log.Printf("%-7s %-7s user entries get unread min: %s", logDebug, logName, min)
		//log.Printf("%-7s %-7s user entries get unread max: %s", logDebug, logName, nxt)

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

// UserEntryGetByID retrieves the specific user entries for a user.
func (z *Service) UserEntryGetByID(userID uint64, ids []uint64) ([]*m.UserEntry, error) {

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
	ue.ID = minID + 1 // minID is exclusive, cursor, inclusive
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
	ue.ID = maxID - 1 // maxID is exclusive, cursor, inclusive
	maxKeys := ue.IndexKeys()[m.UserEntryIndexUser]
	ue.ID = 0
	minKeys := ue.IndexKeys()[m.UserEntryIndexUser]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexUser))
		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))
		bEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.EntryEntity))

		c := bIndex.Cursor()
		// for k, v := c.First(); k != nil; k, v = c.Next() {
		// 	log.Printf("%-7s %-7s user entries get prev: %s: %s", logDebug, logName, k, v)
		// }

		min := []byte(kvKeys(minKeys))
		max := []byte(kvKeys(maxKeys))
		// log.Printf("%-7s %-7s user entries get prev min: %s", logDebug, logName, min)
		// log.Printf("%-7s %-7s user entries get prev max: %s", logDebug, logName, max)

		seekBack := func(key []byte) ([]byte, []byte) {
			k, v := c.Seek(key)
			if k == nil {
				k, v = c.Prev()
			}
			return k, v
		}

		for k, v := seekBack(max); k != nil && bytes.Compare(k, min) >= 0; k, v = c.Prev() {

			//log.Printf("%-7s %-7s user entries get prev: %s: %s", logDebug, logName, k, v)

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

// UserEntryUpdateReadByFeed updated the read flag of user entries
func (z *Service) UserEntryUpdateReadByFeed(userID, userFeedID uint64, maxTime time.Time, read bool) error {

	z.Lock()
	defer z.Unlock()
	return z.db.Update(func(tx *bolt.Tx) error {

		var idCache []uint64

		// define index keys
		ue := &m.UserEntry{}
		ue.UserID = userID
		ue.IsRead = !read
		minKeys := ue.IndexKeys()[m.UserEntryIndexRead]
		ue.Updated = maxTime.Add(1 * time.Second).Truncate(time.Second)
		nxtKeys := ue.IndexKeys()[m.UserEntryIndexRead]

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexRead))
		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		//log.Printf("%-7s %-7s UserEntryUpdateReadByFeed min: %s", logDebug, logName, min)
		//log.Printf("%-7s %-7s UserEntryUpdateReadByFeed max: %s", logDebug, logName, nxt)
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
			//log.Printf("%-7s %-7s UserEntryUpdateReadByFeed %s: %s", logTrace, logName, k, v)

			id, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			idCache = append(idCache, id)

		} // cursor

		for _, id := range idCache {
			if data, ok := kvGet(id, bUserEntry); ok {

				ue := &m.UserEntry{}
				if err := ue.Deserialize(data); err != nil {
					return err
				}

				if userFeedID == 0 || ue.UserFeedID == userFeedID { // additional filter not in index
					ue.IsRead = read
					if err := kvSave(m.UserEntryEntity, ue, tx); err != nil {
						return err
					}
				}

			}
		}

		return nil

	})

}

// UserEntryUpdateStarByFeed updated the star flag of user entries
func (z *Service) UserEntryUpdateStarByFeed(userID, userFeedID uint64, maxTime time.Time, star bool) error {

	z.Lock()
	defer z.Unlock()
	return z.db.Update(func(tx *bolt.Tx) error {

		var idCache []uint64

		// define index keys
		ue := &m.UserEntry{}
		ue.UserID = userID
		ue.IsStar = !star
		minKeys := ue.IndexKeys()[m.UserEntryIndexStar]
		ue.Updated = maxTime.Add(1 * time.Second).Truncate(time.Second)
		nxtKeys := ue.IndexKeys()[m.UserEntryIndexStar]

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.UserEntryEntity)).Bucket([]byte(m.UserEntryIndexStar))
		bUserEntry := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.UserEntryEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		//log.Printf("%-7s %-7s UserEntryUpdateReadByFeed min: %s", logDebug, logName, min)
		//log.Printf("%-7s %-7s UserEntryUpdateReadByFeed max: %s", logDebug, logName, nxt)
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
			//log.Printf("%-7s %-7s UserEntryUpdateReadByFeed %s: %s", logTrace, logName, k, v)

			id, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			idCache = append(idCache, id)

		} // cursor

		for _, id := range idCache {
			if data, ok := kvGet(id, bUserEntry); ok {

				ue := &m.UserEntry{}
				if err := ue.Deserialize(data); err != nil {
					return err
				}

				if userFeedID == 0 || ue.UserFeedID == userFeedID { // additional filter not in index
					ue.IsStar = star
					if err := kvSave(m.UserEntryEntity, ue, tx); err != nil {
						return err
					}
				}

			}
		}

		return nil

	})

}

// UserEntryUpdateReadByGroup updated the read flag of user entries
func (z *Service) UserEntryUpdateReadByGroup(userID, groupID uint64, maxTime time.Time, read bool) error {

	userfeeds, err := z.UserFeedGetAllByUser(userID)
	if err != nil {
		return err
	}
	for _, uf := range userfeeds {
		if groupID == 0 || uf.HasGroup(groupID) {
			if err := z.UserEntryUpdateReadByFeed(userID, uf.ID, maxTime, read); err != nil {
				return err
			}
		}
	}

	return nil

}

// UserEntryUpdateStarByGroup updated the star flag of user entries
func (z *Service) UserEntryUpdateStarByGroup(userID, groupID uint64, maxTime time.Time, star bool) error {

	userfeeds, err := z.UserFeedGetAllByUser(userID)
	if err != nil {
		return err
	}
	for _, uf := range userfeeds {
		if groupID == 0 || uf.HasGroup(groupID) {
			if err := z.UserEntryUpdateStarByFeed(userID, uf.ID, maxTime, star); err != nil {
				return err
			}
		}
	}

	return nil

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
