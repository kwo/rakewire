package modelng

import (
	"bytes"
	"time"
)

// E groups all entry database methods
var E = &entryStore{}

type entryStore struct{}

type entryQuery struct {
	tx     Transaction
	userID string
	feedID string
	min    time.Time
	max    time.Time
}

func (z *entryStore) AddItems(tx Transaction, allItems Items) error {

	mappedItems := allItems.GroupByFeedID()

	for feedID, items := range mappedItems {
		subscriptions := S.GetForFeed(tx, feedID)
		for _, subscription := range subscriptions {
			for _, item := range items {
				entry := z.New(subscription.UserID, item.ID, subscription.FeedID)
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

func (z *entryStore) New(userID, itemID, feedID string) *Entry {
	return &Entry{
		UserID: userID,
		ItemID: itemID,
		FeedID: feedID,
	}
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

func (z *entryStore) Query(tx Transaction, userID string) *entryQuery {
	return &entryQuery{
		tx:     tx,
		userID: userID,
		min:    time.Time{},
		max:    time.Now(),
	}
}

func (z *entryQuery) Feed(feedID string) *entryQuery {
	z.feedID = feedID
	return z
}

func (z *entryQuery) Max(max time.Time) *entryQuery {
	z.max = max
	return z
}

func (z *entryQuery) Min(min time.Time) *entryQuery {
	z.min = min
	return z
}

func (z *entryQuery) Count() int {

	result := 0

	var c Cursor
	var min, max []byte

	if z.feedID != empty {
		// index Entry FeedUpdated = UserID|FeedID|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryFeedUpdated).Cursor()
		min = []byte(keyEncode(z.userID, z.feedID, keyEncodeTime(z.min)))
		max = []byte(keyEncode(z.userID, z.feedID, keyEncodeTime(z.max)))
	} else {
		// index Entry Updated = UserID|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryUpdated).Cursor()
		min = []byte(keyEncode(z.userID, keyEncodeTime(z.min)))
		max = []byte(keyEncode(z.userID, keyEncodeTime(z.max)))
	}

	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, _ = c.Next() {
		result++
	}

	return result

}

func (z *entryQuery) Get() Entries {

	entries := Entries{}

	var c Cursor
	var min, max []byte

	if z.feedID != empty {
		// index Entry FeedUpdated = UserID|FeedID|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryFeedUpdated).Cursor()
		min = []byte(keyEncode(z.userID, z.feedID, keyEncodeTime(z.min)))
		max = []byte(keyEncode(z.userID, z.feedID, keyEncodeTime(z.max)))
	} else {
		// index Entry Updated = UserID|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryUpdated).Cursor()
		min = []byte(keyEncode(z.userID, keyEncodeTime(z.min)))
		max = []byte(keyEncode(z.userID, keyEncodeTime(z.max)))
	}

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		entryID := string(v)
		if entry := E.Get(z.tx, entryID); entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries

}

func (z *entryQuery) Starred() Entries {

	entries := Entries{}

	var c Cursor
	var min, max []byte

	if z.feedID != empty {
		// index Entry FeedStarUpdated = UserID|FeedID|Star|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryFeedStarUpdated).Cursor()
		min = []byte(keyEncode(z.userID, z.feedID, keyEncodeBool(true), keyEncodeTime(z.min)))
		max = []byte(keyEncode(z.userID, z.feedID, keyEncodeBool(true), keyEncodeTime(z.max.Add(-1*time.Second))))
	} else {
		// index Entry StarUpdated = UserID|Star|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryStarUpdated).Cursor()
		min = []byte(keyEncode(z.userID, keyEncodeBool(true), keyEncodeTime(z.min)))
		max = []byte(keyEncode(z.userID, keyEncodeBool(true), keyEncodeTime(z.max.Add(-1*time.Second))))
	}

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		entryID := string(v)
		if entry := E.Get(z.tx, entryID); entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries

}

func (z *entryQuery) Unread() Entries {

	entries := Entries{}

	var c Cursor
	var min, max []byte

	if z.feedID != empty {
		// index Entry FeedReadUpdated = UserID|FeedID|Read|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryFeedReadUpdated).Cursor()
		min = []byte(keyEncode(z.userID, z.feedID, keyEncodeBool(false), keyEncodeTime(z.min)))
		max = []byte(keyEncode(z.userID, z.feedID, keyEncodeBool(false), keyEncodeTime(z.max.Add(-1*time.Second))))
	} else {
		// index Entry ReadUpdated = UserID|Read|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryReadUpdated).Cursor()
		min = []byte(keyEncode(z.userID, keyEncodeBool(false), keyEncodeTime(z.min)))
		max = []byte(keyEncode(z.userID, keyEncodeBool(false), keyEncodeTime(z.max.Add(-1*time.Second))))
	}

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		entryID := string(v)
		if entry := E.Get(z.tx, entryID); entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries

}
