package bolt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire.com/db"
	m "rakewire.com/model"
	"testing"
	"time"
)

func TestSaveFeedLog(t *testing.T) {

	database := Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	now := time.Now().Truncate(time.Second)

	for i := 1; i <= 100; i++ {
		dt := now.Add(time.Hour * time.Duration(-i))
		entry := &m.FeedLog{}
		entry.FeedID = "12345"
		entry.StartTime = &dt
		entry.Duration = time.Duration(i)
		err = database.SaveFeedLog(entry)
		assert.Nil(t, err)
	}

	entries, err := database.GetFeedLog("12345", 10*time.Hour)
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
