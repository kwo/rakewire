package model

import (
	"testing"
	"time"
)

func TestGUID(t *testing.T) {

	e := NewEntry(123, "guid")
	if e.ID != 0 {
		t.Error("entry.ID cannot be set by factory method")
	}
	if e.FeedID != 123 {
		t.Error("entry.feedID not set correctly by factory method")
	}
	if e.GUID != "guid" {
		t.Error("entry.GUID not set correctly by factory method")
	}

}

func TestEntryHash(t *testing.T) {

	e := &Entry{}
	lastHash := e.Hash()

	e.ID = 1
	if h := e.Hash(); h != lastHash {
		t.Fatal("ID should not be part of entry hash")
	}

	e.GUID = "GUID"
	if h := e.Hash(); h != lastHash {
		t.Fatal("guID should not be part of entry hash")
	}

	e.FeedID = 123
	if h := e.Hash(); h != lastHash {
		t.Fatal("guID should not be part of entry hash")
	}

	e.Created = time.Now()
	if h := e.Hash(); h != lastHash {
		t.Fatal("Created should not be part of entry hash")
	}

	e.Updated = time.Now()
	if h := e.Hash(); h != lastHash {
		t.Fatal("Updated should not be part of entry hash")
	}

	e.URL = "url"
	if h := e.Hash(); h == lastHash {
		t.Fatal("URL should be part of entry hash")
	} else {
		lastHash = h
	}

	e.Author = "author"
	if h := e.Hash(); h == lastHash {
		t.Fatal("Author should be part of entry hash")
	} else {
		lastHash = h
	}

	e.Title = "title"
	if h := e.Hash(); h == lastHash {
		t.Fatal("Title should be part of entry hash")
	} else {
		lastHash = h
	}

	e.Content = "content"
	if h := e.Hash(); h == lastHash {
		t.Fatal("Content should be part of entry hash")
	} else {
		lastHash = h
	}

}

func TestEntryHashEmpty(t *testing.T) {

	e1 := NewEntry(123, "guID")
	e2 := NewEntry(123, "guID")

	h1 := e1.Hash()
	h2 := e2.Hash()

	if len(h1) != 64 {
		t.Errorf("bad hash length: %d", len(h1))
	}

	if h1 != h2 {
		t.Errorf("hashes do not match, expected %s, actual %s", h1, h2)
	}

}
