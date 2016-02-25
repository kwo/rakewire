package model

// ItemByGUID retrieve an item object by guid, nil if not found
func ItemByGUID(feedID, guid string, tx Transaction) (item *Item, err error) {

	bItem := tx.Bucket(bucketData, itemEntity)
	bIndex := tx.Bucket(bucketIndex, itemEntity, itemIndexGUID)

	key := kvKeyEncode(feedID, guid)

	if record := bIndex.GetIndex(bItem, key); record != nil {
		item = &Item{}
		err = item.deserialize(record)
	}

	return

}

// ItemsByGUID retrieves items for specific GUIDs
func ItemsByGUID(feedID string, guIDs []string, tx Transaction) (Items, error) {

	items := Items{}
	for _, guid := range guIDs {
		item, err := ItemByGUID(feedID, guid, tx)
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

	err := bIndex.IterateIndex(bItem, min, max, func(id string, record Record) error {
		item := &Item{}
		if err := item.deserialize(record); err != nil {
			return err
		}
		items = append(items, item)
		return nil
	})

	return items, err

}

// Delete removes an item from the database.
func (item *Item) Delete(tx Transaction) error {
	return kvSave(itemEntity, item, tx)
}
