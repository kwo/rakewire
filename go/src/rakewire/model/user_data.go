package model

import (
	"fmt"
	"strings"
)

// UserByFeverHash get a user object by feverhash, nil if not found
func UserByFeverHash(feverhash string, tx Transaction) (user *User, err error) {

	data, ok := kvGetFromIndex(UserEntity, UserIndexFeverHash, []string{strings.ToLower(feverhash)}, tx)
	if ok {
		user = &User{}
		err = user.Deserialize(data)
	}

	return

}

// UserByUsername get a user object by username, nil if not found
func UserByUsername(username string, tx Transaction) (user *User, err error) {

	data, ok := kvGetFromIndex(UserEntity, UserIndexUsername, []string{strings.ToLower(username)}, tx)
	if ok {
		user = &User{}
		err = user.Deserialize(data)
	}

	return

}

// Save persists the user to the database
func (z *User) Save(tx Transaction) error {

	// new user, check for unique username
	if z.GetID() == 0 {
		indexName := UserIndexUsername
		if _, ok := kvGetFromIndex(UserEntity, indexName, z.IndexKeys()[indexName], tx); ok {
			return fmt.Errorf("Cannot save user, username is already taken: %s", strings.ToLower(z.Username))
		}
	}

	return kvSave(UserEntity, z, tx)

}
