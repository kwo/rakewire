package model

import (
	"bytes"
	"strconv"
	"time"
)

// EntriesSave saves entries to the database.
func EntriesSave(entries Entries, tx Transaction) error {

	for _, entry := range entries {
		if err := kvSave(entryEntity, entry, tx); err != nil {
			return err
		}
	}

	return nil

}

// EntriesAddNew saves entries to the database.
func EntriesAddNew(allItems Items, tx Transaction) error {

	keyedItems := allItems.GroupAllByFeedID()

	for feedID, items := range keyedItems {

		// define index keys
		uf := &Subscription{}
		uf.FeedID = feedID
		minKeys := uf.indexKeys()[subscriptionIndexFeed]
		uf.FeedID++
		nxtKeys := uf.indexKeys()[subscriptionIndexFeed]

		bIndex := tx.Bucket(bucketIndex).Bucket(subscriptionEntity).Bucket(subscriptionIndexFeed)
		bSubscription := tx.Bucket(bucketData).Bucket(subscriptionEntity)

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			userID, err := kvKeyElementID(k, 1)
			if err != nil {
				return err
			}

			subscriptionID, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			subscription := &Subscription{}
			if data, ok := kvGet(subscriptionID, bSubscription); ok {
				if err := subscription.deserialize(data); err != nil {
					return err
				}
			}

			for _, item := range items {
				entry := &Entry{
					UserID:         userID,
					SubscriptionID: subscriptionID,
					ItemID:         item.ID,
					Updated:        item.Updated,
					IsRead:         subscription.AutoRead,
					IsStar:         subscription.AutoStar,
				}
				if err := kvSave(entryEntity, entry, tx); err != nil {
					return err
				}
			}

		}

	}

	return nil

}

// EntryTotalByUser retrieves the total count of items for the user.
func EntryTotalByUser(userID uint64, tx Transaction) uint {

	var result uint

	// define index keys
	ue := &Entry{}
	ue.UserID = userID
	minKeys := ue.indexKeys()[entryIndexUser]
	ue.UserID = userID + 1
	nxtKeys := ue.indexKeys()[entryIndexUser]

	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexUser)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
		result++
	} // loop

	return result

}

// EntriesStarredByUser retrieves starred user items for a user.
func EntriesStarredByUser(userID uint64, tx Transaction) (Entries, error) {

	entries := Entries{}

	// define index keys
	ue := &Entry{}
	ue.UserID = userID
	ue.IsStar = true
	minKeys := ue.indexKeys()[entryIndexStar]
	ue.UserID = userID + 1
	ue.IsStar = false
	nxtKeys := ue.indexKeys()[entryIndexStar]

	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexStar)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	// log.Printf("%-7s %-7s user items get unread min: %s", logDebug, logName, min)
	// log.Printf("%-7s %-7s user items get unread max: %s", logDebug, logName, nxt)

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		//log.Printf("%-7s %-7s user items get unread: %s: %s", logDebug, logName, k, v)

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bEntry); ok {
			entry := &Entry{}
			if err := entry.deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(entry.ItemID, bItem); ok {
				item := &Item{}
				if err := item.deserialize(data); err != nil {
					return nil, err
				}
				entry.Item = item
				entries = append(entries, entry)
			}
		}

	}

	return entries, nil

}

// EntriesUnreadByUser retrieves unread user items for a user.
func EntriesUnreadByUser(userID uint64, tx Transaction) (Entries, error) {

	entries := Entries{}

	// define index keys
	ue := &Entry{}
	ue.UserID = userID
	minKeys := ue.indexKeys()[entryIndexRead]
	ue.IsRead = true
	nxtKeys := ue.indexKeys()[entryIndexRead]

	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexRead)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	//log.Printf("%-7s %-7s user items get unread min: %s", logDebug, logName, min)
	//log.Printf("%-7s %-7s user items get unread max: %s", logDebug, logName, nxt)

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		//log.Printf("%-7s %-7s user items get unread: %s: %s", logDebug, logName, k, v)

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bEntry); ok {
			entry := &Entry{}
			if err := entry.deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(entry.ItemID, bItem); ok {
				item := &Item{}
				if err := item.deserialize(data); err != nil {
					return nil, err
				}
				entry.Item = item
				entries = append(entries, entry)
			}
		}

	}

	return entries, nil

}

