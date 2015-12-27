package model

import (
	"time"
)

//go:generate gokv $GOFILE

// UserEntry defines an entry's status for a specific user.
type UserEntry struct {
	ID      uint64
	UserID  uint64    `kv:"Read:1,Star:1"`
	EntryID uint64    `kv:"Read:4,Star:4"`
	Updated time.Time `kv:"Read:3,Star:3"`
	Read    bool      `kv:"Read:2"`
	Starred bool      `kv:"Star:2"`
	Entry   *Entry    `kv:"-"`
}
