package model

import (
	"bytes"
)

// ItemsByGUIDs retrieves items for specific GUIDs
func ItemsByGUIDs(feedID string, guIDs []string, tx Transaction) (Items, error) {

	items := Items{}

	// item index GUID = FeedID|GUID : ItemID
	//min, max := kvKeyMinMax(feedID)

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
			items = append(items, item)
		}

	} // loop

	return items, nil

}

// ItemsByFeed retrieves items for the given feed
func ItemsByFeed(feedID string, tx Transaction) (Items, error) {

	items := Items{}

	// item index GUID = FeedID|GUID : ItemID
	min, max := kvKeyMinMax(feedID)
	bIndex := tx.Bucket(bucketIndex).Bucket(itemEntity).Bucket(itemIndexGUID)
	bItem := tx.Bucket(bucketData).Bucket(itemEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {

		itemID := string(v)

		if data, ok := kvGet(itemID, bItem); ok {
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
