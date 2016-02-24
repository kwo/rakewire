package model

import (
	"testing"
	"time"
)

func TestEntryIndexes(t *testing.T) {

	e := NewEntry("0000000002", "0000000003", "0000000004")
	e.ID = "0000000001"
	e.Updated = time.Now().Truncate(time.Second)
	updatedStr := kvKeyTimeEncode(e.Updated)

	indexes := e.serializeIndexes()

	if len(indexes) != 3 {
		t.Fatalf("invalid number of indexes, expected %d, actual %d", 3, len(indexes))
	}

	for k, record := range indexes {

		t.Logf("%s: %v", k, record)

		if len(record) != 1 {
			t.Errorf("invalid number of record entries, expected %d, actual %d", 1, len(record))
		}

		switch k {

		case entryIndexRead:
			expectedKey := "0000000002.0." + updatedStr + ".0000000001"
			for key, value := range record {
				if key != expectedKey {
					t.Errorf("bad index key for %s: expected %s, actual %s", k, expectedKey, key)
				}
				if value != e.ID {
					t.Errorf("bad index value for %s: expected %s, actual %s", k, e.ID, value)
				}
			}

		case entryIndexStar:
			expectedKey := "0000000002.0." + updatedStr + ".0000000001"
			for key, value := range record {
				if key != expectedKey {
					t.Errorf("bad index key for %s: expected %s, actual %s", k, expectedKey, key)
				}
				if value != e.ID {
					t.Errorf("bad index value for %s: expected %s, actual %s", k, e.ID, value)
				}
			}

		case entryIndexUser:
			expectedKey := "0000000002.0000000001"
			for key, value := range record {
				if key != expectedKey {
					t.Errorf("bad index key for %s: expected %s, actual %s", k, expectedKey, key)
				}
				if value != e.ID {
					t.Errorf("bad index value for %s: expected %s, actual %s", k, e.ID, value)
				}
			}

		}

	}

}

func TestEntriesGetAll(t *testing.T) {

	t.Parallel()

	database := openTestDatabase(t, true)
	defer closeTestDatabase(t, database)

	err := database.Select(func(tx Transaction) error {

		user, err := UserByUsername(testUsername, tx)
		if err != nil {
			return err
		}

		entries, err := EntriesGetAll(user.ID, tx)
		if err != nil {
			return err
		}

		if entries == nil {
			t.Error("Nil entries")
		} else if len(entries) == 0 {
			t.Error("Zero entries")
		}

		return nil

	})

	if err != nil {
		t.Errorf("Error when selecting from database: %s", err.Error())
	}

}
