package model

import (
	"bytes"
	"fmt"
	"strings"
)

// GroupByID retrieves a single group.
func GroupByID(groupID string, tx Transaction) (group *Group, err error) {
	b := tx.Bucket(bucketData, groupEntity)
	if record := b.GetRecord(groupID); record != nil {
		group = &Group{}
		err = group.deserialize(record)
	}
	return
}

// GroupsByUser retrieves the groups belonging to the user.
func GroupsByUser(userID string, tx Transaction) (Groups, error) {

	result := Groups{}

	// group index UserGroup = UserID|Name : GroupID
	min, max := kvKeyMinMax2(userID)
	bIndex := tx.Bucket(bucketIndex).Bucket(groupEntity).Bucket(groupIndexUserGroup)
	b := tx.Bucket(bucketData).Bucket(groupEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		groupID := string(v)
		if record := b.GetRecord(groupID); record != nil {
			g := &Group{}
			if err := g.deserialize(record); err != nil {
				return nil, err
			}
			result = append(result, g)
		}
	}

	return result, nil

}

// GroupByName retrieve a group by user and name.
func GroupByName(userID, groupName string, tx Transaction) (group *Group, err error) {

	bGroup := tx.Bucket(bucketData, groupEntity)
	bIndex := tx.Bucket(bucketIndex, groupEntity, groupIndexUserGroup)

	key := kvKeyEncode(userID, strings.ToLower(groupName))

	if record := bIndex.GetIndex(bGroup, key); record != nil {
		group = &Group{}
		err = group.deserialize(record)
	}

	return

}

// Save saves a group to the database.
func (z *Group) Save(tx Transaction) error {

	// new group, check for unique group name
	if z.getID() == empty {
		group, err := GroupByName(z.UserID, z.Name, tx)
		if err != nil {
			return err
		} else if group != nil {
			return fmt.Errorf("Cannot save group, name is already taken: %s", strings.ToLower(z.Name))
		}
	}

	return kvSave(groupEntity, z, tx)

}

// Delete deletes a group in the database.
func (z *Group) Delete(tx Transaction) error {
	return kvDelete(groupEntity, z, tx)
}
