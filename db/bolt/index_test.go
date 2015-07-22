package bolt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire.com/db"
	"testing"
	"time"
)

func TestIndexFetch(t *testing.T) {

	//t.SkipNow()

	// open database
	database := &Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	// test feeds
	feeds, err := database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 0, feeds.Size())

	maxTime := time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(&maxTime)
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 0, feeds.Size())

	// create new feed, add to database
	feeds = db.NewFeeds()
	feed := db.NewFeed("http://localhost/")
	feeds.Add(feed)
	err = database.SaveFeeds(feeds)
	assert.Nil(t, err)

	// retest
	feeds, err = database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 1, feeds.Size())

	maxTime = time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(&maxTime)
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 1, feeds.Size())

	// modify feed, resave to database
	lastTime := time.Now().Add(-48 * time.Hour)
	// create new feed, add to database
	feeds2 := db.NewFeeds()
	feed2 := &db.Feed{
		ID:        feed.ID,
		URL:       "https://localhost/",
		LastFetch: &lastTime,
	}
	feeds2.Add(feed2)
	f3 := db.NewFeed("http://kangaroo.com/")
	feeds2.Add(f3)
	err = database.SaveFeeds(feeds2)
	assert.Nil(t, err)

	// retest
	feeds, err = database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 2, feeds.Size())

	maxTime = time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(&maxTime)
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 2, feeds.Size())

	// close database
	err = database.Close()
	assert.Nil(t, err)

	// delete test file
	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}

func TestIndexDeletes(t *testing.T) {

	//t.SkipNow()

	// open database
	database := &Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	urls, err := readFile(feedFile)
	require.Nil(t, err)
	require.NotNil(t, urls)

	now := time.Now().Add(-1 * time.Hour)
	feeds := db.NewFeeds()
	for _, url := range urls {
		feed := db.NewFeed(url)
		feed.LastFetch = &now
		feeds.Add(feed)
	}

	err = database.SaveFeeds(feeds)
	assert.Nil(t, err)

	feeds, err = database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 288, feeds.Size())

	for _, feed := range feeds.Values {
		err = database.checkIndexForDuplicates(bucketIndexNextFetch, feed.ID)
		assert.Nil(t, err)
	}

	now = time.Now()
	for _, feed := range feeds.Values {
		feed.LastFetch = &now
		feed.Frequency = 10
	}

	err = database.SaveFeeds(feeds)
	assert.Nil(t, err)

	feeds, err = database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 288, feeds.Size())

	maxTime := time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(&maxTime)
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 288, feeds.Size())

	for _, feed := range feeds.Values {
		err = database.checkIndexForDuplicates(bucketIndexNextFetch, feed.ID)
		assert.Nil(t, err)
	}

	// close database
	err = database.Close()
	assert.Nil(t, err)

	// delete test file
	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}
