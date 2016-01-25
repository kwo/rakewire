package model

import (
	"bytes"
	"strconv"
	"time"
)

// UserEntriesSave saves userentries to the database.
func UserEntriesSave(userentries []*UserEntry, tx Transaction) error {

	for _, userentry := range userentries {
		if err := kvSave(UserEntryEntity, userentry, tx); err != nil {
			return err
		}
	}

	return nil

}

// UserEntriesAddNew saves userentries to the database.
func UserEntriesAddNew(allEntries []*Entry, tx Transaction) error {

	keyedEntries := groupEntriesByFeed(allEntries)

	for feedID, entries := range keyedEntries {

		// define index keys
		uf := &UserFeed{}
		uf.FeedID = feedID
		minKeys := uf.IndexKeys()[UserFeedIndexFeed]
		uf.FeedID++
		nxtKeys := uf.IndexKeys()[UserFeedIndexFeed]

		bIndex := tx.Bucket(bucketIndex).Bucket(UserFeedEntity).Bucket(UserFeedIndexFeed)
		bUserFeed := tx.Bucket(bucketData).Bucket(UserFeedEntity)

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

			userfeed := &UserFeed{}
			if data, ok := kvGet(userFeedID, bUserFeed); ok {
				if err := userfeed.Deserialize(data); err != nil {
					return err
				}
			}

			for _, entry := range entries {
				userentry := &UserEntry{
					UserID:     userID,
					UserFeedID: userFeedID,
					EntryID:    entry.ID,
					Updated:    entry.Updated,
					IsRead:     userfeed.AutoRead,
					IsStar:     userfeed.AutoStar,
				}
				if err := kvSave(UserEntryEntity, userentry, tx); err != nil {
					return err
				}
			}

		}

	}

	return nil

}

// UserEntryTotalByUser retrieves the total count of entries for the user.
func UserEntryTotalByUser(userID uint64, tx Transaction) uint {

	var result uint

	// define index keys
	ue := &UserEntry{}
	ue.UserID = userID
	minKeys := ue.IndexKeys()[UserEntryIndexUser]
	ue.UserID = userID + 1
	nxtKeys := ue.IndexKeys()[UserEntryIndexUser]

	bIndex := tx.Bucket(bucketIndex).Bucket(UserEntryEntity).Bucket(UserEntryIndexUser)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
		result++
	} // loop

	return result

}

// UserEntriesStarredByUser retrieves starred user entries for a user.
func UserEntriesStarredByUser(userID uint64, tx Transaction) ([]*UserEntry, error) {

	var userentries []*UserEntry

	// define index keys
	ue := &UserEntry{}
	ue.UserID = userID
	ue.IsStar = true
	minKeys := ue.IndexKeys()[UserEntryIndexStar]
	ue.UserID = userID + 1
	ue.IsStar = false
	nxtKeys := ue.IndexKeys()[UserEntryIndexStar]

	bIndex := tx.Bucket(bucketIndex).Bucket(UserEntryEntity).Bucket(UserEntryIndexStar)
	bUserEntry := tx.Bucket(bucketData).Bucket(UserEntryEntity)
	bEntry := tx.Bucket(bucketData).Bucket(EntryEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	// log.Printf("%-7s %-7s user entries get unread min: %s", logDebug, logName, min)
	// log.Printf("%-7s %-7s user entries get unread max: %s", logDebug, logName, nxt)

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		//log.Printf("%-7s %-7s user entries get unread: %s: %s", logDebug, logName, k, v)

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bUserEntry); ok {
			userentry := &UserEntry{}
			if err := userentry.Deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(userentry.EntryID, bEntry); ok {
				entry := &Entry{}
				if err := entry.Deserialize(data); err != nil {
					return nil, err
				}
				userentry.Entry = entry
				userentries = append(userentries, userentry)
			}
		}

	}

	return userentries, nil

}

