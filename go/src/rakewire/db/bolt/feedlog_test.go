package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire/db"
	m "rakewire/model"
	"testing"
	"time"
)

func TestFeedLog(t *testing.T) {

	t.SkipNow()

	database := Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	now := time.Now().Truncate(time.Second)
	feedID := "12345"

	err = database.db.Update(func(tx *bolt.Tx) error {
		for i := 1; i <= 100; i++ {
			dt := now.Add(time.Hour * time.Duration(-i))
			entry := &m.FeedLog{}
			entry.StartTime = dt
			entry.Duration = time.Duration(i)
			// FIXME: needs id for feed log
			err := marshal(entry, feedID, tx)
			assert.Nil(t, err)
		}
		return nil
	})
	assert.Nil(t, err)

	entries, err := database.GetFeedLog(feedID, 10*time.Hour)
	assert.Nil(t, err)
	assert.NotNil(t, entries)
	assert.Equal(t, 10, len(entries))

	// test reverse chronological order
	assert.Equal(t, time.Duration(1), entries[0].Duration)
	assert.Equal(t, time.Duration(10), entries[9].Duration)

	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}