// EntriesByUser retrieves the specific user items for a user.
func EntriesByUser(userID uint64, ids []uint64, tx Transaction) (Entries, error) {

	entries := Entries{}

	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	for _, id := range ids {

		if data, ok := kvGet(id, bEntry); ok {
			entry := &Entry{}
			if err := entry.deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(entry.ItemID, bItem); ok {
				item := &Item{}
				if err := item.deserialize(data); err != nil {
					return nil, err
				}
				entry.Item = item
				entries = append(entries, entry)
			}
		}

	}

	return entries, nil

}

// EntriesGetNext retrieves the next X user items for a user.
func EntriesGetNext(userID uint64, minID uint64, count int, tx Transaction) (Entries, error) {

	entries := Entries{}

	// define index keys
	ue := &Entry{}
	ue.UserID = userID
	ue.ID = minID + 1 // minID is exclusive, cursor, inclusive
	minKeys := ue.indexKeys()[entryIndexUser]
	ue.UserID = userID + 1
	nxtKeys := ue.indexKeys()[entryIndexUser]

	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexUser)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		//log.Printf("%-7s %-7s user items get unread: %s: %s", logDebug, logName, k, v)

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bEntry); ok {
			entry := &Entry{}
			if err := entry.deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(entry.ItemID, bItem); ok {
				item := &Item{}
				if err := item.deserialize(data); err != nil {
					return nil, err
				}
				entry.Item = item
				entries = append(entries, entry)
			}
		}

		if count > 0 && len(entries) >= count {
			break
		}

	} // loop

	return entries, nil

}

// EntriesGetPrev retrieves the previous X user items for a user.
func EntriesGetPrev(userID uint64, maxID uint64, count int, tx Transaction) (Entries, error) {

	entries := Entries{}

	// define index keys
	ue := &Entry{}
	ue.UserID = userID
	ue.ID = maxID - 1 // maxID is exclusive, cursor, inclusive
	maxKeys := ue.indexKeys()[entryIndexUser]
	ue.ID = 0
	minKeys := ue.indexKeys()[entryIndexUser]

	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexUser)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	c := bIndex.Cursor()
	// for k, v := c.First(); k != nil; k, v = c.Next() {
	// 	log.Printf("%-7s %-7s user items get prev: %s: %s", logDebug, logName, k, v)
	// }

	min := []byte(kvKeys(minKeys))
	max := []byte(kvKeys(maxKeys))
	// log.Printf("%-7s %-7s user items get prev min: %s", logDebug, logName, min)
	// log.Printf("%-7s %-7s user items get prev max: %s", logDebug, logName, max)

	seekBack := func(key []byte) ([]byte, []byte) {
		k, v := c.Seek(key)
		if k == nil {
			k, v = c.Prev()
		}
		return k, v
	}

	for k, v := seekBack(max); k != nil && bytes.Compare(k, min) >= 0; k, v = c.Prev() {

		//log.Printf("%-7s %-7s user items get prev: %s: %s", logDebug, logName, k, v)

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, bEntry); ok {
			entry := &Entry{}
			if err := entry.deserialize(data); err != nil {
				return nil, err
			}
			if data, ok := kvGet(entry.ItemID, bItem); ok {
				item := &Item{}
				if err := item.deserialize(data); err != nil {
					return nil, err
				}
				entry.Item = item
				entries = append(entries, entry)
			}
		}

		if count > 0 && len(entries) >= count {
			break
		}

	} // loop

	return entries, nil

}

