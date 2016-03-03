package modelng

import (
	"strings"
	"testing"
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
		if feed := F.Get(empty, tx); feed != nil {
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
		if feed := F.GetByURL(empty, tx); feed != nil {
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
		if err := F.Save(feed, tx); err != nil {
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

		feed := F.Get(feedID, tx)
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

		feed := F.GetByURL(feedURL, tx)
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

		feed := F.GetByURL(strings.ToUpper(feedURL), tx)
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
		if err := F.Delete(feedID, tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error deleting feed: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {
		feed := F.Get(feedID, tx)
		if feed != nil {
			t.Error("Expected nil feed")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting feed: %s", err.Error())
	}

}
