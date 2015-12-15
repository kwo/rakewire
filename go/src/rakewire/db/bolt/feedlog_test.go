package bolt

import (
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"testing"
	"time"
)

func TestFeedLog(t *testing.T) {

	// t.SkipNow()

	database := openDatabase(t)
	defer closeDatabase(t, database)
	assertNotNil(t, database)

	now := time.Now().Truncate(time.Second)
	feedID := "12345"

	err := database.db.Update(func(tx *bolt.Tx) error {
		for i := 0; i <= 100; i++ {
			dt := now.Add(time.Hour * time.Duration(-i))
			entry := m.NewFeedLog(feedID)
			entry.StartTime = dt
			entry.Duration = time.Duration(i)
			if err := kvSave(entry, tx); err != nil {
				return err
			}
		}
		return nil
	})
	assertNoError(t, err)

	entries, err := database.GetFeedLog(feedID, 10*time.Hour)
	assertNoError(t, err)
	assertNotNil(t, entries)

	if len(entries) != 11 {
		t.Fatalf("bad entry count, expected %d, actual %d", 11, len(entries))
	}

	// test reverse chronological order
	assertEqual(t, time.Duration(0), entries[0].Duration)
	assertEqual(t, time.Duration(10), entries[10].Duration)

}
