package feedparser

import (
	"testing"
	"time"
)

func TestEntryHash(t *testing.T) {

	e := &Entry{
		Links: make(map[string]string),
	}
	lastHash := HashEntry(e)

	e.ID = "id"
	if h := HashEntry(e); h != lastHash {
		t.Fatal("ID should not be part of entry hash")
	}

	e.Authors = []string{"author1"}
	if h := HashEntry(e); h != lastHash {
		t.Fatal("authors should not be part of entry hash")
	}

	e.Categories = []string{"cat1"}
	if h := HashEntry(e); h != lastHash {
		t.Fatal("categories should not be part of entry hash")
	}

	e.Contributors = []string{"contributors"}
	if h := HashEntry(e); h != lastHash {
		t.Fatal("contributors should not be part of entry hash")
	}

	e.Created = time.Now()
	if h := HashEntry(e); h != lastHash {
		t.Fatal("created should not be part of entry hash")
	}

	e.Updated = time.Now()
	if h := HashEntry(e); h != lastHash {
		t.Fatal("updated should not be part of entry hash")
	}

	e.Links[linkAlternate] = "alternate"
	if h := HashEntry(e); h == lastHash {
		t.Fatal("link[alternate] should be part of entry hash")
	} else {
		lastHash = h
	}

	e.Summary = "summary"
	if h := HashEntry(e); h == lastHash {
		t.Fatal("Summary should be part of entry hash")
	} else {
		lastHash = h
	}

	e.Title = "title"
	if h := HashEntry(e); h == lastHash {
		t.Fatal("Title should be part of entry hash")
	} else {
		lastHash = h
	}

	e.Content = "content"
	if h := HashEntry(e); h == lastHash {
		t.Fatal("Content should be part of entry hash")
		// } else {
		// 	lastHash = h
	}

}

func TestEntryHashEmpty(t *testing.T) {

	e := &Entry{}
	h := HashEntry(e)

	if len(h) != 32 {
		t.Errorf("bad hash length: %d", len(h))
	}

}
