package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"testing"
	"time"
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
	for i := 0; i < 4; i++ {
		feed := m.NewFeed(fmt.Sprintf("http://localhost%d/", i))
		feed.Title = fmt.Sprintf("Feed%d", i)
		if _, err := db.FeedSave(feed); err != nil {
			t.Fatalf("Error saving feed: %s", err.Error())
		}
		feeds = append(feeds, feed)
	}

	// create user groups
	g1 := m.NewGroup(users[0].ID, "group1")
	if err := db.GroupSave(g1); err != nil {
		t.Fatalf("Error creating group: %s", err.Error())
	}
	g2 := m.NewGroup(users[0].ID, "group2")
	if err := db.GroupSave(g2); err != nil {
		t.Fatalf("Error creating group: %s", err.Error())
	}

	// save userfeeds
	uf := m.NewUserFeed(users[0].ID, feeds[0].ID)
	uf.AddGroup(g1.ID)
	if err := db.UserFeedSave(uf); err != nil {
		t.Fatalf("Error saving userfeed: %s", err.Error())
	}
	uf = m.NewUserFeed(users[1].ID, feeds[0].ID)
	if err := db.UserFeedSave(uf); err != nil {
		t.Fatalf("Error saving userfeed: %s", err.Error())
	}
	uf = m.NewUserFeed(users[1].ID, feeds[1].ID)
	if err := db.UserFeedSave(uf); err != nil {
		t.Fatalf("Error saving userfeed: %s", err.Error())
	}
	uf = m.NewUserFeed(users[0].ID, feeds[2].ID)
	uf.AddGroup(g2.ID)
	if err := db.UserFeedSave(uf); err != nil {
		t.Fatalf("Error saving userfeed: %s", err.Error())
	}
	uf = m.NewUserFeed(users[0].ID, feeds[3].ID)
	uf.AddGroup(g1.ID)
	if err := db.UserFeedSave(uf); err != nil {
		t.Fatalf("Error saving userfeed: %s", err.Error())
	}

	// add entries
	now := time.Now().Truncate(time.Second)
	entries := []*m.Entry{}
	for i := 0; i < 50; i++ {
		entry := m.NewEntry(feeds[0].ID, fmt.Sprintf("Entry%d", i))
		entry.Created = now.Add(time.Duration(-i) * 24 * time.Hour)
		entry.Updated = now.Add(time.Duration(-i) * 24 * time.Hour)
		entries = append(entries, entry)
	}
	for i := 0; i < 101; i++ {
		entry := m.NewEntry(feeds[1].ID, fmt.Sprintf("Entry%d", i))
		entry.Created = now.Add(time.Duration(-i) * 24 * time.Hour)
		entry.Updated = now.Add(time.Duration(-i) * 24 * time.Hour)
		entries = append(entries, entry)
	}
	for i := 0; i < 102; i++ {
		entry := m.NewEntry(feeds[2].ID, fmt.Sprintf("Entry%d", i))
		entry.Created = now.Add(time.Duration(-i) * 24 * time.Hour)
		entry.Updated = now.Add(time.Duration(-i) * 24 * time.Hour)
		entries = append(entries, entry)
	}
	for i := 0; i < 103; i++ {
		entry := m.NewEntry(feeds[3].ID, fmt.Sprintf("Entry%d", i))
		entry.Created = now.Add(time.Duration(-i) * 24 * time.Hour)
		entry.Updated = now.Add(time.Duration(-i) * 24 * time.Hour)
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

	// test counts
	count, err := db.UserEntryGetTotalForUser(users[0].ID)
	if err != nil {
		t.Errorf("Error retrieving user count: %s", err.Error())
	}
	if count != 255 {
		t.Errorf("user total mismatch, expected %d, actual %d", 255, count)
	}

	count, err = db.UserEntryGetTotalForUser(users[1].ID)
	if err != nil {
		t.Errorf("Error retrieving user count: %s", err.Error())
	}
	if count != 151 {
		t.Errorf("user total mismatch, expected %d, actual %d", 151, count)
	}

	// test unread
	userentries, err := db.UserEntryGetUnreadForUser(users[0].ID)
	if err != nil {
		t.Errorf("Error retrieving user entries: %s", err.Error())
	}
	if len(userentries) != 255 {
		t.Errorf("bad user entries count, expected %d, actual %d", 255, len(userentries))
	}

	userentries, err = db.UserEntryGetUnreadForUser(users[1].ID)
	if err != nil {
		t.Errorf("Error retrieving user entries: %s", err.Error())
	}
	if len(userentries) != 151 {
		t.Fatalf("bad user entries count, expected %d, actual %d", 151, len(userentries))
	}

	userentries[12].IsRead = true
	userentries[13].IsRead = true
	userentries[14].IsRead = true
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
	if len(userentries) != 148 {
		t.Errorf("bad user entries count, expected %d, actual %d", 148, len(userentries))
	}

	// test get next
	userentries, err = db.UserEntryGetNext(users[1].ID, 0, 100)
	if err != nil {
		t.Fatalf("Error retrieving user entries next: %s", err.Error())
	}
	if len(userentries) != 100 {
		t.Fatalf("bad user entries count, expected %d, actual %d", 100, len(userentries))
	}
	if userentries[0].ID > userentries[99].ID {
		t.Errorf("expected userentries in ascending order")
	}

	userentries, err = db.UserEntryGetNext(users[1].ID, userentries[99].ID, 100)
	if err != nil {
		t.Fatalf("Error retrieving user entries next: %s", err.Error())
	}
	if len(userentries) != 51 {
		t.Errorf("bad user entries count, expected %d, actual %d", 51, len(userentries))
	}

	// test get prev
	userentries, err = db.UserEntryGetPrev(users[1].ID, 99999, 100)
	if err != nil {
		t.Fatalf("Error retrieving user entries prev: %s", err.Error())
	}
	if len(userentries) != 100 {
		t.Fatalf("bad user entries count, expected %d, actual %d", 100, len(userentries))
	}
	if userentries[0].ID < userentries[99].ID {
		t.Errorf("expected userentries in descending order")
	}

	userentries, err = db.UserEntryGetPrev(users[1].ID, userentries[99].ID, 100)
	if err != nil {
		t.Fatalf("Error retrieving user entries next: %s", err.Error())
	}
	if len(userentries) != 51 {
		t.Errorf("bad user entries count, expected %d, actual %d", 51, len(userentries))
	}

	// test get by ID
	userentries, err = db.UserEntryGetByID(users[1].ID, []uint64{0, 1, 2})
	if err != nil {
		t.Fatalf("Error retrieving user entries by ID: %s", err.Error())
	}
	if len(userentries) != 2 {
		t.Fatalf("bad user entries count, expected %d, actual %d", 2, len(userentries))
	}
	if userentries[0].ID != 1 {
		t.Fatalf("bad user entries ID, expected %d, actual %d", 1, userentries[0].ID)
	}
	if userentries[1].ID != 2 {
		t.Fatalf("bad user entries ID, expected %d, actual %d", 2, userentries[1].ID)
	}

	// test by feed
	userfeeds, err := db.UserFeedGetAllByUser(users[1].ID)
	if err != nil {
		t.Fatalf("Cannot retrieve user feeds: %s", err.Error())
	} else if len(userfeeds) != 2 {
		t.Fatalf("Bad userfeed count: expected %d, actual %d", 2, len(userfeeds))
	}

	// unread test feed

	if count, err := countUnreadEntries(db, users[1].ID, userfeeds[0].ID); err != nil {
		t.Fatalf("Cannot count unread entries: %s", err.Error())
	} else if count != 50 {
		t.Errorf("Bad userentry count: expected %d, actual %d", 50, count)
	}

	if err = db.UserEntryUpdateReadByFeed(users[1].ID, userfeeds[0].ID, time.Now(), true); err != nil {
		t.Errorf("Cannot update user entries for feed: %s", err.Error())
	}

	if count, err := countUnreadEntries(db, users[1].ID, userfeeds[0].ID); err != nil {
		t.Fatalf("Cannot count unread entries: %s", err.Error())
	} else if count != 0 {
		t.Errorf("Bad userentry count: expected %d, actual %d", 0, count)
	}

	// unread test all

	if count, err := countUnreadEntries(db, users[1].ID, 0); err != nil {
		t.Fatalf("Cannot count unread entries: %s", err.Error())
	} else if count != 98 {
		t.Errorf("Bad userentry count: expected %d, actual %d", 98, count)
	}

	if err = db.UserEntryUpdateReadByFeed(users[1].ID, 0, time.Now(), true); err != nil {
		t.Errorf("Cannot update user entries for feed: %s", err.Error())
	}

	if count, err := countUnreadEntries(db, users[1].ID, 0); err != nil {
		t.Fatalf("Cannot count unread entries: %s", err.Error())
	} else if count != 0 {
		t.Errorf("Bad userentry count: expected %d, actual %d", 0, count)
	}

}

func countUnreadEntries(db *Service, userID, userfeedID uint64) (int, error) {
	userentries, err := db.UserEntryGetUnreadForUser(userID)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, userentry := range userentries {
		if userfeedID == 0 || userentry.UserFeedID == userfeedID {
			count++
		}
	}
	return count, nil
}