// UserEntriesUnreadByUser retrieves unread user entries for a user.
func UserEntriesUnreadByUser(userID uint64, tx Transaction) ([]*UserEntry, error) {

	var userentries []*UserEntry

	// define index keys
	ue := &UserEntry{}
	ue.UserID = userID
	minKeys := ue.IndexKeys()[UserEntryIndexRead]
	ue.IsRead = true
	nxtKeys := ue.IndexKeys()[UserEntryIndexRead]

	bIndex := tx.Bucket(bucketIndex).Bucket(UserEntryEntity).Bucket(UserEntryIndexRead)
	bUserEntry := tx.Bucket(bucketData).Bucket(UserEntryEntity)
	bEntry := tx.Bucket(bucketData).Bucket(EntryEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	//log.Printf("%-7s %-7s user entries get unread min: %s", logDebug, logName, min)
	//log.Printf("%-7s %-7s user entries get unread max: %s", logDebug, logName, nxt)

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		//log.Printf("%-7s %-7s user entries get unread: %s: %s", logDebug, logName, k, v)

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bUserEntry); ok {
			userentry := &UserEntry{}
			if err := userentry.Deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(userentry.EntryID, bEntry); ok {
				entry := &Entry{}
				if err := entry.Deserialize(data); err != nil {
					return nil, err
				}
				userentry.Entry = entry
				userentries = append(userentries, userentry)
			}
		}

	}

	return userentries, nil

}

// UserEntriesByUser retrieves the specific user entries for a user.
func UserEntriesByUser(userID uint64, ids []uint64, tx Transaction) ([]*UserEntry, error) {

	var userentries []*UserEntry

	bUserEntry := tx.Bucket(bucketData).Bucket(UserEntryEntity)
	bEntry := tx.Bucket(bucketData).Bucket(EntryEntity)

	for _, id := range ids {

		if data, ok := kvGet(id, bUserEntry); ok {
			userentry := &UserEntry{}
			if err := userentry.Deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(userentry.EntryID, bEntry); ok {
				entry := &Entry{}
				if err := entry.Deserialize(data); err != nil {
					return nil, err
				}
				userentry.Entry = entry
				userentries = append(userentries, userentry)
			}
		}

	}

	return userentries, nil

}

// UserEntryGetNext retrieves the next X user entries for a user.
func UserEntryGetNext(userID uint64, minID uint64, count int, tx Transaction) ([]*UserEntry, error) {

	var userentries []*UserEntry

	// define index keys
	ue := &UserEntry{}
	ue.UserID = userID
	ue.ID = minID + 1 // minID is exclusive, cursor, inclusive
	minKeys := ue.IndexKeys()[UserEntryIndexUser]
	ue.UserID = userID + 1
	nxtKeys := ue.IndexKeys()[UserEntryIndexUser]

	bIndex := tx.Bucket(bucketIndex).Bucket(UserEntryEntity).Bucket(UserEntryIndexUser)
	bUserEntry := tx.Bucket(bucketData).Bucket(UserEntryEntity)
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

		if data, ok := kvGet(id, bUserEntry); ok {
			userentry := &UserEntry{}
			if err := userentry.Deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(userentry.EntryID, bEntry); ok {
				entry := &Entry{}
				if err := entry.Deserialize(data); err != nil {
					return nil, err
				}
				userentry.Entry = entry
				userentries = append(userentries, userentry)
			}
		}

		if count > 0 && len(userentries) >= count {
			break
		}

	} // loop

	return userentries, nil

}

// UserEntryGetPrev retrieves the previous X user entries for a user.
func UserEntryGetPrev(userID uint64, maxID uint64, count int, tx Transaction) ([]*UserEntry, error) {

	var userentries []*UserEntry

	// define index keys
	ue := &UserEntry{}
	ue.UserID = userID
	ue.ID = maxID - 1 // maxID is exclusive, cursor, inclusive
	maxKeys := ue.IndexKeys()[UserEntryIndexUser]
	ue.ID = 0
	minKeys := ue.IndexKeys()[UserEntryIndexUser]

	bIndex := tx.Bucket(bucketIndex).Bucket(UserEntryEntity).Bucket(UserEntryIndexUser)
	bUserEntry := tx.Bucket(bucketData).Bucket(UserEntryEntity)
	bEntry := tx.Bucket(bucketData).Bucket(EntryEntity)

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
			return nil, err
		}

		if data, ok := kvGet(id, bUserEntry); ok {
			userentry := &UserEntry{}
			if err := userentry.Deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(userentry.EntryID, bEntry); ok {
				entry := &Entry{}
				if err := entry.Deserialize(data); err != nil {
					return nil, err
				}
				userentry.Entry = entry
				userentries = append(userentries, userentry)
			}
		}

		if count > 0 && len(userentries) >= count {
			break
		}

	} // loop

	return userentries, nil

}

