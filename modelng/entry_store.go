package modelng

// E groups all entry database methods
var E = &entryStore{}

type entryStore struct{}

// addItems
// getItems(userID)
// getItemTotal(userID)
// getUnread(userID) - date, count
// getStarred(userID) - date, count
// markRead - by feed, group, all
// markStarred - individually

func (z *entryStore) Delete(tx Transaction, id string) error {
	return delete(tx, entityEntry, id)
}

func (z *entryStore) Get(id string, tx Transaction) *Entry {
	bData := tx.Bucket(bucketData, entityEntry)
	if data := bData.Get([]byte(id)); data != nil {
		entry := &Entry{}
		if err := entry.decode(data); err == nil {
			return entry
		}
	}
	return nil
}

func (z *entryStore) GetByIDs(userID, itemID string, tx Transaction) *Entry {
	return z.Get(keyEncode(userID, itemID), tx)
}

func (z *entryStore) New(userID, itemID string) *Entry {
	return &Entry{
		UserID: userID,
		ItemID: itemID,
	}
}

func (z *entryStore) Save(tx Transaction, entry *Entry) error {
	return save(tx, entityEntry, entry)
}
