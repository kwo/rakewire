package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"fmt"
	"strconv"
)

// index names
const (
	UserEntryEntity       = "UserEntry"
	UserEntryIndexStarred = "Starred"
	UserEntryIndexUser    = "User"
)

const (
	userentryID      = "ID"
	userentryUserID  = "UserID"
	userentryRead    = "Read"
	userentryStarred = "Starred"
)

// GetID return the primary key of the object.
func (z *UserEntry) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *UserEntry) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *UserEntry) Clear() {
	z.ID = 0
	z.UserID = 0
	z.Read = false
	z.Starred = false

}

// Serialize serializes an object to a list of key-values.
func (z *UserEntry) Serialize() map[string]string {
	result := make(map[string]string)

	if z.ID != 0 {
		result[userentryID] = fmt.Sprintf("%05d", z.ID)
	}

	if z.UserID != 0 {
		result[userentryUserID] = fmt.Sprintf("%05d", z.UserID)
	}

	if z.Read {
		result[userentryRead] = fmt.Sprintf("%t", z.Read)
	}

	if z.Starred {
		result[userentryStarred] = fmt.Sprintf("%t", z.Starred)
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *UserEntry) Deserialize(values map[string]string) error {
	var errors []error

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userentryID, values, errors)

	z.UserID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userentryUserID, values, errors)

	z.Read = func(fieldName string, values map[string]string, errors []error) bool {
		result, err := strconv.ParseBool(values[fieldName])
		if err != nil {
			errors = append(errors, err)
			return false
		}
		return result
	}(userentryRead, values, errors)

	z.Starred = func(fieldName string, values map[string]string, errors []error) bool {
		result, err := strconv.ParseBool(values[fieldName])
		if err != nil {
			errors = append(errors, err)
			return false
		}
		return result
	}(userentryStarred, values, errors)

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *UserEntry) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize()

	result[UserEntryIndexStarred] = []string{

		data[userentryUserID],

		data[userentryStarred],

		data[userentryID],
	}

	result[UserEntryIndexUser] = []string{

		data[userentryUserID],

		data[userentryID],
	}

	return result
}
