package fever

import (
	"rakewire/model"
	"strings"
)

func (z *API) getFeeds(userID string, tx model.Transaction) ([]*Feed, []*FeedGroup, error) {

	mGroups := model.G.GetForUser(tx, userID)
	mSubscriptions := model.S.GetForUser(tx, userID)
	mFeedsByID := model.F.GetBySubscriptions(tx, mSubscriptions).ByID()

	feeds := []*Feed{}
	for _, mSubscription := range mSubscriptions {
		feed := &Feed{
			ID:          parseID(mSubscription.FeedID),
			Title:       mSubscription.Title,
			FaviconID:   0,
			URL:         mFeedsByID[mSubscription.FeedID].URL,
			SiteURL:     mFeedsByID[mSubscription.FeedID].SiteURL,
			IsSpark:     0,
			LastUpdated: mFeedsByID[mSubscription.FeedID].LastUpdated.Unix(),
		}
		feeds = append(feeds, feed)
	}

	feedGroups := makeFeedGroups(mGroups, mSubscriptions)

	return feeds, feedGroups, nil

}

func (z *API) getGroups(userID string, tx model.Transaction) ([]*Group, []*FeedGroup, error) {

	mGroups := model.G.GetForUser(tx, userID)
	mSubscriptions := model.S.GetForUser(tx, userID)

	groups := []*Group{}
	for _, mGroup := range mGroups {
		group := &Group{
			ID:    parseID(mGroup.ID),
			Title: mGroup.Name,
		}
		groups = append(groups, group)
	}

	feedGroups := makeFeedGroups(mGroups, mSubscriptions)

	return groups, feedGroups, nil

}

func makeFeedGroups(mGroups model.Groups, mSubscriptions model.Subscriptions) []*FeedGroup {

	contains := func(i string, a []string) bool {
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
		for _, mSubscription := range mSubscriptions {
			if contains(mGroup.ID, mSubscription.GroupIDs) {
				feedIDs = append(feedIDs, decodeID(mSubscription.FeedID))
			}
		}
		feedGroup := &FeedGroup{
			GroupID: parseID(mGroup.ID),
			FeedIDs: strings.Join(feedIDs, ","),
		}
		feedGroups = append(feedGroups, feedGroup)
	}
	return feedGroups
}
