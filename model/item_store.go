package model

import (
	"bytes"
)

// I groups all item database methods
var I = &itemStore{}

type itemStore struct{}

func (z *itemStore) Delete(tx Transaction, id string) error {
	return deleteObject(tx, entityItem, id)
}

// Get returns the item with the given itemID
func (z *itemStore) Get(tx Transaction, id string) *Item {
	bData := tx.Bucket(bucketData, entityItem)
	if data := bData.Get([]byte(id)); data != nil {
		item := &Item{}
		if err := item.decode(data); err == nil {
			return item
		}
	}
	return nil
}

func (z *itemStore) GetByGUID(tx Transaction, feedID, guid string) *Item {
	// index Item GUID = FeedID|GUID : ItemID
	b := tx.Bucket(bucketIndex, entityItem, indexItemGUID)
	if value := b.Get([]byte(keyEncode(feedID, guid))); value != nil {
		itemID := string(value)
		return z.Get(tx, itemID)
	}
	return nil
}

func (z *itemStore) GetForFeed(tx Transaction, feedID string) Items {
	// index Item GUID = FeedID|GUID : ItemID
	items := Items{}
	min, max := keyMinMax(feedID)
	b := tx.Bucket(bucketIndex, entityItem, indexItemGUID)
	c := b.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		itemID := string(v)
		if item := z.Get(tx, itemID); item != nil {
			items = append(items, item)
		}
	}
	return items
}

func (z *itemStore) GetByEntries(tx Transaction, entries Entries) Items {
	result := Items{}
	for _, entry := range entries {
		if item := z.Get(tx, entry.ItemID); item != nil {
			result = append(result, item)
		}
	}
	return result
}

func (z *itemStore) New(feedID, guid string) *Item {
	return &Item{
		FeedID: feedID,
		GUID:   guid,
	}
}

func (z *itemStore) Save(tx Transaction, item *Item) error {
	return saveObject(tx, entityItem, item)
}

func (z *itemStore) SaveAll(tx Transaction, items Items) error {
	for _, item := range items {
		if err := z.Save(tx, item); err != nil {
			return err
		}
	}
	return nil
}
