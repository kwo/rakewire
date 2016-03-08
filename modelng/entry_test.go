package modelng

import (
	"fmt"
	"testing"
	"time"
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
	feedID := "0000000003"

	// add entry
	err := db.Update(func(tx Transaction) error {

		entry := E.New(userID, itemID, feedID)
		if err := E.Save(tx, entry); err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error adding entry: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {

		entry := E.Get(tx, userID, itemID)
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
		if err := E.Delete(tx, keyEncode(userID, itemID)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error deleting entry: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {
		entry := E.Get(tx, userID, itemID)
		if entry != nil {
			t.Error("Expected nil entry")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting entry: %s", err.Error())
	}

}

func TestEntriesAgain(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	feeds := Feeds{}
	users := Users{}

	err := db.Update(func(tx Transaction) error {

		// add feeds
		for ff := 0; ff < 20; ff++ {
			f := F.New(fmt.Sprintf("Feed%02d", ff+1))
			if err := F.Save(tx, f); err != nil {
				return err
			}
			feeds = append(feeds, f)
		}

		// add users
		for uu := 0; uu < 2; uu++ {
			u := U.New(fmt.Sprintf("User%02d", uu+1))
			u.SetPassword("ancdefg")
			if err := U.Save(tx, u); err != nil {
				return err
			}
			users = append(users, u)
		}

		// add subscriptions
		for _, user := range users {
			for _, feed := range feeds {
				s := S.New(user.ID, feed.ID)
				if err := S.Save(tx, s); err != nil {
					return err
				}
			}
		}

		// add items
		items := Items{}
		now := time.Now().Add(-30 * time.Minute)
		for f, feed := range feeds {
			lastUpdated := now.Add(time.Duration(-f) * time.Hour)
			for ii := 0; ii < 20; ii++ {
				lastUpdated := lastUpdated.Add(time.Duration(ii) * time.Minute)
				item := I.New(feed.ID, fmt.Sprintf("Feed%02dItem%02d", f+1, ii+1))
				item.Updated = lastUpdated
				if err := I.Save(tx, item); err != nil {
					return err
				}
				//t.Logf("Item: %s", item.getID())
				items = append(items, item)
			}
		}

		return E.AddItems(tx, items)

	})
	if err != nil {
		t.Errorf("Error updating database: %s", err.Error())
	}

	err = db.Select(func(tx Transaction) error {

		expectedCount := 400
		for _, user := range users {
			if count := E.Query(tx, user.ID).Count(); count != expectedCount {
				t.Errorf("Bad count: %d, expected %d", count, expectedCount)
			}
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}
