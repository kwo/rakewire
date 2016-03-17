package model

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestFeedSetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityFeed); obj == nil {
		t.Error("missing getObject entry")
	}

	if obj := allEntities[entityFeed]; obj == nil {
		t.Error("missing allEntities entry")
	}

	c := &Config{}
	if obj := c.Sequences.Feed; obj != 0 {
		t.Error("missing sequences entry")
	}

}

func TestFeedGetBadID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	err := db.Select(func(tx Transaction) error {
		if feed := F.Get(tx, empty); feed != nil {
			t.Error("Expected nil feed")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestFeedGetBadURL(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	err := db.Select(func(tx Transaction) error {
		if feed := F.GetByURL(tx, empty); feed != nil {
			t.Error("Expected nil feed")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestFeeds(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	feedID := empty
	feedURL := "https://google.com/feed"

	// add feed
	err := db.Update(func(tx Transaction) error {

		feed := F.New(feedURL)
		if err := F.Save(tx, feed); err != nil {
			return err
		}
		feedID = feed.ID

		return nil

	})
	if err != nil {
		t.Errorf("Error adding feed: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {

		feed := F.Get(tx, feedID)
		if feed == nil {
			t.Fatal("Nil feed, expected valid feed")
		}
		if feed.ID != feedID {
			t.Errorf("feed ID mismatch, expected %s, actual %s", feedID, feed.ID)
		}
		if feed.URL != feedURL {
			t.Errorf("feed URL mismatch, expected %s, actual %s", feedURL, feed.URL)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting feed: %s", err.Error())
	}

	// test by url
	err = db.Select(func(tx Transaction) error {

		feed := F.GetByURL(tx, feedURL)
		if feed == nil {
			t.Fatal("Nil feed, expected valid feed")
		}
		if feed.ID != feedID {
			t.Errorf("feed ID mismatch, expected %s, actual %s", feedID, feed.ID)
		}
		if feed.URL != feedURL {
			t.Errorf("feed URL mismatch, expected %s, actual %s", feedURL, feed.URL)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting feed: %s", err.Error())
	}

	// test by url uppercase
	err = db.Select(func(tx Transaction) error {

		feed := F.GetByURL(tx, strings.ToUpper(feedURL))
		if feed == nil {
			t.Fatal("Nil feed, expected valid feed")
		}
		if feed.ID != feedID {
			t.Errorf("feed ID mismatch, expected %s, actual %s", feedID, feed.ID)
		}
		if feed.URL != feedURL {
			t.Errorf("feed URL mismatch, expected %s, actual %s", feedURL, feed.URL)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting feed: %s", err.Error())
	}

	// delete feed
	err = db.Update(func(tx Transaction) error {
		if err := F.Delete(tx, feedID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error deleting feed: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {
		feed := F.Get(tx, feedID)
		if feed != nil {
			t.Error("Expected nil feed")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting feed: %s", err.Error())
	}

}

func TestFeedGetNext(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	// add 20 feeds with LastUpdated incrementing in 3 minute intervals, starting now
	err := db.Update(func(tx Transaction) error {
		latest := time.Now().Truncate(time.Second)
		for i := 0; i < 20; i++ {
			feed := F.New(fmt.Sprintf("Feed%02d", i+1))
			feed.LastUpdated = latest
			feed.UpdateFetchTime(feed.LastUpdated)
			if err := F.Save(tx, feed); err != nil {
				return err
			}
			t.Logf("Feed Add : %s %s, lastUpdated: %s, nextFetch: %s", feed.ID, feed.URL, keyEncodeTime(feed.LastUpdated), keyEncodeTime(feed.NextFetch))
			latest = latest.Add(-5 * time.Minute)
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error updating database: %s", err.Error())
	}

	err = db.Select(func(tx Transaction) error {
		b := tx.Bucket(bucketIndex, entityFeed, indexFeedNextFetch)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			t.Logf("k: %s, v: %s", k, v)
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

	err = db.Select(func(tx Transaction) error {
		max := time.Now().Add(10 * time.Minute).Truncate(time.Second)
		feeds := F.GetNext(tx, max)
		for _, feed := range feeds {
			t.Logf("Feed Next: %s, lastUpdated: %s, nextFetch: %s", feed.URL, keyEncodeTime(feed.LastUpdated), keyEncodeTime(feed.NextFetch))
		}

		if len(feeds) != 7 {
			t.Errorf("Feed count mismatch, expected %d, actual %d", 7, len(feeds))
		}

		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestAdjustFetchTime(t *testing.T) {

	t.Parallel()

	f := F.New("http://localhost")
	if f == nil {
		t.Fatal("F.New returned a nil feed")
	}
	if f.NextFetch.IsZero() {
		t.Fatalf("F.New must set NextFetch to now, actual: %v", f.NextFetch)
	}

	now := time.Now().Truncate(time.Second)
	f.NextFetch = now

	diff := 3 * time.Hour
	f.AdjustFetchTime(diff)

	expectedNextFetch := now.Add(diff)

	if !f.NextFetch.Equal(expectedNextFetch) {
		t.Errorf("Adjust fetch time error, expected %v, actual %v", expectedNextFetch, f.NextFetch)
	}

}

func TestUpdateFetchTime1(t *testing.T) {

	t.Parallel()

	now := time.Now().Truncate(time.Second)

	f := F.New("http://localhost")
	if f == nil {
		t.Fatal("F.New returned a nil feed")
	}
	if f.NextFetch.IsZero() {
		t.Fatalf("F.New must set NextFetch to now, actual: %v", f.NextFetch)
	}

	f.LastUpdated = now.Add(-47 * time.Second)
	f.UpdateFetchTime(f.LastUpdated)

	expectedTime := f.LastUpdated.Add(15 * time.Minute)
	t.Logf("now:         %v", now)
	t.Logf("lastUpdated: %v", f.LastUpdated)
	t.Logf("nextFetch:   %v", f.NextFetch)
	t.Logf("expected:    %v", expectedTime)

	if !f.NextFetch.Equal(expectedTime) {
		t.Errorf("bad fetch time, expected %v from now, actual %v", expectedTime, f.NextFetch)
	}

}

func TestUpdateFetchTime2(t *testing.T) {

	t.Parallel()

	now := time.Now().Truncate(time.Second)

	f := F.New("http://localhost")
	if f == nil {
		t.Fatal("F.New returned a nil feed")
	}
	if f.NextFetch.IsZero() {
		t.Fatalf("F.New must set NextFetch to now, actual: %v", f.NextFetch)
	}

	f.LastUpdated = now.Add(-29 * time.Minute)
	f.UpdateFetchTime(f.LastUpdated)

	expectedTime := f.LastUpdated.Add(45 * time.Minute)
	t.Logf("now:         %v", now)
	t.Logf("lastUpdated: %v", f.LastUpdated)
	t.Logf("nextFetch:   %v", f.NextFetch)
	t.Logf("expected:    %v", expectedTime)

	if !f.NextFetch.Equal(expectedTime) {
		t.Errorf("bad fetch time, expected %v from now, actual %v", expectedTime, f.NextFetch)
	}

}

func TestUpdateFetchTime3(t *testing.T) {

	t.Parallel()

	now := time.Now().Truncate(time.Second)

	f := F.New("http://localhost")
	if f == nil {
		t.Fatal("F.New returned a nil feed")
	}
	if f.NextFetch.IsZero() {
		t.Fatalf("F.New must set NextFetch to now, actual: %v", f.NextFetch)
	}

	f.LastUpdated = now.Add(-3 * time.Hour).Add(-47 * time.Minute)
	f.UpdateFetchTime(f.LastUpdated)

	expectedTime := f.LastUpdated.Add(4 * time.Hour)
	t.Logf("now:         %v", now)
	t.Logf("lastUpdated: %v", f.LastUpdated)
	t.Logf("nextFetch:   %v", f.NextFetch)
	t.Logf("expected:    %v", expectedTime)

	if !f.NextFetch.Equal(expectedTime) {
		t.Errorf("bad fetch time, expected %v from now, actual %v", expectedTime, f.NextFetch)
	}

}

func TestUpdateFetchTime4(t *testing.T) {

	t.Parallel()

	now := time.Now().Truncate(time.Second)

	f := F.New("http://localhost")
	if f == nil {
		t.Fatal("F.New returned a nil feed")
	}
	if f.NextFetch.IsZero() {
		t.Fatalf("F.New must set NextFetch to now, actual: %v", f.NextFetch)
	}

	f.LastUpdated = now.Add(-4 * 24 * time.Hour)
	f.UpdateFetchTime(f.LastUpdated)

	expectedTime := f.LastUpdated.Add(((4 * 24) + 1) * time.Hour)
	t.Logf("now:         %v", now)
	t.Logf("lastUpdated: %v", f.LastUpdated)
	t.Logf("nextFetch:   %v", f.NextFetch)
	t.Logf("expected:    %v", expectedTime)

	if !f.NextFetch.Equal(expectedTime) {
		t.Errorf("bad fetch time, expected %v from now, actual %v", expectedTime, f.NextFetch)
	}

}
