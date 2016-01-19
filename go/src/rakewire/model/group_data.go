package model

import (
	"bytes"
	"fmt"
	"strconv"
)

// GroupByID retrieves a single group.
func GroupByID(groupID uint64, tx Transaction) (group *Group, err error) {
	b := tx.Bucket(bucketData).Bucket(GroupEntity)
	if values, ok := kvGet(groupID, b); ok {
		group = &Group{}
		err = group.Deserialize(values)
	}
	return
}

// GroupsByUser retrieves the groups belonging to the user.
func GroupsByUser(userID uint64, tx Transaction) ([]*Group, error) {

	result := []*Group{}

	// define index keys
	g := &Group{}
	g.UserID = userID
	minKeys := g.IndexKeys()[GroupIndexUserGroup]
	g.UserID = userID + 1
	nxtKeys := g.IndexKeys()[GroupIndexUserGroup]

	bIndex := tx.Bucket(bucketIndex).Bucket(GroupEntity).Bucket(GroupIndexUserGroup)
	b := tx.Bucket(bucketData).Bucket(GroupEntity)

	c := bIndex.Cursor()
	min := []byte(kvKeys(minKeys))
	nxt := []byte(kvKeys(nxtKeys))
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

		id, err := strconv.ParseUint(string(v), 10, 64)
		if err != nil {
			return nil, err
		}

		if data, ok := kvGet(id, b); ok {
			g := &Group{}
			if err := g.Deserialize(data); err != nil {
				return nil, err
			}
			result = append(result, g)
		}

	}

	return result, nil

}

// Save saves a group to the database.
func (group *Group) Save(tx Transaction) error {

	// new group, check for unique group name
	if group.GetID() == 0 {
		if _, ok := kvGetFromIndex(GroupEntity, GroupIndexUserGroup, group.IndexKeys()[GroupIndexUserGroup], tx); ok {
			return fmt.Errorf("Cannot save group, group name is already taken: %s", group.Name)
		}
	}

	return kvSave(GroupEntity, group, tx)

}

// Delete deletes a group in the database.
func (group *Group) Delete(tx Transaction) error {
	return kvDelete(GroupEntity, group, tx)
}
