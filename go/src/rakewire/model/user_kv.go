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
	userEntity         = "User"
	userIndexFeverHash = "FeverHash"
	userIndexUsername  = "Username"
)

const (
	userID           = "ID"
	userUsername     = "Username"
	userPasswordHash = "PasswordHash"
	userFeverHash    = "FeverHash"
)

var (
	userAllFields = []string{
		userID, userUsername, userPasswordHash, userFeverHash,
	}
)

// Users is a collection of User elements
type Users []*User

func (z Users) Len() int      { return len(z) }
func (z Users) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Users) Less(i, j int) bool {
	return z[i].ID < z[j].ID
}

// First returns the first element in the collection
func (z Users) First() *User { return z[0] }

// Reverse reverses the order of the collection
func (z Users) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

// GetID return the primary key of the object.
func (z *User) getID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *User) setID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *User) clear() {
	z.ID = 0
	z.Username = ""
	z.PasswordHash = ""
	z.FeverHash = ""

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *User) serialize(flags ...bool) map[string]string {
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
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *User) deserialize(values map[string]string, flags ...bool) error {
	flagUnknownCheck := len(flags) > 0 && flags[0]

	var errors []error
	var missing []string
	var unknown []string

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userID, values, errors)

	if !(z.ID != 0) {
		missing = append(missing, userID)
	}

	z.Username = values[userUsername]

	if !(z.Username != "") {
		missing = append(missing, userUsername)
	}

	z.PasswordHash = values[userPasswordHash]

	z.FeverHash = values[userFeverHash]

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, userAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}
	return newDeserializationError(userEntity, errors, missing, unknown)
}

// IndexKeys returns the keys of all indexes for this object.
func (z *User) indexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.serialize(true)

	result[userIndexFeverHash] = []string{

		data[userFeverHash],
	}

	result[userIndexUsername] = []string{

		strings.ToLower(data[userUsername]),
	}

	return result
}

// GroupByUsername groups elements in the Users collection by Username
func (z Users) GroupByUsername() map[string]*User {
	result := make(map[string]*User)
	for _, user := range z {
		result[user.Username] = user
	}
	return result
}
