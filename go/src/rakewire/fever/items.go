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

func (z *API) getItemsNext(userID, minID0 string, tx model.Transaction) ([]*Item, error) {

	minID, err := encodeID(minID0)
	if err != nil {
		return nil, err
	}

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

func (z *API) getItemsPrev(userID, maxID0 string, tx model.Transaction) ([]*Item, error) {

	maxID, err := encodeID(maxID0)
	if err != nil {
		return nil, err
	}

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

func (z *API) getItemsByIds(userID string, ids0 []string, tx model.Transaction) ([]*Item, error) {

	ids := []string{}
	for _, id0 := range ids0 {
		id, err := encodeID(id0)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
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
		ID:             decodeID(entry.ID),
		SubscriptionID: decodeID(entry.SubscriptionID),
		Title:          entry.Item.Title,
		Author:         entry.Item.Author,
		HTML:           entry.Item.Content,
		URL:            entry.Item.URL,
		IsSaved:        boolToUint8(entry.IsStar),
		IsRead:         boolToUint8(entry.IsRead),
		Created:        entry.Item.Created.Unix(),
	}
}
