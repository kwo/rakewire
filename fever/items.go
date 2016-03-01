package fever

import (
	"rakewire/model"
)

func (z *API) getItemsAll(userID string, tx model.Transaction) ([]*Item, error) {

	entries, err := model.EntriesGetAll(userID, tx)
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

func (z *API) getItemsNext(userID, minID string, tx model.Transaction) ([]*Item, error) {

	entries, err := model.EntriesGetNext(userID, encodeID(minID), 50, tx)
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

func (z *API) getItemsPrev(userID, maxID string, tx model.Transaction) ([]*Item, error) {

	entries, err := model.EntriesGetPrev(userID, encodeID(maxID), 50, tx)
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

func (z *API) getItemsByIds(userID string, ids0 []string, tx model.Transaction) ([]*Item, error) {

	ids := []string{}
	for _, id0 := range ids0 {
		ids = append(ids, encodeID(id0))
	}

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
		ID:             parseID(entry.ID),
		SubscriptionID: parseID(entry.SubscriptionID),
		Title:          entry.Item.Title,
		Author:         entry.Item.Author,
		HTML:           entry.Item.Content,
		URL:            entry.Item.URL,
		IsSaved:        boolToUint8(entry.IsStar),
		IsRead:         boolToUint8(entry.IsRead),
		Created:        entry.Item.Created.Unix(),
	}
}
