package model

import (
	"bytes"
	"strings"
)

// ItemByGUID retrieve an item object by guid, nil if not found
func ItemByGUID(guid string, tx Transaction) (item *Item, err error) {

	bItem := tx.Bucket(bucketData, itemEntity)
	bIndex := tx.Bucket(bucketIndex, itemEntity, itemIndexGUID)

	if record := bIndex.GetIndex(bItem, strings.ToLower(guid)); record != nil {
		item = &Item{}
		err = item.deserialize(record)
	}

	return

}

// ItemsByGUID retrieves items for specific GUIDs
func ItemsByGUID(feedID string, guIDs []string, tx Transaction) (Items, error) {

	items := Items{}
	for _, guid := range guIDs {
		item, err := ItemByGUID(guid, tx)
		if err != nil {
			return nil, err
		} else if item != nil {
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
