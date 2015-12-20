package fever

import (
	m "rakewire/model"
	"strconv"
	"strings"
)

func (z *API) getGroupsAndFeedGroups(userID uint64) ([]*Group, []*FeedGroup, error) {

	groups := []*Group{}
	feedGroups := []*FeedGroup{}

	groupFeeds, err := z.getGroupFeeds(userID)
	if err != nil {
		return nil, nil, err
	}

	for group, userfeeds := range groupFeeds {

		g := &Group{
			ID:    group.ID,
			Title: group.Name,
		}
		groups = append(groups, g)

		var feedIDs []string
		for _, userfeed := range userfeeds {
			feedID := strconv.Itoa(int(userfeed.ID))
			feedIDs = append(feedIDs, feedID)
		}

		fg := &FeedGroup{
			GroupID: group.ID,
			FeedIDs: strings.Join(feedIDs, ","),
		}
		feedGroups = append(feedGroups, fg)

	}

	return groups, feedGroups, nil

}

func (z *API) getGroupFeeds(userID uint64) (map[*m.Group][]*m.UserFeed, error) {

	result := make(map[*m.Group][]*m.UserFeed)

	groups, err := z.db.GroupGetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	userfeeds, err := z.db.UserFeedGetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		ufeeds := []*m.UserFeed{}
		result[group] = ufeeds
		for _, userfeed := range userfeeds {
			for _, groupID := range userfeed.GroupIDs {
				if groupID == group.ID {
					ufeeds = append(ufeeds, userfeed)
				}
			}
		}
	}

	return result, nil

}
