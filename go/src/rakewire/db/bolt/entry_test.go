package bolt

import (
	"fmt"
	m "rakewire/model"
	"testing"
)

func TestGetFeedEntriesFromIDs(t *testing.T) {

	db := openDatabase(t)
	defer closeDatabase(t, db)
	if db == nil {
		t.Fatal("cannot open database")
	}

	feed := m.NewFeed("http://localhost/")

	for i := 0; i < 10; i++ {
		entry := feed.AddEntry(fmt.Sprintf("http://localhost/post/%d", i))
		entry.Title = fmt.Sprintf("Post %d", i)
		entry.GenerateNewID()
	}

	if err := db.SaveFeed(feed); err != nil {
		t.Fatal("error saving feed")
	}

	var entryIDs []string
	for _, entry := range feed.Entries {
		entryIDs = append(entryIDs, entry.EntryID)
	}

	if entries, err := db.GetFeedEntriesFromIDs(feed.ID, entryIDs); err != nil {
		t.Fatal("error retrieving entries")
	} else {

		if entries == nil {
			t.Fatal("entries is nil")
		}

		if len(entries) != 10 {
			t.Errorf("bad entry count, expected %d, actual: %d", 10, len(entries))
		}

		for i := 0; i < 10; i++ {
			entry := entries[entryIDs[i]]
			t.Logf("Entry %s: %s", entry.EntryID, entry.Title)
			if entry == nil {
				t.Fatalf("entry is nil: %d", i)
			}
			if entry.EntryID != feed.Entries[i].EntryID {
				t.Fatalf("not equal, expected %s, actual %s", feed.Entries[i].EntryID, entry.EntryID)
			}
			if entry.Title != feed.Entries[i].Title {
				t.Fatalf("not equal, expected %s, actual %s", feed.Entries[i].Title, entry.Title)
			}
		}

	}

}
