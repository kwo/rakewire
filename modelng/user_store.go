package modelng

import (
	"errors"
	"strings"
)

var (
	// ErrUsernameTaken occurs when adding a new user with a non-unique username.
	ErrUsernameTaken = errors.New("Username exists already.")
	// U groups methods for accessing users.
	U = &userStore{}
)

type userStore struct{}

func (z *userStore) Delete(id string, tx Transaction) error {
	return delete(entityUser, id, tx)
}

func (z *userStore) GetByFeverhash(feverhash string, tx Transaction) *User {
	bIndex := tx.Bucket(bucketIndex, entityUser, indexUserFeverhash)
	if value := bIndex.Get(keyEncode(feverhash)); value != nil {
		return U.Get(string(value), tx)
	}
	return nil
}

func (z *userStore) Get(id string, tx Transaction) *User {
	bData := tx.Bucket(bucketData, entityUser)
	if data := bData.Get([]byte(id)); data != nil {
		user := &User{}
		if err := user.decode(data); err == nil {
			return user
		}
	}
	return nil
}

func (z *userStore) GetByUsername(username string, tx Transaction) *User {
	bIndex := tx.Bucket(bucketIndex, entityUser, indexUserUsername)
	if value := bIndex.Get(keyEncode(strings.ToLower(username))); value != nil {
		return U.Get(string(value), tx)
	}
	return nil
}

func (z *userStore) New(username string) *User {
	return &User{
		Username: username,
	}
}

func (z *userStore) Save(user *User, tx Transaction) error {
	if u := U.GetByUsername(user.Username, tx); u != nil {
		return ErrUsernameTaken
	}
	return save(entityUser, user, tx)
}
