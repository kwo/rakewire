package bolt

import (
	"fmt"
	"github.com/boltdb/bolt"
	m "rakewire/model"
)

// UserGetByUsername get a user object by username, nil if not found
func (z *Service) UserGetByUsername(username string) (*m.User, error) {

	users := []*m.User{}
	add := func() interface{} {
		u := &m.User{}
		users = append(users, u)
		return u
	}

	err := z.db.View(func(tx *bolt.Tx) error {
		return Query(bucketUser, "Username", []interface{}{username}, []interface{}{username}, add, tx)
	})

	if err != nil {
		return nil, err
	}

	if len(users) == 1 {
		return users[0], nil
	} else if len(users) > 1 {
		return nil, fmt.Errorf("Multiple users found for username %s", username)
	}

	return nil, nil

}

// UserGetByFeverHash get a user object by feverhash, nil if not found
func (z *Service) UserGetByFeverHash(feverhash string) (*m.User, error) {

	users := []*m.User{}
	add := func() interface{} {
		u := &m.User{}
		users = append(users, u)
		return u
	}

	err := z.db.View(func(tx *bolt.Tx) error {
		return Query(bucketUser, "FeverHash", []interface{}{feverhash}, []interface{}{feverhash}, add, tx)
	})

	if err != nil {
		return nil, err
	}

	if len(users) == 1 {
		return users[0], nil
	} else if len(users) > 1 {
		return nil, fmt.Errorf("Multiple users found for feverhash %s", feverhash)
	}

	return nil, nil

}

// UserSave saves a user to the database. If the user is new, the application must check if the username is unique.
func (z *Service) UserSave(user *m.User) error {

	z.Lock()
	defer z.Unlock()
	err := z.db.Update(func(tx *bolt.Tx) error {
		_, err := Put(user, tx)
		return err
	})

	return err

}
