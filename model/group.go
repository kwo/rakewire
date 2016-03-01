package model

//go:generate gokv $GOFILE

//Group defines a group of feeds for a user
type Group struct {
	ID     string `json:"id" kv:"+groupby"`
	UserID string `json:"userID" kv:"+required,UserGroup:1"`
	Name   string `json:"name" kv:"+required,+groupby,UserGroup:2:lower"`
}

// NewGroup creates a new group with the specified user
func NewGroup(userID string, name string) *Group {
	return &Group{
		UserID: userID,
		Name:   name,
	}
}

func (z *Group) setID(fn fnUniqueID) error {
	if z.ID == empty {
		if id, err := fn(); err == nil {
			z.ID = id
		} else {
			return err
		}
	}
	return nil
}
