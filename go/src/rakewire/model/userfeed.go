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
		UserID: userID,
		FeedID: feedID,
	}
}
