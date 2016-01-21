package model

//go:generate gokv $GOFILE

//Group defines a group of feeds for a user
type Group struct {
	ID     uint64
	UserID uint64 `kv:"+required,UserGroup:1"`
	Name   string `kv:"+required,UserGroup:2"`
}

// NewGroup creates a new group with the specified user
func NewGroup(userID uint64, name string) *Group {
	return &Group{
		UserID: userID,
		Name:   name,
	}
}
