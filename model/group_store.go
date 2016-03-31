package model

import (
	"bytes"
)

// G groups all group database methods
var G = &groupStore{}

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

func (z *groupStore) Range(tx Transaction) Groups {
	groups := Groups{}
	c := tx.Bucket(bucketData, entityGroup).Cursor()
	for k, v := c.First(); k != nil; k, v = c.Next() {
		group := &Group{}
		if err := group.decode(v); err == nil {
			groups = append(groups, group)
		}
	}
	return groups
}

func (z *groupStore) Save(tx Transaction, group *Group) error {
	return saveObject(tx, entityGroup, group)
}
