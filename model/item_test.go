package model

import (
	"fmt"
	"testing"
	"time"
)

func TestItemSetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityItem); obj == nil {
		t.Error("missing getObject entry")
	}

	if obj := allEntities[entityItem]; obj == nil {
		t.Error("missing allEntities entry")
	}

	c := &Config{}
	if obj := c.Sequences.Item; obj != 0 {
		t.Error("missing sequences entry")
	}

}

func TestItemGetBadID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	err := db.Select(func(tx Transaction) error {
		if item := I.Get(tx, empty); item != nil {
			t.Error("Expected nil item")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestItemGetBadGUID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	err := db.Select(func(tx Transaction) error {
		if item := I.GetByGUID(tx, empty, empty); item != nil {
			t.Error("Expected nil item")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestItems(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	itemID := empty
	feedID := "0000000002"
	guid := "http://localhorst/"

	// add item
	err := db.Update(func(tx Transaction) error {

		item := I.New(feedID, guid)
		if err := I.Save(tx, item); err != nil {
			return err
		}
		itemID = item.ID

		return nil

	})
	if err != nil {
		t.Errorf("Error adding item: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {

		item := I.Get(tx, itemID)
		if item == nil {
			t.Fatal("Nil item, expected valid item")
		}
		if item.ID != itemID {
			t.Errorf("bad item ID %s, expected %s", item.ID, itemID)
		}
		if item.FeedID != feedID {
			t.Errorf("bad FeedID: %s, expected %s", item.FeedID, feedID)
		}
		if item.GUID != guid {
			t.Errorf("bad guid: %s, expected %s", item.GUID, guid)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting item: %s", err.Error())
	}

	// delete item
	err = db.Update(func(tx Transaction) error {
		if err := I.Delete(tx, itemID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error deleting item: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {
		item := I.Get(tx, itemID)
		if item != nil {
			t.Error("Expected nil item")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting item: %s", err.Error())
	}

}

func TestItemByGUID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	// add items
	err := db.Update(func(tx Transaction) error {
		for f := 0; f < 10; f++ {
			feedID := keyEncodeUint(uint64(f + 1))
			for i := 0; i < 10; i++ {
				guid := fmt.Sprintf("Feed%02dItem%02d", f, i)
				item := I.New(feedID, guid)
				item.Content = guid
				if err := I.Save(tx, item); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error adding items: %s", err.Error())
	}

	// test by guid
	err = db.Select(func(tx Transaction) error {

		itemID := "0000000024"
		feedID := "0000000003"
		guid := "Feed02Item03"

		item := I.GetByGUID(tx, feedID, guid)
		if item == nil {
			t.Fatal("Nil item, expected valid item")
		}
		if item.ID != itemID {
			t.Errorf("bad item ID %s, expected %s", item.ID, itemID)
		}
		if item.FeedID != feedID {
			t.Errorf("bad FeedID: %s, expected %s", item.FeedID, feedID)
		}
		if item.GUID != guid {
			t.Errorf("bad guid: %s, expected %s", item.GUID, guid)
		}
		if item.Content != guid {
			t.Errorf("bad content: %s, expected %s", item.Content, guid)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting item: %s", err.Error())
	}

}

func TestItemForFeed(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	// add items
	err := db.Update(func(tx Transaction) error {
		for f := 0; f < 10; f++ {
			feedID := keyEncodeUint(uint64(f + 1))
			for i := 0; i < 10; i++ {
				guid := fmt.Sprintf("Feed%02dItem%02d", f, i)
				item := I.New(feedID, guid)
				item.Content = guid
				if err := I.Save(tx, item); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error adding items: %s", err.Error())
	}

	// test by guid
	err = db.Select(func(tx Transaction) error {

		feedID := "0000000003"

		items := I.GetForFeed(tx, feedID)
		if items == nil {
			t.Fatal("Nil items, expected valid items")
		}
		if len(items) != 10 {
			t.Errorf("bad item count %d, expected %d", len(items), 10)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting item: %s", err.Error())
	}

}

func TestItemHash(t *testing.T) {

	e := &Item{}
	lastHash := e.Hash()

	e.ID = keyEncodeUint(1)
	if h := e.Hash(); h != lastHash {
		t.Fatal("ID should not be part of item hash")
	}

	e.GUID = "GUID"
	if h := e.Hash(); h != lastHash {
		t.Fatal("guID should not be part of item hash")
	}

	e.FeedID = keyEncodeUint(123)
	if h := e.Hash(); h != lastHash {
		t.Fatal("guID should not be part of item hash")
	}

	e.Created = time.Now()
	if h := e.Hash(); h != lastHash {
		t.Fatal("Created should not be part of item hash")
	}

	e.Updated = time.Now()
	if h := e.Hash(); h != lastHash {
		t.Fatal("Updated should not be part of item hash")
	}

	e.URL = "url"
	if h := e.Hash(); h == lastHash {
		t.Fatal("URL should be part of item hash")
	} else {
		lastHash = h
	}

	e.Author = "author"
	if h := e.Hash(); h == lastHash {
		t.Fatal("Author should be part of item hash")
	} else {
		lastHash = h
	}

	e.Title = "title"
	if h := e.Hash(); h == lastHash {
		t.Fatal("Title should be part of item hash")
	} else {
		lastHash = h
	}

	e.Content = "content"
	if h := e.Hash(); h == lastHash {
		t.Fatal("Content should be part of item hash")
	} else {
		lastHash = h
	}

}

func TestItemHashEmpty(t *testing.T) {

	e1 := I.New(keyEncodeUint(123), "guID")
	e2 := I.New(keyEncodeUint(123), "guID")

	h1 := e1.Hash()
	h2 := e2.Hash()

	if len(h1) != 64 {
		t.Errorf("bad hash length: %d", len(h1))
	}

	if h1 != h2 {
		t.Errorf("hashes do not match, expected %s, actual %s", h1, h2)
	}

}
