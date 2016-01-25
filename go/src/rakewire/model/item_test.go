package model

import (
	"testing"
	"time"
)

func TestGUID(t *testing.T) {

	e := NewItem(123, "guid")
	if e.ID != 0 {
		t.Error("item.ID cannot be set by factory method")
	}
	if e.FeedID != 123 {
		t.Error("item.feedID not set correctly by factory method")
	}
	if e.GUID != "guid" {
		t.Error("item.GUID not set correctly by factory method")
	}

}

func TestItemHash(t *testing.T) {

	e := &Item{}
	lastHash := e.Hash()

	e.ID = 1
	if h := e.Hash(); h != lastHash {
		t.Fatal("ID should not be part of item hash")
	}

	e.GUID = "GUID"
	if h := e.Hash(); h != lastHash {
		t.Fatal("guID should not be part of item hash")
	}

	e.FeedID = 123
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

	e1 := NewItem(123, "guID")
	e2 := NewItem(123, "guID")

	h1 := e1.Hash()
	h2 := e2.Hash()

	if len(h1) != 64 {
		t.Errorf("bad hash length: %d", len(h1))
	}

	if h1 != h2 {
		t.Errorf("hashes do not match, expected %s, actual %s", h1, h2)
	}

}
