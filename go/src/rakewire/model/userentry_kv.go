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

var (
	userentryAllFields = []string{
		userentryID, userentryUserID, userentryEntryID, userentryUserFeedID, userentryUpdated, userentryIsRead, userentryIsStar,
	}
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
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
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
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *UserEntry) Deserialize(values map[string]string, flags ...bool) error {
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
	}(userentryID, values, errors)

	if !(z.ID != 0) {
		missing = append(missing, userentryID)
	}

	z.UserID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userentryUserID, values, errors)

	if !(z.UserID != 0) {
		missing = append(missing, userentryUserID)
	}

	z.EntryID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userentryEntryID, values, errors)

	if !(z.EntryID != 0) {
		missing = append(missing, userentryEntryID)
	}

	z.UserFeedID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userentryUserFeedID, values, errors)

	if !(z.UserFeedID != 0) {
		missing = append(missing, userentryUserFeedID)
	}

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

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, userentryAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}
	return newDeserializationError(errors, missing, unknown)
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
