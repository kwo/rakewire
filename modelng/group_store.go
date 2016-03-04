package modelng

import (
	"bytes"
)

// G groups all group database methods
var G = &groupStore{}

type groupStore struct{}

func (z *groupStore) Delete(id string, tx Transaction) error {
	return delete(entityGroup, id, tx)
}

func (z *groupStore) Get(id string, tx Transaction) *Group {
	bData := tx.Bucket(bucketData, entityGroup)
	if data := bData.Get([]byte(id)); data != nil {
		group := &Group{}
		if err := group.decode(data); err == nil {
			return group
		}
	}
	return nil
}

func (z *groupStore) GetForUser(userID string, tx Transaction) Groups {
	// index Group UserName = UserID|Name : GroupID
	groups := Groups{}
	min, max := keyMinMax(userID)
	b := tx.Bucket(bucketIndex, entityGroup, indexGroupUserName)
	c := b.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		groupID := string(v)
		if group := z.Get(groupID, tx); group != nil {
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

func (z *groupStore) Save(group *Group, tx Transaction) error {
	return save(entityGroup, group, tx)
}
