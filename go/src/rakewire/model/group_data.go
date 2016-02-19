package model

import (
	"bytes"
)

// GroupByID retrieves a single group.
func GroupByID(groupID string, tx Transaction) (group *Group, err error) {
	b := tx.Bucket(bucketData).Bucket(groupEntity)
	if values, ok := kvGet(groupID, b); ok {
		group = &Group{}
		err = group.deserialize(values)
	}
	return
}

// GroupsByUser retrieves the groups belonging to the user.
func GroupsByUser(userID string, tx Transaction) (Groups, error) {

	result := Groups{}

	// group index UserGroup = UserID|Name : GroupID
	min, max := kvKeyMinMax(userID)
	bIndex := tx.Bucket(bucketIndex).Bucket(groupEntity).Bucket(groupIndexUserGroup)
	b := tx.Bucket(bucketData).Bucket(groupEntity)

	c := bIndex.Cursor()
	for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
		groupID := string(v)
		if data, ok := kvGet(groupID, b); ok {
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

	// TODO: new group, check for unique group name
	// if group.getID() ==  {
	// 	if _, ok := kvGetFromIndex(groupEntity, groupIndexUserGroup, group.indexKeys()[groupIndexUserGroup], tx); ok {
	// 		return fmt.Errorf("Cannot save group, group name is already taken: %s", group.Name)
	// 	}
	// }

	return kvSave(groupEntity, group, tx)

}

// Delete deletes a group in the database.
func (group *Group) Delete(tx Transaction) error {
	return kvDelete(groupEntity, group, tx)
}
