package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"fmt"
	"strconv"
	"strings"
)

// index names
const (
	UserEntity         = "User"
	UserIndexFeverHash = "FeverHash"
	UserIndexUsername  = "Username"
)

const (
	userID           = "ID"
	userUsername     = "Username"
	userPasswordHash = "PasswordHash"
	userFeverHash    = "FeverHash"
)

// GetID return the primary key of the object.
func (z *User) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *User) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *User) Clear() {
	z.ID = 0
	z.Username = ""
	z.PasswordHash = ""
	z.FeverHash = ""

}

// Serialize serializes an object to a list of key-values.
func (z *User) Serialize(flags ...bool) map[string]string {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != 0 {
		result[userID] = fmt.Sprintf("%05d", z.ID)
	}

	if flagNoZeroCheck || z.Username != "" {
		result[userUsername] = z.Username
	}

	if flagNoZeroCheck || z.PasswordHash != "" {
		result[userPasswordHash] = z.PasswordHash
	}

	if flagNoZeroCheck || z.FeverHash != "" {
		result[userFeverHash] = z.FeverHash
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *User) Deserialize(values map[string]string) error {
	var errors []error

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userID, values, errors)

	z.Username = values[userUsername]

	z.PasswordHash = values[userPasswordHash]

	z.FeverHash = values[userFeverHash]

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *User) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize(true)

	result[UserIndexFeverHash] = []string{

		data[userFeverHash],
	}

	result[UserIndexUsername] = []string{

		strings.ToLower(data[userUsername]),
	}

	return result
}
