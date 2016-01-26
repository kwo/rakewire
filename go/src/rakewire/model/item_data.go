package model

import (
	"bytes"
	"strconv"
)

// ItemsByGUIDs retrieves items for specific GUIDs
func ItemsByGUIDs(feedID uint64, guIDs []string, tx Transaction) (map[string]*Item, error) {

	items := make(map[string]*Item)

	e := &Item{}
	e.FeedID = feedID

	for _, guID := range guIDs {

		e.GUID = guID
		indexKeys := e.indexKeys()[itemIndexGUID]

		if data, ok := kvGetFromIndex(itemEntity, itemIndexGUID, indexKeys, tx); ok {
			item := &Item{}
			if err := item.deserialize(data); err != nil {
				return nil, err
			}
			items[item.GUID] = item
		}

	} // loop

	return items, nil

}

// ItemsByFeed retrieves items for the given feed
func ItemsByFeed(feedID uint64, tx Transaction) ([]*Item, error) {

	items := []*Item{}

	// define index keys
	e := &Item{}
	e.FeedID = feedID
	minKeys := e.indexKeys()[itemIndexGUID]
	e.FeedID = feedID + 1
	nxtKeys := e.indexKeys()[itemIndexGUID]

	bIndex := tx.Bucket(bucketIndex).Bucket(itemEntity).Bucket(itemIndexGUID)
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

		if data, ok := kvGet(id, bItem); ok {
			item := &Item{}
			if err := item.deserialize(data); err != nil {
				return nil, err
			}
			items = append(items, item)
		}

	} // loop

	return items, nil

}

// Delete removes an item from the database.
func (item *Item) Delete(tx Transaction) error {
	return kvSave(itemEntity, item, tx)
}
