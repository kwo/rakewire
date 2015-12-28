package fever

import (
	m "rakewire/model"
)

func (z *API) getItemsAll(userID uint64) ([]*Item, error) {

	entries, err := z.db.UserEntryGetNext(userID, 0, 0)
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

func (z *API) getItemsNext(userID uint64, minID uint64) ([]*Item, error) {

	entries, err := z.db.UserEntryGetNext(userID, minID, 50)
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

func (z *API) getItemsPrev(userID uint64, maxID uint64) ([]*Item, error) {

	entries, err := z.db.UserEntryGetPrev(userID, maxID, 50)
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

func (z *API) getItemsByIds(userID uint64, ids []uint64) ([]*Item, error) {

	entries, err := z.db.UserEntryGet(userID, ids)
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

func toItem(entry *m.UserEntry) *Item {
	return &Item{
		ID:         entry.ID,
		UserFeedID: entry.UserFeedID,
		Title:      entry.Entry.Title,
		Author:     entry.Entry.Author,
		HTML:       entry.Entry.Content,
		URL:        entry.Entry.URL,
		IsSaved:    boolToUint8(entry.IsStarred),
		IsRead:     boolToUint8(entry.IsRead),
		Created:    entry.Entry.Created.Unix(),
	}
}
