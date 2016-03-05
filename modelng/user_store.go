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

func (z *userStore) Delete(tx Transaction, id string) error {
	return delete(tx, entityUser, id)
}

func (z *userStore) GetByFeverhash(tx Transaction, feverhash string) *User {
	// index User FeverHash = FeverHash : UserID
	bIndex := tx.Bucket(bucketIndex, entityUser, indexUserFeverhash)
	if value := bIndex.Get([]byte(feverhash)); value != nil {
		return z.Get(tx, string(value))
	}
	return nil
}

func (z *userStore) Get(tx Transaction, id string) *User {
	bData := tx.Bucket(bucketData, entityUser)
	if data := bData.Get([]byte(id)); data != nil {
		user := &User{}
		if err := user.decode(data); err == nil {
			return user
		}
	}
	return nil
}

func (z *userStore) GetByUsername(tx Transaction, username string) *User {
	// index User Username = Username (lowercase) : UserID
	bIndex := tx.Bucket(bucketIndex, entityUser, indexUserUsername)
	if value := bIndex.Get([]byte(strings.ToLower(username))); value != nil {
		return z.Get(tx, string(value))
	}
	return nil
}

func (z *userStore) New(username string) *User {
	return &User{
		Username: username,
	}
}

func (z *userStore) Save(tx Transaction, user *User) error {
	if u := z.GetByUsername(tx, user.Username); u != nil {
		return ErrUsernameTaken
	}
	return save(tx, entityUser, user)
}
