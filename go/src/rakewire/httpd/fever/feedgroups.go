package fever

import (
	m "rakewire/model"
	"strconv"
	"strings"
)

func (z *API) getFeeds(userID uint64) ([]*Feed, []*FeedGroup, error) {

	mGroups, err := z.db.GroupGetAllByUser(userID)
	if err != nil {
		return nil, nil, err
	}

	mFeeds, err := z.db.UserFeedGetAllByUser(userID)
	if err != nil {
		return nil, nil, err
	}

	feeds := []*Feed{}
	for _, mFeed := range mFeeds {
		feed := &Feed{
			ID:          mFeed.ID,
			Title:       mFeed.Title,
			FaviconID:   0, // TODO: favicon ID
			URL:         mFeed.Feed.URL,
			SiteURL:     mFeed.Feed.SiteURL,
			IsSpark:     0,
			LastUpdated: mFeed.Feed.LastUpdated.Unix(),
		}
		feeds = append(feeds, feed)
	}

	feedGroups := makeFeedGroups(mGroups, mFeeds)

	return feeds, feedGroups, nil

}

func (z *API) getGroups(userID uint64) ([]*Group, []*FeedGroup, error) {

	mGroups, err := z.db.GroupGetAllByUser(userID)
	if err != nil {
		return nil, nil, err
	}

	mFeeds, err := z.db.UserFeedGetAllByUser(userID)
	if err != nil {
		return nil, nil, err
	}

	groups := []*Group{}
	for _, mGroup := range mGroups {
		group := &Group{
			ID:    mGroup.ID,
			Title: mGroup.Name,
		}
		groups = append(groups, group)
	}

	feedGroups := makeFeedGroups(mGroups, mFeeds)

	return groups, feedGroups, nil

}

func makeFeedGroups(mGroups []*m.Group, mFeeds []*m.UserFeed) []*FeedGroup {

	contains := func(i uint64, a []uint64) bool {
		for _, x := range a {
			if x == i {
				return true
			}
		}
		return false
	}

	feedGroups := []*FeedGroup{}
	for _, mGroup := range mGroups {
		feedIDs := []string{}
		for _, mFeed := range mFeeds {
			if contains(mGroup.ID, mFeed.GroupIDs) {
				feedID := strconv.Itoa(int(mFeed.ID))
				feedIDs = append(feedIDs, feedID)
			}
		}
		feedGroup := &FeedGroup{
			GroupID: mGroup.ID,
			FeedIDs: strings.Join(feedIDs, ","),
		}
		feedGroups = append(feedGroups, feedGroup)
	}
	return feedGroups
}
