package model

//go:generate gokv $GOFILE

// UserFeed defines a feed specific to a user.
type UserFeed struct {
	ID       uint64
	UserID   uint64 `kv:"User:1"`
	FeedID   uint64 `kv:"User:2"`
	GroupIDs []uint64
	Title    string
	Notes    string
}
