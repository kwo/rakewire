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
	groupEntity         = "Group"
	groupIndexUserGroup = "UserGroup"
)

const (
	groupID     = "ID"
	groupUserID = "UserID"
	groupName   = "Name"
)

var (
	groupAllFields = []string{
		groupID, groupUserID, groupName,
	}
)

// GetID return the primary key of the object.
func (z *Group) getID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *Group) setID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Group) clear() {
	z.ID = 0
	z.UserID = 0
	z.Name = ""

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *Group) serialize(flags ...bool) map[string]string {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != 0 {
		result[groupID] = fmt.Sprintf("%05d", z.ID)
	}

	if flagNoZeroCheck || z.UserID != 0 {
		result[groupUserID] = fmt.Sprintf("%05d", z.UserID)
	}

	if flagNoZeroCheck || z.Name != "" {
		result[groupName] = z.Name
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *Group) deserialize(values map[string]string, flags ...bool) error {
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
	}(groupID, values, errors)

	if !(z.ID != 0) {
		missing = append(missing, groupID)
	}

	z.UserID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(groupUserID, values, errors)

	if !(z.UserID != 0) {
		missing = append(missing, groupUserID)
	}

	z.Name = values[groupName]

	if !(z.Name != "") {
		missing = append(missing, groupName)
	}

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, groupAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}
	return newDeserializationError(groupEntity, errors, missing, unknown)
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Group) indexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.serialize(true)

	result[groupIndexUserGroup] = []string{

		data[groupUserID],

		data[groupName],
	}

	return result
}
