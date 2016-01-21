package model

import (
	"time"
)

//go:generate gokv $GOFILE

// UserEntry defines an entry's status for a specific user.
type UserEntry struct {
	ID         uint64    `kv:"Read:4,Star:4,User:2"`
	UserID     uint64    `kv:"+required,Read:1,Star:1,User:1"`
	EntryID    uint64    `kv:"+required"`
	UserFeedID uint64    `kv:"+required"`
	Updated    time.Time `kv:"Read:3,Star:3"`
	IsRead     bool      `kv:"Read:2"`
	IsStar     bool      `kv:"Star:2"`
	Entry      *Entry    `kv:"-"`
}
