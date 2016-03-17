package model

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestEntrySetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityEntry); obj == nil {
		t.Error("missing getObject entry")
	}

	if obj := allEntities[entityEntry]; obj == nil {
		t.Error("missing allEntities entry")
	}

}

func TestIndex(t *testing.T) {

	t.Parallel()

	r1 := []byte("0000000002|0000000012|20160309210900Z|0000000230") // -1
	mx := []byte("0000000002|0000000012|20160309211000Z")
	r2 := []byte("0000000002|0000000012|20160309211000Z|0000000230") // 1

	var expectedResult int

	expectedResult = -1
	if c := bytes.Compare(r1, mx); c != expectedResult {
		t.Errorf("Bad compare result: %d, expected %d", c, expectedResult)
	}

	expectedResult = 1
	if c := bytes.Compare(r2, mx); c != expectedResult {
		t.Errorf("Bad compare result: %d, expected %d", c, expectedResult)
	}

}

func TestEntryBasics(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	userID := "0000000001"
	itemID := "0000000002"
	feedID := "0000000003"

	// add entry
	err := db.Update(func(tx Transaction) error {

		entry := E.New(userID, itemID, feedID)
		if err := E.Save(tx, entry); err != nil {
			return err
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error adding entry: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {

		entry := E.Get(tx, userID, itemID)
		if entry == nil {
			t.Fatal("Nil entry, expected valid entry")
		}
		if entry.UserID != userID {
			t.Errorf("bad userID: %s, expected %s", entry.UserID, userID)
		}
		if entry.ItemID != itemID {
			t.Errorf("bad itemID: %s, expected %s", entry.ItemID, itemID)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting entry: %s", err.Error())
	}

	// delete entry
	err = db.Update(func(tx Transaction) error {
		if err := E.Delete(tx, keyEncode(userID, itemID)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error deleting entry: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {
		entry := E.Get(tx, userID, itemID)
		if entry != nil {
			t.Error("Expected nil entry")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting entry: %s", err.Error())
	}

}

func TestEntryQueries(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	t0 := time.Now().Truncate(time.Hour).Add(-24 * time.Hour)
	feeds := Feeds{}
	users := Users{}
	items := Items{}

	// add entries
	err := db.Update(func(tx Transaction) error {

		// add users
		for uu := 0; uu < 2; uu++ {
			u := U.New(fmt.Sprintf("User%02d", uu+1))
			u.SetPassword("ancdefg")
			if err := U.Save(tx, u); err != nil {
				return err
			}
			users = append(users, u)
		}

		// add feeds
		for ff := 0; ff < 20; ff++ {
			f := F.New(fmt.Sprintf("Feed%02d", ff+1))
			if err := F.Save(tx, f); err != nil {
				return err
			}
			feeds = append(feeds, f)
		}

		// add subscriptions
		for _, user := range users {
			for _, feed := range feeds {
				s := S.New(user.ID, feed.ID)
				if err := S.Save(tx, s); err != nil {
					return err
				}
			}
		}

		// add items
		// there are 20 feeds with 20 items each
		// feeds start on the hour, the first one at t0 - 24hrs
		// items increment by 1 minute, starting on the hour
		// Feed 0000000001 : -24:00:00 to -23:41:00 +00
		// Feed 0000000002 : -23:00:00 to -22:41:00 +01
		// Feed 0000000003 : -22:00:00 to -21:41:00 +02
		// Feed 0000000004 : -21:00:00 to -20:41:00 +03
		// Feed 0000000005 : -20:00:00 to -19:41:00 +04
		// Feed 0000000006 : -19:00:00 to -18:41:00 +05
		// Feed 0000000007 : -18:00:00 to -17:41:00 +06
		// Feed 0000000008 : -17:00:00 to -16:41:00 +07
		// Feed 0000000009 : -16:00:00 to -15:41:00 +08
		// Feed 0000000010 : -15:00:00 to -14:41:00 +09
		// Feed 0000000011 : -14:00:00 to -13:41:00 +10
		// Feed 0000000012 : -13:00:00 to -12:41:00 +11
		// Feed 0000000013 : -12:00:00 to -11:41:00 +12
		// Feed 0000000014 : -11:00:00 to -10:41:00 +13
		// Feed 0000000015 : -10:00:00 to -09:41:00 +14
		// Feed 0000000016 : -09:00:00 to -08:41:00 +15
		// Feed 0000000017 : -08:00:00 to -07:41:00 +16
		// Feed 0000000018 : -07:00:00 to -06:41:00 +17
		// Feed 0000000019 : -06:00:00 to -05:41:00 +18
		// Feed 0000000020 : -05:00:00 to -04:41:00 +19
		for f, feed := range feeds {
			lastUpdated := t0.Add(time.Duration(f) * time.Hour)
			for ii := 0; ii < 20; ii++ {
				lastUpdated := lastUpdated.Add(time.Duration(ii) * time.Minute)
				item := I.New(feed.ID, fmt.Sprintf("Feed%02dItem%02d", f+1, ii+1))
				item.Updated = lastUpdated
				if err := I.Save(tx, item); err != nil {
					return err
				}
				//t.Logf("Item: %s - %s, %v", item.ID, item.FeedID, item.Updated)
				items = append(items, item)
			}
		}

		return E.AddItems(tx, items)

	})
	if err != nil {
		t.Errorf("Error updating database: %s", err.Error())
	}

	// test counts
	err = db.Select(func(tx Transaction) error {

		// test count
		expectedCount := 400
		for _, user := range users {
			if count := E.Query(tx, user.ID).Count(); count != expectedCount {
				t.Errorf("Bad count: %d, expected %d", count, expectedCount)
			}
		}

		// test count with timeframe
		expectedCount = 120
		min := t0.Add(6 * time.Hour)
		max := t0.Add(12 * time.Hour)
		for _, user := range users {
			if count := E.Query(tx, user.ID).Min(min).Max(max).Count(); count != expectedCount {
				t.Errorf("Bad count: %d, expected %d", count, expectedCount)
			}
		}

		// test count with feed
		expectedCount = 20
		for _, user := range users {
			for _, feed := range feeds {
				if count := E.Query(tx, user.ID).Feed(feed.ID).Count(); count != expectedCount {
					t.Errorf("Bad count: %d, expected %d", count, expectedCount)
				}
			}
		}

		// test count with timeframe, with feed
		expectedCount = 5
		min = t0.Add(11 * time.Hour).Add(5 * time.Minute)
		max = t0.Add(11 * time.Hour).Add(10 * time.Minute)
		for _, user := range users {
			if count := E.Query(tx, user.ID).Feed("0000000012").Min(min).Max(max).Count(); count != expectedCount {
				t.Errorf("Bad count: %d, expected %d", count, expectedCount)
			}
		}

		// test get with timeframe, with feed + 1 second to show max exclusivity
		expectedCount = 6
		min = t0.Add(11 * time.Hour).Add(5 * time.Minute)
		max = t0.Add(11 * time.Hour).Add(10 * time.Minute).Add(1 * time.Second)
		t.Logf("min: %v, max: %v", min, max)
		for _, user := range users {
			entries := E.Query(tx, user.ID).Feed("0000000012").Min(min).Max(max).Get()
			if len(entries) != expectedCount {
				t.Errorf("Bad count: %d, expected %d", len(entries), expectedCount)
				for _, entry := range entries {
					t.Logf("entry: %s: %v", entry.GetID(), entry.Updated)
				}
			}
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

	// test gets all
	err = db.Select(func(tx Transaction) error {

		// get entries in updated order
		expectedCount := 360
		min := t0.Add(1 * time.Hour)
		max := t0.Add(19 * time.Hour)
		for _, user := range users {
			if entries := E.Query(tx, user.ID).Min(min).Max(max).Get(); entries != nil {

				if len(entries.Unique()) != expectedCount {
					t.Errorf("Bad count: %d, expected %d", len(entries), expectedCount)
				}

				entries.Reverse()
				lastUpdated := time.Now()
				itemsByID := items.ByID()
				for _, entry := range entries {
					if entry.Updated.After(lastUpdated) {
						t.Errorf("Entries not sorted, updated %s, should be before %s", entry.Updated.Format(fmtTime), lastUpdated.Format(fmtTime))
					}
					lastUpdated = entry.Updated
					item := itemsByID[entry.ItemID]
					if !entry.Updated.Equal(item.Updated) {
						t.Errorf("Bad updated: %s, expected: %s", entry.Updated.Format(fmtTime), item.Updated.Format(fmtTime))
					}
				}

			} else {
				t.Fatal("nil entries")
			}
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

	// test gets per feed
	err = db.Select(func(tx Transaction) error {

		// get entries in updated order
		expectedCount := 20
		for _, user := range users {
			for _, feed := range feeds {

				if entries := E.Query(tx, user.ID).Feed(feed.ID).Get(); entries != nil {

					if len(entries.Unique()) != expectedCount {
						t.Errorf("Bad count: %d, expected %d", len(entries), expectedCount)
					}

					entries.Reverse()
					lastUpdated := time.Now()
					itemsByID := items.ByID()
					for _, entry := range entries {
						if entry.FeedID != feed.ID {
							t.Errorf("Bad feedID %s, expected %s", entry.FeedID, feed.ID)
						}
						if entry.Updated.After(lastUpdated) {
							t.Errorf("Entries not sorted, updated %s, should be before %s", entry.Updated.Format(fmtTime), lastUpdated.Format(fmtTime))
						}
						lastUpdated = entry.Updated
						item := itemsByID[entry.ItemID]
						if !entry.Updated.Equal(item.Updated) {
							t.Errorf("Bad updated: %s, expected: %s", entry.Updated.Format(fmtTime), item.Updated.Format(fmtTime))
						}
					}

				} else {
					t.Fatal("nil entries")
				}
			}
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestStarred(t *testing.T) {

	/*

		User Star Updated Item
		0000000001 0 t0 0000000001
		0000000001 0 t1 0000000002
		0000000001 0 t1 0000000006
		0000000001 1 t0 0000000003
		0000000001 1 t1 0000000004
		0000000001 1 t1 0000000005

		User Feed Star Updated Item
		0000000001 0000000002 0 t0 0000000001
		0000000001 0000000002 0 t1 0000000002
		0000000001 0000000002 1 t0 0000000003
		0000000001 0000000002 1 t1 0000000004
		0000000001 0000000003 1 t1 0000000005
		0000000001 0000000004 0 t1 0000000006

	*/

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	now := time.Now().Truncate(time.Second)
	t0 := now.Add(-15 * time.Minute)
	t1 := now.Add(-10 * time.Minute)

	err := db.Update(func(tx Transaction) error {

		entries := Entries{}
		var entry *Entry

		entry = E.New(keyEncodeUint(1), keyEncodeUint(1), keyEncodeUint(2))
		entry.Updated = t0
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(2), keyEncodeUint(2))
		entry.Updated = t1
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(3), keyEncodeUint(2))
		entry.Updated = t0
		entry.Star = true
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(4), keyEncodeUint(2))
		entry.Updated = t1
		entry.Star = true
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(5), keyEncodeUint(3))
		entry.Updated = t1
		entry.Star = true
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(6), keyEncodeUint(4))
		entry.Updated = t1
		entries = append(entries, entry)

		return E.SaveAll(tx, entries)

	})
	if err != nil {
		t.Errorf("Error updating database: %s", err.Error())
	}

	err = db.Select(func(tx Transaction) error {

		// get total starred
		if entries := E.Query(tx, keyEncodeUint(1)).Starred(); entries != nil {
			expectedCount := 3
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total starred for time up til t1
		if entries := E.Query(tx, keyEncodeUint(1)).Max(t1).Starred(); entries != nil {
			expectedCount := 1
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total starred for time starting at t1
		if entries := E.Query(tx, keyEncodeUint(1)).Min(t1).Starred(); entries != nil {
			expectedCount := 2
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total starred for feed 2
		if entries := E.Query(tx, keyEncodeUint(1)).Feed(keyEncodeUint(2)).Starred(); entries != nil {
			expectedCount := 2
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total starred for feed 2 up til t1
		if entries := E.Query(tx, keyEncodeUint(1)).Feed(keyEncodeUint(2)).Max(t1).Starred(); entries != nil {
			expectedCount := 1
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total starred for feed 2 starting at t1
		if entries := E.Query(tx, keyEncodeUint(1)).Feed(keyEncodeUint(2)).Min(t1).Starred(); entries != nil {
			expectedCount := 1
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestUnread(t *testing.T) {

	/*

		User Read Updated Item
		0000000001 0 t0 0000000001
		0000000001 0 t1 0000000002
		0000000001 0 t1 0000000006
		0000000001 1 t0 0000000003
		0000000001 1 t1 0000000004
		0000000001 1 t1 0000000005

		User Feed Read Updated Item
		0000000001 0000000002 0 t0 0000000001
		0000000001 0000000002 0 t1 0000000002
		0000000001 0000000002 1 t0 0000000003
		0000000001 0000000002 1 t1 0000000004
		0000000001 0000000003 1 t1 0000000005
		0000000001 0000000004 0 t1 0000000006

	*/

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	now := time.Now().Truncate(time.Second)
	t0 := now.Add(-15 * time.Minute)
	t1 := now.Add(-10 * time.Minute)

	err := db.Update(func(tx Transaction) error {

		entries := Entries{}
		var entry *Entry

		entry = E.New(keyEncodeUint(1), keyEncodeUint(1), keyEncodeUint(2))
		entry.Updated = t0
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(2), keyEncodeUint(2))
		entry.Updated = t1
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(3), keyEncodeUint(2))
		entry.Updated = t0
		entry.Read = true
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(4), keyEncodeUint(2))
		entry.Updated = t1
		entry.Read = true
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(5), keyEncodeUint(3))
		entry.Updated = t1
		entry.Read = true
		entries = append(entries, entry)

		entry = E.New(keyEncodeUint(1), keyEncodeUint(6), keyEncodeUint(4))
		entry.Updated = t1
		entries = append(entries, entry)

		return E.SaveAll(tx, entries)

	})
	if err != nil {
		t.Errorf("Error updating database: %s", err.Error())
	}

	err = db.Select(func(tx Transaction) error {

		// get total unread
		if entries := E.Query(tx, keyEncodeUint(1)).Unread(); entries != nil {
			expectedCount := 3
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total unread for time up til t1
		if entries := E.Query(tx, keyEncodeUint(1)).Max(t1).Unread(); entries != nil {
			expectedCount := 1
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total unread for time starting at t1
		if entries := E.Query(tx, keyEncodeUint(1)).Min(t1).Unread(); entries != nil {
			expectedCount := 2
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total unread for feed 2
		if entries := E.Query(tx, keyEncodeUint(1)).Feed(keyEncodeUint(2)).Unread(); entries != nil {
			expectedCount := 2
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total unread for feed 2 up til t1
		if entries := E.Query(tx, keyEncodeUint(1)).Feed(keyEncodeUint(2)).Max(t1).Unread(); entries != nil {
			expectedCount := 1
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		// get total unread for feed 2 starting at t1
		if entries := E.Query(tx, keyEncodeUint(1)).Feed(keyEncodeUint(2)).Min(t1).Unread(); entries != nil {
			expectedCount := 1
			if len(entries) != expectedCount {
				t.Errorf("Bad entry count: %d, expected %d", len(entries), expectedCount)
			}
		} else {
			t.Error("Nil entries")
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}
