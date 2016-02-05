package model

//go:generate gokv $GOFILE

//Group defines a group of feeds for a user
type Group struct {
	ID     uint64 `kv:"+groupby"`
	UserID uint64 `kv:"+required,UserGroup:1"`
	Name   string `kv:"+required,+groupby,UserGroup:2"`
}

// NewGroup creates a new group with the specified user
func NewGroup(userID uint64, name string) *Group {
	return &Group{
		UserID: userID,
		Name:   name,
	}
}

func (z *Group) setIDIfNecessary(fn fnNextID, tx Transaction) error {
	if z.ID == 0 {
		if id, _, err := fn(groupEntity, tx); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}
