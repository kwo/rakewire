package model

import (
	"fmt"
	"strings"
)

// UserByFeverHash get a user object by feverhash, nil if not found
func UserByFeverHash(feverhash string, tx Transaction) (user *User, err error) {

	bUser := tx.Bucket(bucketData, userEntity)
	bIndex := tx.Bucket(bucketIndex, userEntity, userIndexFeverHash)

	if record := bIndex.GetIndex(bUser, feverhash); record != nil {
		user = &User{}
		err = user.deserialize(record)
	}

	return

}

// UserByUsername get a user object by username, nil if not found
func UserByUsername(username string, tx Transaction) (user *User, err error) {

	bUser := tx.Bucket(bucketData, userEntity)
	bIndex := tx.Bucket(bucketIndex, userEntity, userIndexUsername)

	if record := bIndex.GetIndex(bUser, strings.ToLower(username)); record != nil {
		user = &User{}
		err = user.deserialize(record)
	}

	return

}

// Save persists the user to the database
func (z *User) Save(tx Transaction) error {

	// new user, check for unique username
	if z.getID() == empty {
		user, err := UserByUsername(z.Username, tx)
		if err != nil {
			return err
		} else if user != nil {
			return fmt.Errorf("Cannot save user, username is already taken: %s", strings.ToLower(z.Username))
		}
	}

	return kvSave(userEntity, z, tx)

}