// EntriesUpdateReadByFeed updated the read flag of user items
func EntriesUpdateReadByFeed(userID, subscriptionID uint64, maxTime time.Time, read bool, tx Transaction) error {

	var idCache []uint64

	// define index keys
	ue := &Entry{}
	ue.UserID = userID
	ue.IsRead = !read
	minKeys := ue.indexKeys()[entryIndexRead]
	ue.Updated = maxTime.Add(1 * time.Second).Truncate(time.Second)
	nxtKeys := ue.indexKeys()[entryIndexRead]

	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexRead)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	//log.Printf("%-7s %-7s EntryUpdateReadByFeed min: %s", logDebug, logName, min)
	//log.Printf("%-7s %-7s EntryUpdateReadByFeed max: %s", logDebug, logName, nxt)
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
		//log.Printf("%-7s %-7s EntryUpdateReadByFeed %s: %s", logTrace, logName, k, v)

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return err
		}

		idCache = append(idCache, id)

	} // cursor

	for _, id := range idCache {
		if data, ok := kvGet(id, bEntry); ok {

			entry := &Entry{}
			if err := entry.deserialize(data); err != nil {
				return err
			}

			if subscriptionID == 0 || entry.SubscriptionID == subscriptionID { // additional filter not in index
				entry.IsRead = read
				if err := kvSave(entryEntity, entry, tx); err != nil {
					return err
				}
			}

		}
	}

	return nil

}

// EntriesUpdateStarByFeed updated the star flag of user items
func EntriesUpdateStarByFeed(userID, subscriptionID uint64, maxTime time.Time, star bool, tx Transaction) error {

	var idCache []uint64

	// define index keys
	ue := &Entry{}
	ue.UserID = userID
	ue.IsStar = !star
	minKeys := ue.indexKeys()[entryIndexStar]
	ue.Updated = maxTime.Add(1 * time.Second).Truncate(time.Second)
	nxtKeys := ue.indexKeys()[entryIndexStar]

	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexStar)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	//log.Printf("%-7s %-7s EntryUpdateReadByFeed min: %s", logDebug, logName, min)
	//log.Printf("%-7s %-7s EntryUpdateReadByFeed max: %s", logDebug, logName, nxt)
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
		//log.Printf("%-7s %-7s EntryUpdateReadByFeed %s: %s", logTrace, logName, k, v)

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return err
		}

		idCache = append(idCache, id)

	} // cursor

	for _, id := range idCache {
		if data, ok := kvGet(id, bEntry); ok {

			entry := &Entry{}
			if err := entry.deserialize(data); err != nil {
				return err
			}

			if subscriptionID == 0 || entry.SubscriptionID == subscriptionID { // additional filter not in index
				entry.IsStar = star
				if err := kvSave(entryEntity, entry, tx); err != nil {
					return err
				}
			}

		}
	}

	return nil

}

// EntriesUpdateReadByGroup updated the read flag of user items
func EntriesUpdateReadByGroup(userID, groupID uint64, maxTime time.Time, read bool, tx Transaction) error {

	subscriptions, err := SubscriptionsByUser(userID, tx)
	if err != nil {
		return err
	}
	for _, uf := range subscriptions {
		if groupID == 0 || uf.HasGroup(groupID) {
			if err := EntriesUpdateReadByFeed(userID, uf.ID, maxTime, read, tx); err != nil {
				return err
			}
		}
	}

	return nil

}

// EntriesUpdateStarByGroup updated the star flag of user items
func EntriesUpdateStarByGroup(userID, groupID uint64, maxTime time.Time, star bool, tx Transaction) error {

	subscriptions, err := SubscriptionsByUser(userID, tx)
	if err != nil {
		return err
	}
	for _, uf := range subscriptions {
		if groupID == 0 || uf.HasGroup(groupID) {
			if err := EntriesUpdateStarByFeed(userID, uf.ID, maxTime, star, tx); err != nil {
				return err
			}
		}
	}

	return nil

}
