package bolt

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire.com/db"
	m "rakewire.com/model"
	"testing"
	"time"
)

func TestFeeds(t *testing.T) {

	//t.SkipNow()

	urls, err := readFile(feedFile)
	require.Nil(t, err)
	require.NotNil(t, urls)

	feeds := m.NewFeeds()
	for _, url := range urls {
		feeds.Add(m.NewFeed(url))
	}

	database := Database{}
	err = database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	err = database.SaveFeeds(feeds)
	assert.Nil(t, err)

	feeds2, err := database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds2)

	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	err = os.Remove(databaseFile)
	assert.Nil(t, err)

	assert.Equal(t, len(urls), feeds.Size())
	assert.Equal(t, feeds.Size(), feeds2.Size())
	// for k, v := range feedmap {
	// 	fmt.Printf("Feed %s: %v\n", k, v.URL)
	// }

}

func TestURLIndex(t *testing.T) {

	const URL1 = "http://localhost/"
	const URL2 = "http://localhost:8888/"

	// open database
	database := Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	// create feeds, feed
	feeds := m.NewFeeds()
	feed := m.NewFeed(URL1)
	feeds.Add(feed)
	assert.Equal(t, 1, feeds.Size())
	assert.Equal(t, URL1, feed.URL)

	// save feeds
	err = database.SaveFeeds(feeds)
	require.Nil(t, err)

	// get feeds2, feed2
	feeds2, err := database.GetFeeds()
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, feeds.Size(), feeds2.Size())
	feed2 := feeds2.Values[0]
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
	feed2 = feeds2.Values[0]
	feed2.URL = URL2
	err = database.SaveFeeds(feeds2)
	require.Nil(t, err)

	// get feeds2, feed2
	feeds2, err = database.GetFeeds()
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, feeds.Size(), feeds2.Size())
	feed2 = feeds2.Values[0]
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

	// close database
	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	// remove file
	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}

func TestNextFetchKeyCompare(t *testing.T) {

	assert.Equal(t, 1, bytes.Compare([]byte("2015-07-21T07:00:24Z#"), []byte("2015-07-21T07:00:24Z!c35d9174-2f74-11e5-baf1-5cf938992b62")))

	f := m.NewFeed("http://localhost/")
	assert.NotNil(t, f.NextFetch)

	key := fetchKey(f)
	max := formatMaxTime(*f.NextFetch)
	assert.Equal(t, 1, bytes.Compare([]byte(max), []byte(key)))

}

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
	feeds = m.NewFeeds()
	feed := m.NewFeed("http://localhost/")
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
	// create new feed, add to database
	feeds2 := m.NewFeeds()
	feed2 := m.NewFeed("https://localhost/")
	feed2.ID = feed.ID
	feeds2.Add(feed2)
	f3 := m.NewFeed("http://kangaroo.com/")
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

	feeds := m.NewFeeds()
	for _, url := range urls {
		feed := m.NewFeed(url)
		feeds.Add(feed)
	}

	err = database.SaveFeeds(feeds)
	assert.Nil(t, err)

	feeds, err = database.GetFeeds()
	assert.Nil(t, err)
	assert.NotNil(t, feeds)
	assert.Equal(t, 288, feeds.Size())

	for _, feed := range feeds.Values {
		err = database.checkIndexForEntries(bucketIndexNextFetch, feed.ID, 1)
		assert.Nil(t, err)
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
		err = database.checkIndexForEntries(bucketIndexNextFetch, feed.ID, 1)
		assert.Nil(t, err)
	}

	// close database
	err = database.Close()
	assert.Nil(t, err)

	// delete test file
	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}
