package feed

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFeed(t *testing.T) {
	var feed, err = Parse("https://ostendorf.com/feed.xml")
	assert.Nil(t, err)
	assert.NotNil(t, feed)
	assert.NotNil(t, feed.Date)
	assert.Equal(t, time.Date(2013, time.May, 31, 13, 54, 0, 0, time.UTC), *feed.Date)
	assert.NotNil(t, feed.Author)
	assert.NotNil(t, feed.Entries)
	assert.Equal(t, 6, len(feed.Entries))
}
