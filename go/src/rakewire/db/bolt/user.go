package bolt

import (
	"github.com/boltdb/bolt"
	m "rakewire/model"
)

// UserGetByUsername get a user object by username, nil if not found
func (z *Service) UserGetByUsername(username string) (*m.User, error) {

	user := &m.User{}

	err := z.db.View(func(tx *bolt.Tx) error {
		data, err := kvGetIndex(user.GetName(), "Username", []string{username}, tx)
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
		data, err := kvGetIndex(user.GetName(), "FeverHash", []string{feverhash}, tx)
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
		return kvSave(user, tx)
	})

	return err

}
