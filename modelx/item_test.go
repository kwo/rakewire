package model

import (
	"testing"
	"time"
)

func TestGUID(t *testing.T) {

	e := NewItem(kvKeyUintEncode(123), "guid")
	if e.ID != empty {
		t.Error("item.ID cannot be set by factory method")
	}
	if e.FeedID != kvKeyUintEncode(123) {
		t.Error("item.feedID not set correctly by factory method")
	}
	if e.GUID != "guid" {
		t.Error("item.GUID not set correctly by factory method")
	}

}

func TestItemHash(t *testing.T) {

	e := &Item{}
	lastHash := e.Hash()

	e.ID = kvKeyUintEncode(1)
	if h := e.Hash(); h != lastHash {
		t.Fatal("ID should not be part of item hash")
	}

	e.GUID = "GUID"
	if h := e.Hash(); h != lastHash {
		t.Fatal("guID should not be part of item hash")
	}

	e.FeedID = kvKeyUintEncode(123)
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

	e1 := NewItem(kvKeyUintEncode(123), "guID")
	e2 := NewItem(kvKeyUintEncode(123), "guID")

	h1 := e1.Hash()
	h2 := e2.Hash()

	if len(h1) != 64 {
		t.Errorf("bad hash length: %d", len(h1))
	}

	if h1 != h2 {
		t.Errorf("hashes do not match, expected %s, actual %s", h1, h2)
	}

}

func TestItemsByGUID(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t, true)
	defer closeTestDatabase(t, database)

	err := database.Select(func(tx Transaction) error {

		bIndex := tx.Bucket(bucketIndex, itemEntity, itemIndexGUID)
		c := bIndex.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			t.Logf("k: %v, v: %v", string(k), string(v))
		}

		item, err := ItemByGUID("0000000001", "Feed0Item0", tx)
		if err != nil {
			return err
		}

		if item == nil {
			t.Error("Expected item but is nil")
		}

		return nil

	})

	if err != nil {
		t.Errorf("Error when selecting from database: %s", err.Error())
	}

}
