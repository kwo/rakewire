package fever

import (
	m "rakewire/model"
)

func (z *API) getGroups(userID uint64) ([]*Group, error) {

	result := []*Group{}

	groups, err := z.db.GroupGetAllByUser(userID)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		g := &Group{
			ID:    group.ID,
			Title: group.Name,
		}
		result = append(result, g)
	}

	return result, nil

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
