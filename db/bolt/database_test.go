package bolt

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"rakewire.com/db"
	"strings"
	"testing"
	"time"
)

const (
	feedFile     = "../../test/feedlist.txt"
	databaseFile = "../../test/bolt.db"
)

func TestInterface(t *testing.T) {

	var d db.Database = &Database{}
	assert.NotNil(t, d)

}

func TestFeeds(t *testing.T) {

	urls, err := readFile(feedFile)
	require.Nil(t, err)
	require.NotNil(t, urls)

	feeds := db.NewFeeds()
	for _, url := range urls {
		feeds.Add(db.NewFeed(url))
	}
	assert.Equal(t, len(urls), feeds.Size())

	database := Database{}
	err = database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	err = database.SaveFeeds(feeds)
	require.Nil(t, err)

	feeds2, err := database.GetFeeds()
	require.Nil(t, err)
	require.NotNil(t, feeds2)

	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	err = os.Remove(databaseFile)
	assert.Nil(t, err)

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
	feeds := db.NewFeeds()
	feed := db.NewFeed(URL1)
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

	assert.Equal(t, -1, bytes.Compare([]byte("2015-07-21T06:59:24Z#c35d9174-2f74-11e5-baf1-5cf938992b62"), []byte("2015-07-21T07:00:24Z")))

	now := time.Now()
	f := db.NewFeed("http://localhost/")
	f.LastFetch = &now
	f.Frequency = 5

	key := fetchKey(f)
	max := formatFetchTime(now.Add(5 * time.Minute))
	assert.Equal(t, 1, bytes.Compare([]byte(key), []byte(max)))
	max = formatFetchTime(now.Add(5 * time.Minute).Add(1 * time.Second))
	assert.Equal(t, -1, bytes.Compare([]byte(key), []byte(max)))

}

func TestNextFetch(t *testing.T) {

	// open database
	database := Database{}
	err := database.Open(&db.Configuration{
		Location: databaseFile,
	})
	require.Nil(t, err)

	// create feeds, feed
	feeds := db.NewFeeds()
	now := time.Now()
	lastFetch := &now
	for i := 0; i < 10; i++ {
		url := fmt.Sprintf("http://localhost:888%d", i)
		feed := db.NewFeed(url)
		feed.LastFetch = lastFetch
		feed.Frequency = 5 + i
		feeds.Add(feed)
	}

	assert.Equal(t, 10, feeds.Size())

	// save feeds
	err = database.SaveFeeds(feeds)
	require.Nil(t, err)

	// right now there should be no feeds up for fetch
	feeds2, err := database.GetFetchFeeds(nil)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 0, feeds2.Size())

	maxTime := now.Add(1 * time.Minute)
	feeds2, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 0, feeds2.Size())

	maxTime = now.Add(5 * time.Minute)
	feeds2, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 1, feeds2.Size())

	maxTime = now.Add(10 * time.Minute)
	feeds2, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 6, feeds2.Size())

	maxTime = now.Add(14 * time.Minute)
	feeds2, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 10, feeds2.Size())

	// again

	now = time.Now()
	for _, f := range feeds.Values {
		nt := now.Add(-90 * time.Second)
		f.LastFetch = &nt
	}
	// save feeds
	err = database.SaveFeeds(feeds)
	require.Nil(t, err)

	// right now there should be no feeds up for fetch
	feeds2, err = database.GetFetchFeeds(nil)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 0, feeds2.Size())

	maxTime = now.Add(1 * time.Minute)
	feeds2, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 0, feeds2.Size())

	maxTime = now.Add(5 * time.Minute)
	feeds2, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 2, feeds2.Size())

	maxTime = now.Add(10 * time.Minute)
	feeds2, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 7, feeds2.Size())

	maxTime = now.Add(20 * time.Minute)
	feeds2, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 10, feeds2.Size())

	// logger.Printf("max: %s\n", formatMaxTime(maxTime))
	// for _, f := range feeds2.Values {
	// 	logger.Printf("%s: %d %s\n", f.URL, f.Frequency, formatFetchTime(*f.GetNextFetchTime()))
	// }

	feeds2, err = database.GetFeeds()
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 10, feeds2.Size())

	maxTime = now.Add(48 * time.Hour)
	feeds2, err = database.GetFetchFeeds(&maxTime)
	require.Nil(t, err)
	require.NotNil(t, feeds2)
	assert.Equal(t, 10, feeds2.Size())

	// close database
	err = database.Close()
	assert.Nil(t, err)
	assert.Nil(t, database.db)

	// remove file
	err = os.Remove(databaseFile)
	assert.Nil(t, err)

}

func readFile(feedfile string) ([]string, error) {

	var result []string

	f, err1 := os.Open(feedfile)
	if err1 != nil {
		return nil, err1
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var url = strings.TrimSpace(scanner.Text())
		if url != "" && url[:1] != "#" {
			result = append(result, url)
		}
	}
	f.Close()

	return result, nil

}
