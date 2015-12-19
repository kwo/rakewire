package model

//go:generate gokv $GOFILE

// UserEntry defines an entry's status for a specific user.
type UserEntry struct {
	ID      uint64 `kv:"User:2,Starred:3"`
	UserID  uint64 `kv:"User:1,Starred:1"`
	Read    bool
	Starred bool `kv:"Starred:2"`
}
