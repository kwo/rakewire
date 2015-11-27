package bolt

import (
	m "rakewire/model"
	"testing"
	"time"
)

func TestFeeds(t *testing.T) {

	//t.SkipNow()
	database := openDatabase(t)
	defer closeDatabase(t, database)
	assertNotNil(t, database)

	feeds, err := m.ParseFeedsFromFile(feedFile)
	assertNoError(t, err)
	assertNotNil(t, feeds)

	for _, feed := range feeds {
		// t.Logf("URL (%d): %s", n, feed.URL)
		feed.Attempt = m.NewFeedLog(feed.ID)
	}

	for _, feed := range feeds {
		err = database.SaveFeed(feed)
		assertNoError(t, err)
	}

	feeds2, err := database.GetFeeds()
	assertNoError(t, err)
	assertNotNil(t, feeds2)

	assertEqual(t, len(feeds), len(feeds2))
	// for n, feed := range feeds2 {
	// t.Logf("Feed (%d): %s", n, feed.URL)
	// }

}

func TestURLIndex(t *testing.T) {

	//t.SkipNow()

	const URL1 = "http://localhost/"
	const URL2 = "http://localhost:8888/"

	database := openDatabase(t)
	defer closeDatabase(t, database)
	assertNotNil(t, database)

	// create feed
	feed := m.NewFeed(URL1)
	feed.Attempt = m.NewFeedLog(feed.ID)
	assertEqual(t, URL1, feed.URL)
	assertEqual(t, feed.ID, feed.Attempt.FeedID)

	// save feeds
	err := database.SaveFeed(feed)
	assertNoError(t, err)

	var feed2 *m.Feed
	var feeds2 []*m.Feed

	// get feed2
	feeds2, err = database.GetFeeds()
	assertNoError(t, err)
	assertNotNil(t, feeds2)
	assertEqual(t, 1, len(feeds2))
	feed2 = feeds2[0]
	assertNotNil(t, feed2)
	assertEqual(t, feed.ID, feed2.ID)
	assertEqual(t, URL1, feed2.URL)

	// get by URL
	feed2, err = database.GetFeedByURL(URL1)
	assertNoError(t, err)
	assertNotNil(t, feed2)
	assertEqual(t, feed.ID, feed2.ID)
	assertEqual(t, URL1, feed2.URL)

	// update URL
	feed2 = feeds2[0]
	feed2.URL = URL2
	err = database.SaveFeed(feed2)
	assertNoError(t, err)

	// get feeds2, feed2
	feeds2, err = database.GetFeeds()
	assertNoError(t, err)
	assertNotNil(t, feeds2)
	assertEqual(t, 1, len(feeds2))
	feed2 = feeds2[0]
	assertNotNil(t, feed2)
	assertEqual(t, feed.ID, feed2.ID)
	assertEqual(t, URL2, feed2.URL)

	// get by old URL
	feed2, err = database.GetFeedByURL(URL1)
	assertNoError(t, err)
	assertNil(t, feed2)

	// get by new URL
	feed2, err = database.GetFeedByURL(URL2)
	assertNoError(t, err)
	assertNotNil(t, feed2)
	assertEqual(t, feed.ID, feed2.ID)
	assertEqual(t, URL2, feed2.URL)

}

func TestIndexFetch(t *testing.T) {

	//t.SkipNow()

	database := openDatabase(t)
	defer closeDatabase(t, database)
	assertNotNil(t, database)

	// test feeds
	feeds, err := database.GetFeeds()
	assertNoError(t, err)
	assertNotNil(t, feeds)
	assertEqual(t, 0, len(feeds))

	maxTime := time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(maxTime)
	assertNoError(t, err)
	assertNotNil(t, feeds)
	assertEqual(t, 0, len(feeds))

	// create new feed, add to database
	feed := m.NewFeed("http://localhost/")
	err = database.SaveFeed(feed)
	assertNoError(t, err)

	// retest
	feeds, err = database.GetFeeds()
	assertNoError(t, err)
	assertNotNil(t, feeds)
	assertEqual(t, 1, len(feeds))

	maxTime = time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(maxTime)
	assertNoError(t, err)
	assertNotNil(t, feeds)
	assertEqual(t, 1, len(feeds))

	// modify feed, resave to database
	// create new feed, add to database
	feed2 := m.NewFeed("https://localhost/")
	feed2.ID = feed.ID
	err = database.SaveFeed(feed2)
	f3 := m.NewFeed("http://kangaroo.com/")
	err = database.SaveFeed(f3)
	assertNoError(t, err)

	// retest
	feeds, err = database.GetFeeds()
	assertNoError(t, err)
	assertNotNil(t, feeds)
	assertEqual(t, 2, len(feeds))

	maxTime = time.Now().Add(48 * time.Hour)
	feeds, err = database.GetFetchFeeds(maxTime)
	assertNoError(t, err)
	assertNotNil(t, feeds)
	assertEqual(t, 2, len(feeds))

}
