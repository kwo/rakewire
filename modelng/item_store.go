package modelng

import (
	"bytes"
)

// I groups all item database methods
var I = &itemStore{}

type itemStore struct{}

func (z *itemStore) Delete(id string, tx Transaction) error {
	return delete(entityItem, id, tx)
}

func (z *itemStore) Get(id string, tx Transaction) *Item {
	bData := tx.Bucket(bucketData, entityItem)
	if data := bData.Get([]byte(id)); data != nil {
		item := &Item{}
		if err := item.decode(data); err == nil {
			return item
		}
	}
	return nil
}

func (z *itemStore) GetByGUID(feedID, guid string, tx Transaction) *Item {
	// index Item GUID = FeedID|GUID : ItemID
	b := tx.Bucket(bucketIndex, entityItem, indexItemGUID)
	if value := b.Get([]byte(keyEncode(feedID, guid))); value != nil {
		itemID := string(value)
		return I.Get(itemID, tx)
	}
	return nil
}

func (z *itemStore) GetForFeed(feedID string, tx Transaction) Items {
	// index Item GUID = FeedID|GUID : ItemID
	items := Items{}
	min, max := keyMinMax(feedID)
	b := tx.Bucket(bucketIndex, entityItem, indexItemGUID)
	c := b.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		itemID := string(v)
		if item := I.Get(itemID, tx); item != nil {
			items = append(items, item)
		}
	}
	return items
}

func (z *itemStore) New(feedID, guid string) *Item {
	return &Item{
		FeedID: feedID,
		GUID:   guid,
	}
}

func (z *itemStore) Save(item *Item, tx Transaction) error {
	return save(entityItem, item, tx)
}
