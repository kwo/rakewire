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
	GroupEntity         = "Group"
	GroupIndexUserGroup = "UserGroup"
)

const (
	groupID     = "ID"
	groupUserID = "UserID"
	groupName   = "Name"
)

// GetID return the primary key of the object.
func (z *Group) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *Group) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Group) Clear() {
	z.ID = 0
	z.UserID = 0
	z.Name = ""

}

// Serialize serializes an object to a list of key-values.
func (z *Group) Serialize() map[string]string {
	result := make(map[string]string)

	if z.ID != 0 {
		result[groupID] = fmt.Sprintf("%05d", z.ID)
	}

	if z.UserID != 0 {
		result[groupUserID] = fmt.Sprintf("%05d", z.UserID)
	}

	if z.Name != "" {
		result[groupName] = z.Name
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *Group) Deserialize(values map[string]string) error {
	var errors []error

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(groupID, values, errors)

	z.UserID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(groupUserID, values, errors)

	z.Name = values[groupName]

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Group) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize()

	result[GroupIndexUserGroup] = []string{

		data[groupUserID],

		data[groupName],
	}

	return result
}
