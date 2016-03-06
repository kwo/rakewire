package modelng

import (
	"bytes"
	"time"
)

// E groups all entry database methods
var E = &entryStore{}

type entryStore struct{}

func (z *entryStore) AddItems(tx Transaction, allItems Items) error {

	mappedItems := allItems.GroupByFeedID()

	for feedID, items := range mappedItems {
		subscriptions := S.GetForFeed(tx, feedID)
		for _, subscription := range subscriptions {
			for _, item := range items {
				entry := z.New(subscription.UserID, item.ID)
				entry.Updated = item.Updated
				entry.Read = subscription.AutoRead
				entry.Star = subscription.AutoStar
				if err := z.Save(tx, entry); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (z *entryStore) Count(tx Transaction, userID string, minTime, maxTime time.Time) int {

	result := 0

	// index Entry Updated = UserID|Updated|ItemID : ItemID
	min := []byte(keyEncode(userID, keyEncodeTime(minTime)))
	max := []byte(keyEncode(userID, keyEncodeTime(maxTime)))
	c := tx.Bucket(bucketIndex, entityEntry, indexEntryUpdated).Cursor()
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
		result++
	}

	return result

}

func (z *entryStore) Delete(tx Transaction, id string) error {
	return delete(tx, entityEntry, id)
}

func (z *entryStore) Get(tx Transaction, id ...string) *Entry {
	compoundID := ""
	switch len(id) {
	case 1:
		compoundID = id[0]
	case 2:
		compoundID = keyEncode(id...)
	default:
		return nil
	}
	// bucket Entry key = UserID|ItemID
	bData := tx.Bucket(bucketData, entityEntry)
	if data := bData.Get([]byte(compoundID)); data != nil {
		entry := &Entry{}
		if err := entry.decode(data); err == nil {
			return entry
		}
	}
	return nil
}

func (z *entryStore) New(userID, itemID string) *Entry {
	return &Entry{
		UserID: userID,
		ItemID: itemID,
	}
}

func (z *entryStore) Range(tx Transaction, userID string, minTime, maxTime time.Time) Entries {

	entries := Entries{}

	// index Entry Updated = UserID|Updated|ItemID : ItemID
	min := []byte(keyEncode(userID, keyEncodeTime(minTime)))
	max := []byte(keyEncode(userID, keyEncodeTime(maxTime)))
	c := tx.Bucket(bucketIndex, entityEntry, indexEntryUpdated).Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		entryID := string(v)
		if entry := z.Get(tx, entryID); entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries

}

func (z *entryStore) Save(tx Transaction, entry *Entry) error {
	return save(tx, entityEntry, entry)
}

func (z *entryStore) SaveAll(tx Transaction, entries Entries) error {
	for _, entry := range entries {
		if err := z.Save(tx, entry); err != nil {
			return err
		}
	}
	return nil
}

func (z *entryStore) Starred(tx Transaction, userID string, minTime, maxTime time.Time) Entries {

	entries := Entries{}

	// index Entry Star = UserID|Star|Updated|ItemID : ItemID
	min := []byte(keyEncode(userID, keyEncodeBool(true), keyEncodeTime(minTime)))
	max := []byte(keyEncode(userID, keyEncodeBool(true), keyEncodeTime(maxTime)))
	c := tx.Bucket(bucketIndex, entityEntry, indexEntryStar).Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		entryID := string(v)
		if entry := z.Get(tx, entryID); entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries

}

func (z *entryStore) Unread(tx Transaction, userID string, minTime, maxTime time.Time) Entries {

	entries := Entries{}

	// index Entry Read = UserID|Read|Updated|ItemID : ItemID
	min := []byte(keyEncode(userID, keyEncodeBool(false), keyEncodeTime(minTime)))
	max := []byte(keyEncode(userID, keyEncodeBool(false), keyEncodeTime(maxTime)))
	c := tx.Bucket(bucketIndex, entityEntry, indexEntryRead).Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		entryID := string(v)
		if entry := z.Get(tx, entryID); entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries

}
