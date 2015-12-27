package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"testing"
)

func TestUserEntry(t *testing.T) {

	t.Parallel()

	db := openDatabase(t)
	defer closeDatabase(t, db)
	if db == nil {
		t.Fatal("cannot open database")
	}

	var users []*m.User
	for i := 0; i < 2; i++ {
		user := m.NewUser(fmt.Sprintf("User%d", i))
		user.SetPassword("abcdefg")
		if err := db.UserSave(user); err != nil {
			t.Fatalf("Error saving user: %s", err.Error())
		}
		users = append(users, user)
	}

	var feeds []*m.Feed
	for i := 0; i < 2; i++ {
		feed := m.NewFeed(fmt.Sprintf("http://localhost%d/", i))
		feed.Title = fmt.Sprintf("Feed%d", i)
		if _, err := db.FeedSave(feed); err != nil {
			t.Fatalf("Error saving feed: %s", err.Error())
		}
		feeds = append(feeds, feed)
	}

	// save userfeeds
	if err := db.UserFeedSave(m.NewUserFeed(users[0].ID, feeds[0].ID)); err != nil {
		t.Fatalf("Error saving userfeed: %s", err.Error())
	}
	if err := db.UserFeedSave(m.NewUserFeed(users[1].ID, feeds[0].ID)); err != nil {
		t.Fatalf("Error saving userfeed: %s", err.Error())
	}
	if err := db.UserFeedSave(m.NewUserFeed(users[1].ID, feeds[1].ID)); err != nil {
		t.Fatalf("Error saving userfeed: %s", err.Error())
	}

	entries := []*m.Entry{}
	for i := 0; i < 5; i++ {
		entry := m.NewEntry(feeds[0].ID, fmt.Sprintf("Entry%d", i))
		entries = append(entries, entry)
	}
	for i := 0; i < 10; i++ {
		entry := m.NewEntry(feeds[1].ID, fmt.Sprintf("Entry%d", i))
		entries = append(entries, entry)
	}
	err := db.db.Update(func(tx *bolt.Tx) error {
		for _, entry := range entries {
			if err := kvSave(m.EntryEntity, entry, tx); err != nil {
				return err
			}
		}
		return db.UserEntryAddNew(entries, tx)
	})
	if err != nil {
		t.Errorf("Error saving user entries: %s", err.Error())
	}

	userentries, err := db.UserEntryGetUnreadForUser(users[0].ID)
	if err != nil {
		t.Fatalf("Error retrieving user entries: %s", err.Error())
	}
	if len(userentries) != 5 {
		t.Errorf("bad user entries count, expected %d, actual %d", 5, len(userentries))
	}

	userentries, err = db.UserEntryGetUnreadForUser(users[1].ID)
	if err != nil {
		t.Fatalf("Error retrieving user entries: %s", err.Error())
	}
	if len(userentries) != 15 {
		t.Fatalf("bad user entries count, expected %d, actual %d", 15, len(userentries))
	}

	userentries[12].Read = true
	userentries[13].Read = true
	userentries[14].Read = true
	readEntries := []*m.UserEntry{
		userentries[12], userentries[13], userentries[14],
	}

	if err := db.UserEntrySave(readEntries); err != nil {
		t.Fatalf("err saving user entries: %s", err.Error())
	}

	userentries, err = db.UserEntryGetUnreadForUser(users[1].ID)
	if err != nil {
		t.Fatalf("Error retrieving user entries: %s", err.Error())
	}
	if len(userentries) != 12 {
		t.Fatalf("bad user entries count, expected %d, actual %d", 12, len(userentries))
	}

}
