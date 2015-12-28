package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"fmt"
	"strconv"
	"time"
)

// index names
const (
	UserEntryEntity    = "UserEntry"
	UserEntryIndexRead = "Read"
	UserEntryIndexStar = "Star"
	UserEntryIndexUser = "User"
)

const (
	userentryID         = "ID"
	userentryUserID     = "UserID"
	userentryEntryID    = "EntryID"
	userentryUserFeedID = "UserFeedID"
	userentryUpdated    = "Updated"
	userentryIsRead     = "IsRead"
	userentryIsStarred  = "IsStarred"
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
	z.EntryID = 0
	z.UserFeedID = 0
	z.Updated = time.Time{}
	z.IsRead = false
	z.IsStarred = false

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

	if z.EntryID != 0 {
		result[userentryEntryID] = fmt.Sprintf("%05d", z.EntryID)
	}

	if z.UserFeedID != 0 {
		result[userentryUserFeedID] = fmt.Sprintf("%05d", z.UserFeedID)
	}

	if !z.Updated.IsZero() {
		result[userentryUpdated] = z.Updated.UTC().Format(time.RFC3339)
	}

	if z.IsRead {
		result[userentryIsRead] = fmt.Sprintf("%t", z.IsRead)
	}

	if z.IsStarred {
		result[userentryIsStarred] = fmt.Sprintf("%t", z.IsStarred)
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

	z.EntryID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userentryEntryID, values, errors)

	z.UserFeedID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userentryUserFeedID, values, errors)

	z.Updated = func(fieldName string, values map[string]string, errors []error) time.Time {
		result := time.Time{}
		if value, ok := values[fieldName]; ok {
			t, err := time.Parse(time.RFC3339, value)
			if err != nil {
				errors = append(errors, err)
			} else {
				result = t
			}
		}
		return result
	}(userentryUpdated, values, errors)

	z.IsRead = func(fieldName string, values map[string]string, errors []error) bool {
		result, err := strconv.ParseBool(values[fieldName])
		if err != nil {
			errors = append(errors, err)
			return false
		}
		return result
	}(userentryIsRead, values, errors)

	z.IsStarred = func(fieldName string, values map[string]string, errors []error) bool {
		result, err := strconv.ParseBool(values[fieldName])
		if err != nil {
			errors = append(errors, err)
			return false
		}
		return result
	}(userentryIsStarred, values, errors)

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *UserEntry) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize()

	result[UserEntryIndexRead] = []string{

		data[userentryUserID],

		data[userentryIsRead],

		data[userentryUpdated],

		data[userentryEntryID],
	}

	result[UserEntryIndexStar] = []string{

		data[userentryUserID],

		data[userentryIsStarred],

		data[userentryUpdated],

		data[userentryEntryID],
	}

	result[UserEntryIndexUser] = []string{

		data[userentryUserID],

		data[userentryID],
	}

	return result
}
