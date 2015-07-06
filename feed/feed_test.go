package feed

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

const (
	FeedFilename = "../test/feed.xml"
)

func TestFeed(t *testing.T) {

	f, err1 := os.Open(FeedFilename)
	assert.Nil(t, err1)
	assert.NotNil(t, f)
	defer f.Close()

	var feed, err2 = Parse(f)
	assert.Nil(t, err2)
	assert.NotNil(t, feed)

	assert.NotEmpty(t, feed.Title)
	assert.Empty(t, feed.Subtitle)

	assert.NotNil(t, feed.Updated)
	assert.Equal(t, time.Date(2013, time.May, 31, 13, 54, 0, 0, time.UTC), *feed.Updated)

	assert.NotNil(t, feed.Author)

	assert.Empty(t, feed.Icon)
	assert.Empty(t, feed.Generator)

	assert.NotNil(t, feed.Links)
	assert.Equal(t, 1, len(feed.Links))
	assert.Equal(t, "self", feed.Links[0].Rel)
	assert.Equal(t, "https://ostendorf.com/feed.xml", feed.Links[0].Href)

	assert.NotNil(t, feed.Entries)
	assert.Equal(t, 6, len(feed.Entries))

	assert.Nil(t, feed.Entries[0].Created)
	assert.NotEmpty(t, feed.Entries[0].Summary)
	assert.NotEmpty(t, feed.Entries[0].Content)

}
