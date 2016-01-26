package model

import (
	"bytes"
	"fmt"
	"strconv"
)

// GroupByID retrieves a single group.
func GroupByID(groupID uint64, tx Transaction) (group *Group, err error) {
	b := tx.Bucket(bucketData).Bucket(groupEntity)
	if values, ok := kvGet(groupID, b); ok {
		group = &Group{}
		err = group.deserialize(values)
	}
	return
}

// GroupsByUser retrieves the groups belonging to the user.
func GroupsByUser(userID uint64, tx Transaction) (Groups, error) {

	result := Groups{}

	// define index keys
	g := &Group{}
	g.UserID = userID
	minKeys := g.indexKeys()[groupIndexUserGroup]
	g.UserID = userID + 1
	nxtKeys := g.indexKeys()[groupIndexUserGroup]

	bIndex := tx.Bucket(bucketIndex).Bucket(groupEntity).Bucket(groupIndexUserGroup)
	b := tx.Bucket(bucketData).Bucket(groupEntity)

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
			if err := g.deserialize(data); err != nil {
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
	if group.getID() == 0 {
		if _, ok := kvGetFromIndex(groupEntity, groupIndexUserGroup, group.indexKeys()[groupIndexUserGroup], tx); ok {
			return fmt.Errorf("Cannot save group, group name is already taken: %s", group.Name)
		}
	}

	return kvSave(groupEntity, group, tx)

}

// Delete deletes a group in the database.
func (group *Group) Delete(tx Transaction) error {
	return kvDelete(groupEntity, group, tx)
}
