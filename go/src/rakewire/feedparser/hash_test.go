package feedparser

import (
	"testing"
	"time"
)

func TestEntryHash(t *testing.T) {

	e := &Entry{
		Links: make(map[string]string),
	}
	lastHash := e.hash()

	for i := 0; i < 9; i++ {

		switch i {

		case 0:
			e.ID = "id"
			if h := e.hash(); h != lastHash {
				t.Fatal("ID should not be part of entry hash")
			}

		case 1:
			e.Authors = []string{"author1"}
			if h := e.hash(); h != lastHash {
				t.Fatal("authors should not be part of entry hash")
			}

		case 2:
			e.Categories = []string{"cat1"}
			if h := e.hash(); h != lastHash {
				t.Fatal("categories should not be part of entry hash")
			}

		case 3:
			e.Contributors = []string{"contributors"}
			if h := e.hash(); h != lastHash {
				t.Fatal("contributors should not be part of entry hash")
			}

		case 4:
			e.Created = time.Now()
			if h := e.hash(); h != lastHash {
				t.Fatal("created should not be part of entry hash")
			}

		case 5:
			e.Updated = time.Now()
			if h := e.hash(); h != lastHash {
				t.Fatal("updated should not be part of entry hash")
			}

		case 6:
			e.Links[linkAlternate] = "alternate"
			if h := e.hash(); h == lastHash {
				t.Fatal("link[alternate] should be part of entry hash")
			} else {
				lastHash = h
			}

		case 7:
			e.Summary = "summary"
			if h := e.hash(); h == lastHash {
				t.Fatal("summary should be part of entry hash")
			} else {
				lastHash = h
			}

		case 8:
			e.Title = "title"
			if h := e.hash(); h == lastHash {
				t.Fatal("Title should be part of entry hash")
			} else {
				lastHash = h
			}

		case 9:
			e.Content = "content"
			if h := e.hash(); h == lastHash {
				t.Fatal("Content should be part of entry hash")
			} else {
				lastHash = h
			}

		}

	}

}

func TestEntryHashEmpty(t *testing.T) {

	e := &Entry{}
	h := e.hash()

	if len(h) != 32 {
		t.Errorf("bad hash length: %d", len(h))
	}

}
