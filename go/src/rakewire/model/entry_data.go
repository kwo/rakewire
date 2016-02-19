package model

import (
	"bytes"
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

		// subscription Feed index = FeedID|UserID : SubscriptionID
		min, max := kvKeyMinMax(feedID)
		bIndex := tx.Bucket(bucketIndex).Bucket(subscriptionEntity).Bucket(subscriptionIndexFeed)
		bSubscription := tx.Bucket(bucketData).Bucket(subscriptionEntity)

		c := bIndex.Cursor()
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

			keys := kvKeyDecode(k)
			userID := keys[1]
			subscriptionID := string(v)

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
func EntryTotalByUser(userID string, tx Transaction) uint {

	var result uint

	// entry User index = UserID|EntryID : EntryID
	min, max := kvKeyMinMax(userID)
	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexUser)

	c := bIndex.Cursor()
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
		result++
	} // loop

	return result

}

// EntriesStarredByUser retrieves starred user items for a user.
func EntriesStarredByUser(userID string, tx Transaction) (Entries, error) {

	entries := Entries{}

	// entry index Star = UserID|IsStar|Updated|EntryID : EntryID
	min := kvKeyEncode(userID, kvKeyBoolEncode(false))
	nxt := kvKeyEncode(userID, kvKeyBoolEncode(true))
	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexStar)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		entryID := string(v)

		if data, ok := kvGet(entryID, bEntry); ok {
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
func EntriesUnreadByUser(userID string, tx Transaction) (Entries, error) {

	entries := Entries{}

	// entry Read index = UserID|IsRead|Updated|EntryID : EntryID
	min := kvKeyEncode(userID, kvKeyBoolEncode(false))
	nxt := kvKeyEncode(userID, kvKeyBoolEncode(true))
	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexRead)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		entryID := string(v)

		if data, ok := kvGet(entryID, bEntry); ok {
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
func EntriesByUser(userID string, ids []string, tx Transaction) (Entries, error) {

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

// EntriesGetAll retrieves the all entries for a user.
func EntriesGetAll(userID string, tx Transaction) (Entries, error) {
	return EntriesGetNext(userID, kvKeyUintEncode(0), 0, tx)
}

// EntriesGetNext retrieves the next X user items for a user.
func EntriesGetNext(userID, minID string, count int, tx Transaction) (Entries, error) {

	entries := Entries{}

	// entry User index = UserID|EntryID : EntryID
	min := kvKeyMax(userID, minID) // minID is exclusive, cursor, inclusive
	max := kvKeyMax(userID)
	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexUser)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		entryID := string(v)

		if data, ok := kvGet(entryID, bEntry); ok {
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
func EntriesGetPrev(userID, maxID string, count int, tx Transaction) (Entries, error) {

	entries := Entries{}

	// maxID is exclusive, cursor, inclusive
	if maxIDUnit64, err := kvKeyUintDecode(maxID); err == nil {
		maxIDUnit64--
		maxID = kvKeyUintEncode(maxIDUnit64)
	} else {
		return nil, err
	}

	// entry User index = UserID|EntryID : EntryID
	min := kvKeyEncode(userID)
	max := kvKeyEncode(userID, maxID) // maxID is exclusive, cursor, inclusive
	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexUser)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	c := bIndex.Cursor()
	seekBack := func(key []byte) ([]byte, []byte) {
		k, v := c.Seek(key)
		if k == nil {
			k, v = c.Prev()
		}
		return k, v
	}

	for k, v := seekBack(max); k != nil && bytes.Compare(k, min) >= 0; k, v = c.Prev() {

		entryID := string(v)

		if data, ok := kvGet(entryID, bEntry); ok {
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

// EntriesUpdateReadByFeed updated the read flag of user items.
// MaxTime prevents marking new, not yet synced, entries from being marked as read.
func EntriesUpdateReadByFeed(userID, subscriptionID string, maxTime time.Time, read bool, tx Transaction) error {

	var idCache []string

	// entry Read index = UserID|IsRead|Updated|EntryID : EntryID
	min := kvKeyEncode(userID, kvKeyBoolEncode(!read))
	max := kvKeyEncode(userID, kvKeyBoolEncode(!read), kvKeyTimeEncode(maxTime))
	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexRead)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		entryID := string(v)
		idCache = append(idCache, entryID)
	} // cursor

	for _, id := range idCache {
		if data, ok := kvGet(id, bEntry); ok {

			entry := &Entry{}
			if err := entry.deserialize(data); err != nil {
				return err
			}

			// subscriptionID 0: denotes "Kindling" or all subscriptions
			if subscriptionID == kvKeyUintEncode(0) || entry.SubscriptionID == subscriptionID { // additional filter not in index
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
func EntriesUpdateStarByFeed(userID, subscriptionID string, maxTime time.Time, star bool, tx Transaction) error {

	var idCache []string

	// entry index Star = UserID|IsStar|Updated|EntryID : EntryID
	min := kvKeyEncode(userID, kvKeyBoolEncode(!star))
	max := kvKeyEncode(userID, kvKeyBoolEncode(!star), kvKeyTimeEncode(maxTime))
	bIndex := tx.Bucket(bucketIndex).Bucket(entryEntity).Bucket(entryIndexStar)
	bEntry := tx.Bucket(bucketData).Bucket(entryEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		entryID := string(v)
		idCache = append(idCache, entryID)
	} // cursor

	for _, id := range idCache {
		if data, ok := kvGet(id, bEntry); ok {

			entry := &Entry{}
			if err := entry.deserialize(data); err != nil {
				return err
			}

			// subscriptionID 0: denotes "Kindling" or all subscriptions
			if subscriptionID == kvKeyUintEncode(0) || entry.SubscriptionID == subscriptionID { // additional filter not in index
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
func EntriesUpdateReadByGroup(userID, groupID string, maxTime time.Time, read bool, tx Transaction) error {

	subscriptions, err := SubscriptionsByUser(userID, tx)
	if err != nil {
		return err
	}
	for _, uf := range subscriptions {
		if groupID == kvKeyUintEncode(0) || uf.HasGroup(groupID) {
			if err := EntriesUpdateReadByFeed(userID, uf.ID, maxTime, read, tx); err != nil {
				return err
			}
		}
	}

	return nil

}

// EntriesUpdateStarByGroup updated the star flag of user items
func EntriesUpdateStarByGroup(userID, groupID string, maxTime time.Time, star bool, tx Transaction) error {

	subscriptions, err := SubscriptionsByUser(userID, tx)
	if err != nil {
		return err
	}
	for _, uf := range subscriptions {
		if groupID == kvKeyUintEncode(0) || uf.HasGroup(groupID) {
			if err := EntriesUpdateStarByFeed(userID, uf.ID, maxTime, star, tx); err != nil {
				return err
			}
		}
	}

	return nil

}
