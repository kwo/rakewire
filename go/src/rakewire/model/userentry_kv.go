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
	userentryIsStar     = "IsStar"
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
	z.IsStar = false

}

// Serialize serializes an object to a list of key-values.
func (z *UserEntry) Serialize(flags ...bool) map[string]string {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != 0 {
		result[userentryID] = fmt.Sprintf("%05d", z.ID)
	}

	if flagNoZeroCheck || z.UserID != 0 {
		result[userentryUserID] = fmt.Sprintf("%05d", z.UserID)
	}

	if flagNoZeroCheck || z.EntryID != 0 {
		result[userentryEntryID] = fmt.Sprintf("%05d", z.EntryID)
	}

	if flagNoZeroCheck || z.UserFeedID != 0 {
		result[userentryUserFeedID] = fmt.Sprintf("%05d", z.UserFeedID)
	}

	if flagNoZeroCheck || !z.Updated.IsZero() {
		result[userentryUpdated] = z.Updated.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || z.IsRead {
		result[userentryIsRead] = func(value bool) string {
			if value {
				return "1"
			}
			return "0"
		}(z.IsRead)
	}

	if flagNoZeroCheck || z.IsStar {
		result[userentryIsStar] = func(value bool) string {
			if value {
				return "1"
			}
			return "0"
		}(z.IsStar)
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
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(userentryIsRead, values, errors)

	z.IsStar = func(fieldName string, values map[string]string, errors []error) bool {
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(userentryIsStar, values, errors)

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *UserEntry) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize(true)

	result[UserEntryIndexRead] = []string{

		data[userentryUserID],

		data[userentryIsRead],

		data[userentryUpdated],

		data[userentryID],
	}

	result[UserEntryIndexStar] = []string{

		data[userentryUserID],

		data[userentryIsStar],

		data[userentryUpdated],

		data[userentryID],
	}

	result[UserEntryIndexUser] = []string{

		data[userentryUserID],

		data[userentryID],
	}

	return result
}
