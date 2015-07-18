package feedpump

import (
	"rakewire.com/db"
	"rakewire.com/fetch"
)

func databaseFeedsToFetchRequests(dbfeeds *db.Feeds) []*fetch.Request {
	var feeds []*fetch.Request
	for _, v := range dbfeeds.Values {
		feed := &fetch.Request{
			ID:           v.ID,
			ETag:         v.ETag,
			LastModified: v.LastModified,
			URL:          v.URL,
		}
		feeds = append(feeds, feed)
	}
	return feeds
}
