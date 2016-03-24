package model

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
	return deleteObject(tx, entityEntry, id)
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

// Range returns Entries for the given user by internal ID.
// The optional string parameters are for minID (inclusive) and maxID (exclusive) respectively.
func (z *entryStore) Range(tx Transaction, userID string, minmax ...string) Entries {

	entries := Entries{}

	// bucket Entry = UserID|ItemID : value
	min, max := keyMinMax(userID)
	switch len(minmax) {
	case 1:
		min = []byte(keyEncode(userID, minmax[0]))
	case 2:
		min = []byte(keyEncode(userID, minmax[0]))
		max = []byte(keyEncode(userID, minmax[1]))
	}

	//fmt.Printf("Range: min: %s, max: %s\n", string(min), string(max))
	c := tx.Bucket(bucketData, entityEntry).Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) < 0; k, v = c.Next() {
		entry := &Entry{}
		if err := entry.decode(v); err == nil {
			entries = append(entries, entry)
		}
	}

	return entries

}

func (z *entryStore) Save(tx Transaction, entry *Entry) error {
	return saveObject(tx, entityEntry, entry)
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
		max:    time.Now().Add(1 * time.Second).Truncate(time.Second), // prevents losing entries when saving and then querying
	}
}

func (z *entryQuery) Feed(feedID string) *entryQuery {
	z.feedID = feedID
	return z
}

// Max sets the latest update time for entries in the query, exclusive.
// Time resolution is precise to the second.
func (z *entryQuery) Max(max time.Time) *entryQuery {
	z.max = max
	return z
}

// Min sets the earliest update time for entries in the query.
// Time resolution is precise to the second.
func (z *entryQuery) Min(min time.Time) *entryQuery {
	z.min = min
	return z
}

func (z *entryQuery) Count() uint {

	var result uint

	var c Cursor
	var min, nxt []byte

	if z.feedID != empty {
		// index Entry FeedUpdated = UserID|FeedID|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryFeedUpdated).Cursor()
		min = []byte(keyEncode(z.userID, z.feedID, keyEncodeTime(z.min)))
		nxt = []byte(keyEncode(z.userID, z.feedID, keyEncodeTime(z.max)))
	} else {
		// index Entry Updated = UserID|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryUpdated).Cursor()
		min = []byte(keyEncode(z.userID, keyEncodeTime(z.min)))
		nxt = []byte(keyEncode(z.userID, keyEncodeTime(z.max)))
	}

	//fmt.Printf("min: %s, max: %s\n", string(min), string(max))
	for k, _ := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, _ = c.Next() {
		//fmt.Printf("key: %s\n", string(k))
		result++
	}

	return result

}

func (z *entryQuery) Get() Entries {

	entries := Entries{}

	var c Cursor
	var min, nxt []byte

	if z.feedID != empty {
		// index Entry FeedUpdated = UserID|FeedID|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryFeedUpdated).Cursor()
		min = []byte(keyEncode(z.userID, z.feedID, keyEncodeTime(z.min)))
		nxt = []byte(keyEncode(z.userID, z.feedID, keyEncodeTime(z.max)))
	} else {
		// index Entry Updated = UserID|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryUpdated).Cursor()
		min = []byte(keyEncode(z.userID, keyEncodeTime(z.min)))
		nxt = []byte(keyEncode(z.userID, keyEncodeTime(z.max)))
	}

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
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
	var min, nxt []byte

	if z.feedID != empty {
		// index Entry FeedStarUpdated = UserID|FeedID|Star|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryFeedStarUpdated).Cursor()
		min = []byte(keyEncode(z.userID, z.feedID, keyEncodeBool(true), keyEncodeTime(z.min)))
		nxt = []byte(keyEncode(z.userID, z.feedID, keyEncodeBool(true), keyEncodeTime(z.max)))
	} else {
		// index Entry StarUpdated = UserID|Star|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryStarUpdated).Cursor()
		min = []byte(keyEncode(z.userID, keyEncodeBool(true), keyEncodeTime(z.min)))
		nxt = []byte(keyEncode(z.userID, keyEncodeBool(true), keyEncodeTime(z.max)))
	}

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
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
	var min, nxt []byte

	if z.feedID != empty {
		// index Entry FeedReadUpdated = UserID|FeedID|Read|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryFeedReadUpdated).Cursor()
		min = []byte(keyEncode(z.userID, z.feedID, keyEncodeBool(false), keyEncodeTime(z.min)))
		nxt = []byte(keyEncode(z.userID, z.feedID, keyEncodeBool(false), keyEncodeTime(z.max)))
	} else {
		// index Entry ReadUpdated = UserID|Read|Updated|ItemID : ItemID
		c = z.tx.Bucket(bucketIndex, entityEntry, indexEntryReadUpdated).Cursor()
		min = []byte(keyEncode(z.userID, keyEncodeBool(false), keyEncodeTime(z.min)))
		nxt = []byte(keyEncode(z.userID, keyEncodeBool(false), keyEncodeTime(z.max)))
	}

	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {
		entryID := string(v)
		if entry := E.Get(z.tx, entryID); entry != nil {
			entries = append(entries, entry)
		}
	}

	return entries

}
