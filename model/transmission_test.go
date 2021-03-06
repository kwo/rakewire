package model

import (
	"testing"
	"time"
)

func TestTransmissionSetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityTransmission); obj == nil {
		t.Error("missing getObject entry")
	} else if !obj.hasIncrementingID() {
		t.Error("transmissions have incrementing IDs")
	}

	if obj := allEntities[entityTransmission]; obj == nil {
		t.Error("missing allEntities entry")
	}

}

func TestTransmissionIncrementingID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	if err := db.Update(func(tx Transaction) error {
		t := T.New("0000000001")
		return T.Save(tx, t)
	}); err != nil {
		t.Fatalf("error adding new transmission: %s", err.Error())
	}

	if err := db.Update(func(tx Transaction) error {
		id, err := tx.NextID(entityTransmission)
		if err != nil {
			return err
		}
		if id != 2 {
			t.Errorf("bad next id: %d, expected %d", id, 2)
		}
		return nil
	}); err != nil {
		t.Fatalf("error getting next id: %s", err.Error())
	}

}

func TestTransmissionGetBadID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	err := db.Select(func(tx Transaction) error {
		if transmission := T.Get(tx, empty); transmission != nil {
			t.Error("Expected nil transmission")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting from database: %s", err.Error())
	}

}

func TestTransmissions(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	transmissionID := empty
	feedID := "0000000002"

	// add transmission
	err := db.Update(func(tx Transaction) error {

		transmission := T.New(feedID)
		if err := T.Save(tx, transmission); err != nil {
			return err
		}
		transmissionID = transmission.ID

		return nil

	})
	if err != nil {
		t.Errorf("Error adding transmission: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {

		transmission := T.Get(tx, transmissionID)
		if transmission == nil {
			t.Fatal("Nil transmission, expected valid transmission")
		}
		if transmission.ID != transmissionID {
			t.Errorf("transmission ID mismatch, expected %s, actual %s", transmissionID, transmission.ID)
		}
		if transmission.FeedID != feedID {
			t.Errorf("transmission FeedID mismatch, expected %s, actual %s", feedID, transmission.URL)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting transmission: %s", err.Error())
	}

	// delete transmission
	err = db.Update(func(tx Transaction) error {
		if err := T.Delete(tx, transmissionID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error deleting transmission: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {
		transmission := T.Get(tx, transmissionID)
		if transmission != nil {
			t.Error("Expected nil transmission")
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error selecting transmission: %s", err.Error())
	}

}

func TestTransmissionRanges(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	feedID1 := "0000000010"
	feedID2 := "0000000020"
	now := time.Now().Truncate(time.Second)

	// add transmissions
	err := db.Update(func(tx Transaction) error {

		latest := now.Add(-10 * 100 * time.Minute)

		for i := 0; i < 100; i++ {
			var feedID string
			if i%2 == 0 {
				feedID = feedID2
			} else {
				feedID = feedID1
			}
			latest = latest.Add(10 * time.Minute)
			transmission := T.New(feedID)
			transmission.StartTime = latest
			if err := T.Save(tx, transmission); err != nil {
				return err
			}
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error adding transmissions: %s", err.Error())
	}

	// get last transmission
	err = db.Select(func(tx Transaction) error {

		transmission := T.GetLast(tx)

		if transmission == nil {
			t.Fatalf("Cannot get last transmission")
		}

		if transmission.ID != "0000000100" {
			t.Errorf("Bad transmissionID: %s", transmission.ID)
		}

		if transmission.FeedID != "0000000010" {
			t.Errorf("Bad transmission FeedID: %s", transmission.FeedID)
		}

		if transmission.StartTime.Unix() != now.Unix() {
			t.Errorf("Bad transmission StartTime: %v, expected: %v", transmission.StartTime, now)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting transmissions: %s", err.Error())
	}

	// get feed2 in the last hour transmissions
	err = db.Select(func(tx Transaction) error {

		transmissions := T.GetForFeed(tx, feedID2, time.Hour)

		if transmissions == nil {
			t.Fatalf("Cannot get last transmissions")
		}

		if len(transmissions) != 3 {
			t.Errorf("Bad transmission count: %d, expected %d", len(transmissions), 3)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting transmissions: %s", err.Error())
	}

	// get all in the last hour
	err = db.Select(func(tx Transaction) error {

		transmissions := T.GetRange(tx, now, time.Hour)

		if transmissions == nil {
			t.Fatalf("Cannot get last transmissions")
		}

		if len(transmissions) != 6 {
			t.Errorf("Bad transmission count: %d, expected %d", len(transmissions), 6)
		}

		return nil

	})
	if err != nil {
		t.Errorf("Error selecting transmissions: %s", err.Error())
	}

}
