package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

// index names
const (
	entryEntity    = "Entry"
	entryIndexRead = "Read"
	entryIndexStar = "Star"
	entryIndexUser = "User"
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
	entryAllIndexes = []string{
		entryIndexRead, entryIndexStar, entryIndexUser,
	}
)

// Entries is a collection of Entry elements
type Entries []*Entry

func (z Entries) Len() int      { return len(z) }
func (z Entries) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Entries) Less(i, j int) bool {
	return z[i].ID < z[j].ID
}

// SortByID sort collection by ID
func (z Entries) SortByID() {
	sort.Stable(z)
}

// First returns the first element in the collection
func (z Entries) First() *Entry { return z[0] }

// Reverse reverses the order of the collection
func (z Entries) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

// getID return the primary key of the object.
func (z *Entry) getID() uint64 {
	return z.ID
}

// setID sets the primary key of the object.
func (z *Entry) setID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *Entry) clear() {
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
func (z *Entry) serialize(flags ...bool) Record {
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
func (z *Entry) deserialize(values Record, flags ...bool) error {
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
	return newDeserializationError(entryEntity, errors, missing, unknown)
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Entry) indexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.serialize(true)

	result[entryIndexRead] = []string{

		data[entryUserID],

		data[entryIsRead],

		data[entryUpdated],

		data[entryID],
	}

	result[entryIndexStar] = []string{

		data[entryUserID],

		data[entryIsStar],

		data[entryUpdated],

		data[entryID],
	}

	result[entryIndexUser] = []string{

		data[entryUserID],

		data[entryID],
	}

	return result
}

func newEntryID(tx Transaction) (string, error) {
	return kvNextID(entryEntity, tx)
}
