package model

import (
	"fmt"
	"strings"
)

// UserByFeverHash get a user object by feverhash, nil if not found
func UserByFeverHash(feverhash string, tx Transaction) (user *User, err error) {

	data, ok := kvGetFromIndex(userEntity, userIndexFeverHash, []string{strings.ToLower(feverhash)}, tx)
	if ok {
		user = &User{}
		err = user.deserialize(data)
	}

	return

}

// UserByUsername get a user object by username, nil if not found
func UserByUsername(username string, tx Transaction) (user *User, err error) {

	data, ok := kvGetFromIndex(userEntity, userIndexUsername, []string{strings.ToLower(username)}, tx)
	if ok {
		user = &User{}
		err = user.deserialize(data)
	}

	return

}

// Save persists the user to the database
func (z *User) Save(tx Transaction) error {

	// new user, check for unique username
	if z.getID() == empty {
		indexName := userIndexUsername
		if _, ok := kvGetFromIndex(userEntity, indexName, z.indexKeys()[indexName], tx); ok {
			return fmt.Errorf("Cannot save user, username is already taken: %s", strings.ToLower(z.Username))
		}
	}

	return kvSave(userEntity, z, tx)

}
