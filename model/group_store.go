package model

import (
	"bytes"
	"errors"
)

var (
	// ErrGroupnameTaken occurs when adding a new group with a non-unique group name (per user).
	ErrGroupnameTaken = errors.New("Group name exists already.")
	// G groups all group database methods
	G = &groupStore{}
)

type groupStore struct{}

func (z *groupStore) Delete(tx Transaction, id string) error {
	return deleteObject(tx, entityGroup, id)
}

func (z *groupStore) Get(tx Transaction, id string) *Group {
	bData := tx.Bucket(bucketData, entityGroup)
	if data := bData.Get([]byte(id)); data != nil {
		group := &Group{}
		if err := group.decode(data); err == nil {
			return group
		}
	}
	return nil
}

func (z *groupStore) GetForUser(tx Transaction, userID string) Groups {
	// index Group UserName = UserID|Name : GroupID
	groups := Groups{}
	min, max := keyMinMax(userID)
	b := tx.Bucket(bucketIndex, entityGroup, indexGroupUserName)
	c := b.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		groupID := string(v)
		if group := z.Get(tx, groupID); group != nil {
			groups = append(groups, group)
		}
	}
	return groups
}

func (z *groupStore) New(userID, name string) *Group {
	return &Group{
		UserID: userID,
		Name:   name,
	}
}

func (z *groupStore) Save(tx Transaction, group *Group) error {
	if group.GetID() == empty {
		groupsByName := z.GetForUser(tx, group.UserID).ByName()
		if group := groupsByName[group.Name]; group != nil {
			return ErrGroupnameTaken
		}
	}
	return saveObject(tx, entityGroup, group)
}
