package modelng

import (
	"testing"
)

func TestEntrySetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityEntry); obj == nil {
		t.Error("missing getObject entry")
	}

	if obj := allEntities[entityEntry]; obj == nil {
		t.Error("missing allEntities entry")
	}

}

func TestEntries(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	userID := "0000000001"
	itemID := "0000000002"

	// add entry
	err := db.Update(func(tx Transaction) error {

		entry := E.New(userID, itemID)
		if err := E.Save(entry, tx); err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error adding entry: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {

		entry := E.GetByIDs(userID, itemID, tx)
		if entry == nil {
			t.Fatal("Nil entry, expected valid entry")
		}
		if entry.UserID != userID {
			t.Errorf("bad userID: %s, expected %s", entry.UserID, userID)
		}
		if entry.ItemID != itemID {
			t.Errorf("bad itemID: %s, expected %s", entry.ItemID, itemID)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting entry: %s", err.Error())
	}

	// delete entry
	err = db.Update(func(tx Transaction) error {
		if err := E.Delete(keyEncode(userID, itemID), tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error deleting entry: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {
		entry := E.GetByIDs(userID, itemID, tx)
		if entry != nil {
			t.Error("Expected nil entry")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting entry: %s", err.Error())
	}

}
