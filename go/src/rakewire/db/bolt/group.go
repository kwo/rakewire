package bolt

import (
	"bytes"
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"strconv"
)

// GroupGet retrieves a single group.
func (z *Service) GroupGet(groupID uint64) (*m.Group, error) {
	var group *m.Group
	z.Lock()
	defer z.Unlock()
	err := z.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.GroupEntity))
		if values, ok := kvGet(groupID, b); ok {
			group = &m.Group{}
			if err := group.Deserialize(values); err != nil {
				return err
			}
		}
		return nil
	})
	return group, err
}

// GroupGetAllByUser retrieves the groups belonging to the user.
func (z *Service) GroupGetAllByUser(userID uint64) ([]*m.Group, error) {

	result := []*m.Group{}

	// define index keys
	g := &m.Group{}
	g.UserID = userID
	minKeys := g.IndexKeys()[m.GroupIndexUserGroup]
	g.UserID = userID + 1
	nxtKeys := g.IndexKeys()[m.GroupIndexUserGroup]

	err := z.db.View(func(tx *bolt.Tx) error {

		bIndex := tx.Bucket([]byte(bucketIndex)).Bucket([]byte(m.GroupEntity)).Bucket([]byte(m.GroupIndexUserGroup))
		b := tx.Bucket([]byte(bucketData)).Bucket([]byte(m.GroupEntity))

		c := bIndex.Cursor()
		min := []byte(kvKeys(minKeys))
		nxt := []byte(kvKeys(nxtKeys))
		for k, v := c.Seek(min); k != nil && bytes.Compare(k, nxt) < 0; k, v = c.Next() {

			id, err := strconv.ParseUint(string(v), 10, 64)
			if err != nil {
				return err
			}

			if data, ok := kvGet(id, b); ok {
				g := &m.Group{}
				if err := g.Deserialize(data); err != nil {
					return err
				}
				result = append(result, g)
			}

		}

		return nil

	})

	return result, err

}

// GroupSave saves a group to the database.
func (z *Service) GroupSave(group *m.Group) error {

	z.Lock()
	defer z.Unlock()
	err := z.db.Update(func(tx *bolt.Tx) error {

		// new group, check for unique group name
		if group.GetID() == 0 {
			if _, ok := kvGetFromIndex(m.GroupEntity, m.GroupIndexUserGroup, group.IndexKeys()[m.GroupIndexUserGroup], tx); ok {
				return fmt.Errorf("Cannot save group, group name is already taken: %s", group.Name)
			}
		}

		return kvSave(m.GroupEntity, group, tx)

	})

	return err

}

// GroupDelete deletes a group in the database.
func (z *Service) GroupDelete(group *m.Group) error {
	z.Lock()
	defer z.Unlock()
	return z.db.Update(func(tx *bolt.Tx) error {
		return kvDelete(m.GroupEntity, group, tx)
	})
}
