package fever

import (
	"rakewire/model"
)

func (z *API) getItemsAll(userID string, tx model.Transaction) ([]*Item, error) {

	entries := model.E.Range(tx, userID)
	itemsByID := model.I.GetByEntries(tx, entries).ByID()

	feverItems := []*Item{}
	for _, entry := range entries {
		feverItem := toItem(entry, itemsByID[entry.ItemID])
		feverItems = append(feverItems, feverItem)
	}
	return feverItems, nil

}

func (z *API) getItemsNext(userID, minID string, tx model.Transaction) ([]*Item, error) {

	min := parseID(minID)
	max := min + 100

	entries := model.E.Range(tx, userID, encodeID(minID), formatID(max)).Limit(50)
	itemsByID := model.I.GetByEntries(tx, entries).ByID()

	feverItems := []*Item{}
	for _, entry := range entries {
		feverItem := toItem(entry, itemsByID[entry.ItemID])
		feverItems = append(feverItems, feverItem)
	}
	return feverItems, nil

}

func (z *API) getItemsPrev(userID, maxID string, tx model.Transaction) ([]*Item, error) {

	max := parseID(maxID)
	min := max - 100

	entries := model.E.Range(tx, userID, formatID(min), encodeID(maxID)).Reverse().Limit(50)
	itemsByID := model.I.GetByEntries(tx, entries).ByID()

	feverItems := []*Item{}
	for _, entry := range entries {
		feverItem := toItem(entry, itemsByID[entry.ItemID])
		feverItems = append(feverItems, feverItem)
	}
	return feverItems, nil

}

func (z *API) getItemsByIds(userID string, ids []string, tx model.Transaction) ([]*Item, error) {

	feverItems := []*Item{}
	for _, id := range ids {
		itemID := encodeID(id)
		if entry := model.E.Get(tx, userID, itemID); entry != nil {
			item := model.I.Get(tx, itemID)
			feverItem := toItem(entry, item)
			feverItems = append(feverItems, feverItem)
		}
	}

	return feverItems, nil

}

func toItem(entry *model.Entry, item *model.Item) *Item {
	return &Item{
		ID:             parseID(entry.ItemID),
		SubscriptionID: parseID(entry.FeedID),
		Title:          item.Title,
		Author:         item.Author,
		HTML:           item.Content,
		URL:            item.URL,
		IsSaved:        boolToUint8(entry.Star),
		IsRead:         boolToUint8(entry.Read),
		Created:        item.Created.Unix(),
	}
}
