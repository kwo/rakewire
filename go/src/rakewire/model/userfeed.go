package model

//go:generate gokv $GOFILE

// UserFeed defines a feed specific to a user.
type UserFeed struct {
	ID       uint64
	UserID   uint64 `kv:"Feed:2,User:1"`
	FeedID   uint64 `kv:"Feed:1,User:2"`
	GroupIDs []uint64
	Title    string
	Notes    string
	Feed     *Feed `kv:"-"`
}

// NewUserFeed associates a feed with a user.
func NewUserFeed(userID, feedID uint64) *UserFeed {
	return &UserFeed{
		UserID:   userID,
		FeedID:   feedID,
		GroupIDs: []uint64{},
	}
}

// AddGroup adds the userfeed to the given group.
func (z *UserFeed) AddGroup(groupID uint64) {
	if !z.HasGroup(groupID) {
		z.GroupIDs = append(z.GroupIDs, groupID)
	}
}

// RemoveGroup removes the UserFeed from the given group.
func (z *UserFeed) RemoveGroup(groupID uint64) {
	for i, value := range z.GroupIDs {
		if value == groupID {
			z.GroupIDs = append(z.GroupIDs[:i], z.GroupIDs[i+1:]...)
		}
	}
}

// HasGroup tests if the UserFeed belongs to the given group
func (z *UserFeed) HasGroup(groupID uint64) bool {
	result := false
	if len(z.GroupIDs) > 0 {
		for _, value := range z.GroupIDs {
			if value == groupID {
				return true
			}
		}
	}
	return result
}
