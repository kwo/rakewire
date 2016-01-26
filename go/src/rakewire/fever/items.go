package fever

import (
	"rakewire/model"
)

func (z *API) getItemsAll(userID uint64, tx model.Transaction) ([]*Item, error) {

	entries, err := model.EntriesGetNext(userID, 0, 0, tx)
	if err != nil {
		return nil, err
	}

	items := []*Item{}
	for _, entry := range entries {
		item := toItem(entry)
		items = append(items, item)
	}
	return items, nil

}

func (z *API) getItemsNext(userID uint64, minID uint64, tx model.Transaction) ([]*Item, error) {

	entries, err := model.EntriesGetNext(userID, minID, 50, tx)
	if err != nil {
		return nil, err
	}

	items := []*Item{}
	for _, entry := range entries {
		item := toItem(entry)
		items = append(items, item)
	}
	return items, nil

}

func (z *API) getItemsPrev(userID uint64, maxID uint64, tx model.Transaction) ([]*Item, error) {

	entries, err := model.EntriesGetPrev(userID, maxID, 50, tx)
	if err != nil {
		return nil, err
	}

	items := []*Item{}
	for _, entry := range entries {
		item := toItem(entry)
		items = append(items, item)
	}
	return items, nil

}

func (z *API) getItemsByIds(userID uint64, ids []uint64, tx model.Transaction) ([]*Item, error) {

	entries, err := model.EntriesByUser(userID, ids, tx)
	if err != nil {
		return nil, err
	}

	items := []*Item{}
	for _, entry := range entries {
		item := toItem(entry)
		items = append(items, item)
	}
	return items, nil

}

func toItem(entry *model.Entry) *Item {
	return &Item{
		ID:             entry.ID,
		SubscriptionID: entry.SubscriptionID,
		Title:          entry.Item.Title,
		Author:         entry.Item.Author,
		HTML:           entry.Item.Content,
		URL:            entry.Item.URL,
		IsSaved:        boolToUint8(entry.IsStar),
		IsRead:         boolToUint8(entry.IsRead),
		Created:        entry.Item.Created.Unix(),
	}
}