// UserEntriesUpdateReadByFeed updated the read flag of user entries
func UserEntriesUpdateReadByFeed(userID, userFeedID uint64, maxTime time.Time, read bool, tx Transaction) error {

	var idCache []uint64

	// define index keys
	ue := &UserEntry{}
	ue.UserID = userID
	ue.IsRead = !read
	minKeys := ue.IndexKeys()[UserEntryIndexRead]
	ue.Updated = maxTime.Add(1 * time.Second).Truncate(time.Second)
	nxtKeys := ue.IndexKeys()[UserEntryIndexRead]

	bIndex := tx.Bucket(bucketIndex).Bucket(UserEntryEntity).Bucket(UserEntryIndexRead)
	bUserEntry := tx.Bucket(bucketData).Bucket(UserEntryEntity)

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

			userentry := &UserEntry{}
			if err := userentry.Deserialize(data); err != nil {
				return err
			}

			if userFeedID == 0 || userentry.UserFeedID == userFeedID { // additional filter not in index
				userentry.IsRead = read
				if err := kvSave(UserEntryEntity, userentry, tx); err != nil {
					return err
				}
			}

		}
	}

	return nil

}

// UserEntriesUpdateStarByFeed updated the star flag of user entries
func UserEntriesUpdateStarByFeed(userID, userFeedID uint64, maxTime time.Time, star bool, tx Transaction) error {

	var idCache []uint64

	// define index keys
	ue := &UserEntry{}
	ue.UserID = userID
	ue.IsStar = !star
	minKeys := ue.IndexKeys()[UserEntryIndexStar]
	ue.Updated = maxTime.Add(1 * time.Second).Truncate(time.Second)
	nxtKeys := ue.IndexKeys()[UserEntryIndexStar]

	bIndex := tx.Bucket(bucketIndex).Bucket(UserEntryEntity).Bucket(UserEntryIndexStar)
	bUserEntry := tx.Bucket(bucketData).Bucket(UserEntryEntity)

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

			userentry := &UserEntry{}
			if err := userentry.Deserialize(data); err != nil {
				return err
			}

			if userFeedID == 0 || userentry.UserFeedID == userFeedID { // additional filter not in index
				userentry.IsStar = star
				if err := kvSave(UserEntryEntity, userentry, tx); err != nil {
					return err
				}
			}

		}
	}

	return nil

}

// UserEntriesUpdateReadByGroup updated the read flag of user entries
func UserEntriesUpdateReadByGroup(userID, groupID uint64, maxTime time.Time, read bool, tx Transaction) error {

	userfeeds, err := UserFeedsByUser(userID, tx)
	if err != nil {
		return err
	}
	for _, uf := range userfeeds {
		if groupID == 0 || uf.HasGroup(groupID) {
			if err := UserEntriesUpdateReadByFeed(userID, uf.ID, maxTime, read, tx); err != nil {
				return err
			}
		}
	}

	return nil

}

// UserEntriesUpdateStarByGroup updated the star flag of user entries
func UserEntriesUpdateStarByGroup(userID, groupID uint64, maxTime time.Time, star bool, tx Transaction) error {

	userfeeds, err := UserFeedsByUser(userID, tx)
	if err != nil {
		return err
	}
	for _, uf := range userfeeds {
		if groupID == 0 || uf.HasGroup(groupID) {
			if err := UserEntriesUpdateStarByFeed(userID, uf.ID, maxTime, star, tx); err != nil {
				return err
			}
		}
	}

	return nil

}

func groupEntriesByFeed(entries []*Entry) map[uint64][]*Entry {

	result := make(map[uint64][]*Entry)

	for _, entry := range entries {
		a := result[entry.FeedID]
		a = append(a, entry)
		result[entry.FeedID] = a
	}

	return result

}
