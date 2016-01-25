package fever

import (
	"rakewire/model"
)

func (z *API) getItemsAll(userID uint64, tx model.Transaction) ([]*Item, error) {

	entries, err := model.UserEntryGetNext(userID, 0, 0, tx)
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

	entries, err := model.UserEntryGetNext(userID, minID, 50, tx)
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

	entries, err := model.UserEntryGetPrev(userID, maxID, 50, tx)
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

	entries, err := model.UserEntriesByUser(userID, ids, tx)
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

func toItem(entry *model.UserEntry) *Item {
	return &Item{
		ID:         entry.ID,
		UserFeedID: entry.UserFeedID,
		Title:      entry.Entry.Title,
		Author:     entry.Entry.Author,
		HTML:       entry.Entry.Content,
		URL:        entry.Entry.URL,
		IsSaved:    boolToUint8(entry.IsStar),
		IsRead:     boolToUint8(entry.IsRead),
		Created:    entry.Entry.Created.Unix(),
	}
}
