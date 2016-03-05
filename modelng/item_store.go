package modelng

import (
	"bytes"
)

// I groups all item database methods
var I = &itemStore{}

type itemStore struct{}

func (z *itemStore) Delete(tx Transaction, id string) error {
	return delete(tx, entityItem, id)
}

// Get returns the item with the given compoundID or the given userID and itemID
func (z *itemStore) Get(tx Transaction, id ...string) *Item {
	compoundID := ""
	switch len(id) {
	case 1:
		compoundID = id[0]
	case 2:
		compoundID = keyEncode(id...)
	default:
		return nil
	}
	bData := tx.Bucket(bucketData, entityItem)
	if data := bData.Get([]byte(compoundID)); data != nil {
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

func (z *itemStore) New(feedID, guid string) *Item {
	return &Item{
		FeedID: feedID,
		GUID:   guid,
	}
}

func (z *itemStore) Save(tx Transaction, item *Item) error {
	return save(tx, entityItem, item)
}
