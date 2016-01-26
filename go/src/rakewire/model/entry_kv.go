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
	EntryEntity    = "Entry"
	EntryIndexRead = "Read"
	EntryIndexStar = "Star"
	EntryIndexUser = "User"
)

const (
	entryID             = "ID"
	entryUserID         = "UserID"
	entryItemID         = "ItemID"
	entrySubscriptionID = "SubscriptionID"
	entryUpdated        = "Updated"
	entryIsRead         = "IsRead"
	entryIsStar         = "IsStar"
)

var (
	entryAllFields = []string{
		entryID, entryUserID, entryItemID, entrySubscriptionID, entryUpdated, entryIsRead, entryIsStar,
	}
)

// GetID return the primary key of the object.
func (z *Entry) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *Entry) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Entry) Clear() {
	z.ID = 0
	z.UserID = 0
	z.ItemID = 0
	z.SubscriptionID = 0
	z.Updated = time.Time{}
	z.IsRead = false
	z.IsStar = false

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *Entry) Serialize(flags ...bool) map[string]string {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != 0 {
		result[entryID] = fmt.Sprintf("%05d", z.ID)
	}

	if flagNoZeroCheck || z.UserID != 0 {
		result[entryUserID] = fmt.Sprintf("%05d", z.UserID)
	}

	if flagNoZeroCheck || z.ItemID != 0 {
		result[entryItemID] = fmt.Sprintf("%05d", z.ItemID)
	}

	if flagNoZeroCheck || z.SubscriptionID != 0 {
		result[entrySubscriptionID] = fmt.Sprintf("%05d", z.SubscriptionID)
	}

	if flagNoZeroCheck || !z.Updated.IsZero() {
		result[entryUpdated] = z.Updated.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || z.IsRead {
		result[entryIsRead] = func(value bool) string {
			if value {
				return "1"
			}
			return "0"
		}(z.IsRead)
	}

	if flagNoZeroCheck || z.IsStar {
		result[entryIsStar] = func(value bool) string {
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
func (z *Entry) Deserialize(values map[string]string, flags ...bool) error {
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
	}(entryID, values, errors)

	if !(z.ID != 0) {
		missing = append(missing, entryID)
	}

	z.UserID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(entryUserID, values, errors)

	if !(z.UserID != 0) {
		missing = append(missing, entryUserID)
	}

	z.ItemID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(entryItemID, values, errors)

	if !(z.ItemID != 0) {
		missing = append(missing, entryItemID)
	}

	z.SubscriptionID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(entrySubscriptionID, values, errors)

	if !(z.SubscriptionID != 0) {
		missing = append(missing, entrySubscriptionID)
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
	}(entryUpdated, values, errors)

	z.IsRead = func(fieldName string, values map[string]string, errors []error) bool {
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(entryIsRead, values, errors)

	z.IsStar = func(fieldName string, values map[string]string, errors []error) bool {
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(entryIsStar, values, errors)

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, entryAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}
	return newDeserializationError(EntryEntity, errors, missing, unknown)
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Entry) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize(true)

	result[EntryIndexRead] = []string{

		data[entryUserID],

		data[entryIsRead],

		data[entryUpdated],

		data[entryID],
	}

	result[EntryIndexStar] = []string{

		data[entryUserID],

		data[entryIsStar],

		data[entryUpdated],

		data[entryID],
	}

	result[EntryIndexUser] = []string{

		data[entryUserID],

		data[entryID],
	}

	return result
}
