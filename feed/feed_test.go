package feed

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestAtomFeed(t *testing.T) {

	f, err1 := os.Open("../test/feed.xml")
	assert.Nil(t, err1)
	assert.NotNil(t, f)
	defer f.Close()

	var feed, err2 = Parse(f)
	require.Nil(t, err2)
	require.NotNil(t, feed)

	assert.NotEmpty(t, feed.ID)
	assert.NotEmpty(t, feed.Title)
	assert.Empty(t, feed.Subtitle)

	assert.NotNil(t, feed.Updated)
	assert.Equal(t, time.Date(2013, time.May, 31, 13, 54, 0, 0, time.UTC), *feed.Updated)

	assert.NotNil(t, feed.Author)

	assert.Empty(t, feed.Icon)
	assert.Empty(t, feed.Generator)

	assert.NotEmpty(t, feed.Links)
	assert.Equal(t, 1, len(feed.Links))
	assert.Equal(t, "https://ostendorf.com/feed.xml", feed.Links["self"])

	assert.NotNil(t, feed.Entries)
	assert.Equal(t, 6, len(feed.Entries))

	assert.Nil(t, feed.Entries[0].Created)
	assert.NotEmpty(t, feed.Entries[0].Summary)
	assert.NotEmpty(t, feed.Entries[0].Content)

}

func TestRSSFeed(t *testing.T) {

	f, err1 := os.Open("../test/wordpress.xml")
	assert.Nil(t, err1)
	assert.NotNil(t, f)
	defer f.Close()

	var feed, err2 = Parse(f)
	require.NoError(t, err2)
	require.NotNil(t, feed)

	//assert.NotEmpty(t, feed.ID)
	assert.NotEmpty(t, feed.Title)
	assert.NotEmpty(t, feed.Subtitle)

}
