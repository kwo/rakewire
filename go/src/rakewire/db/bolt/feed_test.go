package bolt

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"rakewire/db"
	m "rakewire/model"
	"testing"
	"time"
)

func TestFeeds(t *testing.T) {

	//t.SkipNow()

	feeds, err := m.ParseFeedsFromFile(feedFile)
	require.Nil(t, err)
	require.NotNil(t, feeds)

	for _, feed := range feeds {
		// logger.Debugf("URL (%d): %s", n, feed.URL)
		feed.Attempt = m.NewFeedLog(feed.ID)
	}

	database := &Database{}
	err = database.Open(&db.Configuration{
		Location: getTempFile(t),
	})
	require.Nil(t, err)
	defer cleanUp(t, database)

	for _, feed := range feeds {
		err = database.SaveFeed(feed)
		assert.Nil(t, err)
	}

	feeds2, err := database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds2)

	assert.Equal(t, len(feeds), len(feeds2))
	// for n, feed := range feeds2 {
	// logger.Debugf("Feed (%d): %s", n, feed.URL)
	// }

}

func TestURLIndex(t *testing.T) {

	//t.SkipNow()

	const URL1 = "http://localhost/"
	const URL2 = "http://localhost:8888/"

	// open database
	database := &Database{}
	err := database.Open(&db.Configuration{
		Location: getTempFile(t),
	})
	require.Nil(t, err)
	defer cleanUp(t, database)

	// create feed
	feed := m.NewFeed(URL1)
	feed.Attempt = m.NewFeedLog(feed.ID)
	assert.Equal(t, URL1, feed.URL)
	assert.Equal(t, feed.ID, feed.Attempt.FeedID)

	// save feeds
	err = database.SaveFeed(feed)
	require.Nil(t, err)

	var feed2 *m.Feed
	var feeds2 []*m.Feed

	// get feed2
	feeds2, err = database.GetFeeds()
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	require.Equal(t, 1, len(feeds2))
	feed2 = feeds2[0]
	assert.NotNil(t, feed2)
	assert.Equal(t, feed.ID, feed2.ID)
	assert.Equal(t, URL1, feed2.URL)

	// get by URL
	feed2, err = database.GetFeedByURL(URL1)
	require.Nil(t, err)
	require.NotNil(t, feed2)
	assert.Equal(t, feed.ID, feed2.ID)
	assert.Equal(t, URL1, feed2.URL)

	// update URL
	feed2 = feeds2[0]
	feed2.URL = URL2
	err = database.SaveFeed(feed2)
	require.Nil(t, err)

	// get feeds2, feed2
	feeds2, err = database.GetFeeds()
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 1, len(feeds2))
	feed2 = feeds2[0]
	assert.NotNil(t, feed2)
	assert.Equal(t, feed.ID, feed2.ID)
	assert.Equal(t, URL2, feed2.URL)

	// get by old URL
	feed2, err = database.GetFeedByURL(URL1)
	require.Nil(t, err)
	require.Nil(t, feed2)

	// get by new URL
	feed2, err = database.GetFeedByURL(URL2)
	require.Nil(t, err)
	require.NotNil(t, feed2)
	assert.Equal(t, feed.ID, feed2.ID)
	assert.Equal(t, URL2, feed2.URL)

}

func TestIndexFetch(t *testing.T) {

	//t.SkipNow()

	// open database
	database := &Database{}
	err := database.Open(&db.Configuration{
		Location: getTempFile(t),
	})
	require.Nil(t, err)
	defer cleanUp(t, database)

	// test feeds
	feeds, err := database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 0, len(feeds))

	maxTime := time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(maxTime)
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 0, len(feeds))

	// create new feed, add to database
	feed := m.NewFeed("http://localhost/")
	err = database.SaveFeed(feed)
	assert.Nil(t, err)

	// retest
	feeds, err = database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 1, len(feeds))

	maxTime = time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(maxTime)
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 1, len(feeds))

	// modify feed, resave to database
	// create new feed, add to database
	feed2 := m.NewFeed("https://localhost/")
	feed2.ID = feed.ID
	err = database.SaveFeed(feed2)
	f3 := m.NewFeed("http://kangaroo.com/")
	err = database.SaveFeed(f3)
	assert.Nil(t, err)

	// retest
	feeds, err = database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 2, len(feeds))

	maxTime = time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(maxTime)
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 2, len(feeds))

}
