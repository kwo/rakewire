package model

import (
	"testing"
	"time"
)

func TestEntryID(t *testing.T) {

	e := NewEntry("feedID", "entryID")
	if e.ID != "" {
		t.Error("entry.ID cannot be set by factory method")
	}
	if e.FeedID != "feedID" {
		t.Error("entry.feedID not set correctly by factory method")
	}
	if e.EntryID != "entryID" {
		t.Error("entry.entryID not set correctly by factory method")
	}

	e.GenerateNewID()
	if len(e.ID) != 36 {
		t.Errorf("entry.ID not generated property, expected length: %d, actual: %d", 36, len(e.ID))
	}

}

func TestEntryHash(t *testing.T) {

	e := &Entry{}
	lastHash := e.Hash()

	e.GenerateNewID()
	if h := e.Hash(); h != lastHash {
		t.Fatal("ID should not be part of entry hash")
	}

	e.EntryID = "entryID"
	if h := e.Hash(); h != lastHash {
		t.Fatal("entryID should not be part of entry hash")
	}

	e.FeedID = "feedID"
	if h := e.Hash(); h != lastHash {
		t.Fatal("entryID should not be part of entry hash")
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

	e1 := NewEntry("feedID", "entryID")
	e2 := NewEntry("feedID", "entryID")

	h1 := e1.Hash()
	h2 := e2.Hash()

	if len(h1) != 64 {
		t.Errorf("bad hash length: %d", len(h1))
	}

	if h1 != h2 {
		t.Errorf("hashes do not match, expected %s, actual %s", h1, h2)
	}

}
