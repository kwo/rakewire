package feed

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestAtomFeed(t *testing.T) {

	f, err := os.Open("../test/feed.xml")
	assert.Nil(t, err)
	assert.NotNil(t, f)
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	assert.Nil(t, err)

	feed, err := Parse(body)
	require.Nil(t, err)
	require.NotNil(t, feed)

	assert.NotEmpty(t, feed.ID)
	assert.NotEmpty(t, feed.Title)
	assert.Empty(t, feed.Subtitle)

	assert.NotNil(t, feed.Updated)
	assert.Equal(t, time.Date(2013, time.May, 31, 13, 54, 0, 0, time.UTC), *feed.Updated)

	//assert.NotNil(t, feed.Author)

	assert.Empty(t, feed.Icon)
	assert.Empty(t, feed.Generator)

	assert.NotEmpty(t, feed.Links)
	assert.Equal(t, 1, len(feed.Links))
	assert.Equal(t, "https://ostendorf.com/feed.xml", feed.Links["self"])

	assert.NotNil(t, feed.Entries)
	assert.Equal(t, 6, len(feed.Entries))

	assert.True(t, feed.Entries[0].Created.IsZero())
	assert.NotEmpty(t, feed.Entries[0].Summary)
	assert.NotEmpty(t, feed.Entries[0].Content)

	feedFmt := "%-12s %s"
	t.Logf(feedFmt, "ID", feed.ID)
	t.Logf(feedFmt, "Title", feed.Title)
	t.Logf(feedFmt, "Flavor", feed.Flavor)
	t.Logf(feedFmt, "Updated", feed.Updated.Format("2006-01-02 15:04:05"))
	t.Logf(feedFmt, "Generator", feed.Generator)
	for _, e := range feed.Entries {
		t.Logf("%s %s %s", e.ID, e.Updated.Format("2006-01-02 15:04:05"), e.Title)
	}

}

func TestRSSFeed(t *testing.T) {

	f, err := os.Open("../test/wordpress.xml")
	assert.Nil(t, err)
	assert.NotNil(t, f)
	defer f.Close()
	body, err := ioutil.ReadAll(f)
	assert.Nil(t, err)

	feed, err := Parse(body)
	require.Nil(t, err)
	require.NotNil(t, feed)

	//assert.NotEmpty(t, feed.ID)
	assert.NotEmpty(t, feed.Title)
	assert.NotEmpty(t, feed.Subtitle)

	feedFmt := "%-12s %s"
	t.Logf(feedFmt, "ID", feed.ID)
	t.Logf(feedFmt, "Title", feed.Title)
	t.Logf(feedFmt, "Flavor", feed.Flavor)
	t.Logf(feedFmt, "Updated", feed.Updated)
	t.Logf(feedFmt, "Generator", feed.Generator)
	for _, e := range feed.Entries {
		t.Logf("%s %s %s", e.ID, e.Updated.Format("2006-01-02 15:04:05"), e.Title)
	}

}
