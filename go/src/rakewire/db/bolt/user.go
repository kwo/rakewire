package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"strings"
)

// UserGetByUsername get a user object by username, nil if not found
func (z *Service) UserGetByUsername(username string) (*m.User, error) {

	found := false
	user := &m.User{}

	err := z.db.View(func(tx *bolt.Tx) error {
		if data, ok := kvGetFromIndex(m.UserEntity, m.UserIndexUsername, []string{username}, tx); ok {
			found = true
			return user.Deserialize(data)
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else if !found {
		return nil, nil
	}
	return user, nil

}

// UserGetByFeverHash get a user object by feverhash, nil if not found
func (z *Service) UserGetByFeverHash(feverhash string) (*m.User, error) {

	found := false
	user := &m.User{}

	err := z.db.View(func(tx *bolt.Tx) error {
		if data, ok := kvGetFromIndex(m.UserEntity, m.UserIndexFeverHash, []string{feverhash}, tx); ok {
			found = true
			return user.Deserialize(data)
		}
		return nil
	})

	if err != nil {
		return nil, err
	} else if !found {
		return nil, nil
	}
	return user, nil

}

// UserSave saves a user to the database.
func (z *Service) UserSave(user *m.User) error {

	z.Lock()
	defer z.Unlock()
	err := z.db.Update(func(tx *bolt.Tx) error {

		// new user, check for unique username
		if user.GetID() == 0 {
			indexName := m.UserIndexUsername
			if _, ok := kvGetFromIndex(m.UserEntity, indexName, user.IndexKeys()[indexName], tx); ok {
				return fmt.Errorf("Cannot save user, username is already taken: %s", strings.ToLower(user.Username))
			}
		}

		return kvSave(m.UserEntity, user, tx)

	})

	return err

}
