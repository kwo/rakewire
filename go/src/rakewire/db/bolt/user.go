package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
	"strings"
)

// UserGetByUsername get a user object by username, nil if not found
func (z *Service) UserGetByUsername(username string) (*m.User, error) {

	user := &m.User{}

	err := z.db.View(func(tx *bolt.Tx) error {
		data, err := kvGetIndex(user.GetName(), m.IndexUserUsername, []string{username}, tx)
		if err != nil {
			return err
		}
		return user.Deserialize(data)
	})

	return user, err

}

// UserGetByFeverHash get a user object by feverhash, nil if not found
func (z *Service) UserGetByFeverHash(feverhash string) (*m.User, error) {

	user := &m.User{}

	err := z.db.View(func(tx *bolt.Tx) error {
		data, err := kvGetIndex(user.GetName(), m.IndexUserFeverHash, []string{feverhash}, tx)
		if err != nil {
			return err
		}
		return user.Deserialize(data)
	})

	return user, err

}

// UserSave saves a user to the database.
func (z *Service) UserSave(user *m.User) error {

	z.Lock()
	defer z.Unlock()
	err := z.db.Update(func(tx *bolt.Tx) error {

		// new user, check for unique username
		if user.GetID() == 0 {
			indexName := m.IndexUserUsername
			data, err := kvGetIndex(user.GetName(), indexName, user.IndexKeys()[indexName], tx)
			if err != nil {
				return err
			}
			if data != nil {
				return fmt.Errorf("Cannot save user, username is already taken: %s", strings.ToLower(user.Username))
			}
		}

		return kvSave(user, tx)

	})

	return err

}
