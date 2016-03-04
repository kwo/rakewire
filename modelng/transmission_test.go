package modelng

import (
	"testing"
	"time"
)

func TestTransmissionSetup(t *testing.T) {

	t.Parallel()

	if obj := getObject(entityTransmission); obj == nil {
		t.Error("missing getObject entry")
	}

	if obj := allEntities[entityTransmission]; obj == nil {
		t.Error("missing allEntities entry")
	}

	c := &Config{}
	if obj := c.Sequences.Transmission; obj != 0 {
		t.Error("missing sequences entry")
	}

}

func TestTransmissionGetBadID(t *testing.T) {

	t.Parallel()

	db := openTestDatabase(t)
	defer closeTestDatabase(t, db)

	err := db.Select(func(tx Transaction) error {
		if transmission := T.Get(empty, tx); transmission != nil {
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
		if err := T.Save(transmission, tx); err != nil {
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

		transmission := T.Get(transmissionID, tx)
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
		if err := T.Delete(transmissionID, tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Errorf("Error deleting transmission: %s", err.Error())
	}

	// test by id
	err = db.Select(func(tx Transaction) error {
		transmission := T.Get(transmissionID, tx)
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
			if err := T.Save(transmission, tx); err != nil {
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

		transmissions := T.GetForFeed(feedID2, time.Hour, tx)

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

		transmissions := T.GetRange(now, time.Hour, tx)

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
