package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"sort"
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
func (z *Entry) getID() string {
	return z.ID
}

// Clear reset all fields to zero/empty
func (z *Entry) clear() {
	z.ID = ""
	z.UserID = ""
	z.ItemID = ""
	z.SubscriptionID = ""
	z.Updated = time.Time{}
	z.IsRead = false
	z.IsStar = false

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *Entry) serialize(flags ...bool) Record {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != "" {
		result[entryID] = z.ID
	}

	if flagNoZeroCheck || z.UserID != "" {
		result[entryUserID] = z.UserID
	}

	if flagNoZeroCheck || z.ItemID != "" {
		result[entryItemID] = z.ItemID
	}

	if flagNoZeroCheck || z.SubscriptionID != "" {
		result[entrySubscriptionID] = z.SubscriptionID
	}

	if flagNoZeroCheck || !z.Updated.IsZero() {
		result[entryUpdated] = z.Updated.UTC().Format(fmtTime)
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

	z.ID = values[entryID]

	if !(z.ID != "") {
		missing = append(missing, entryID)
	}

	z.UserID = values[entryUserID]

	if !(z.UserID != "") {
		missing = append(missing, entryUserID)
	}

	z.ItemID = values[entryItemID]

	if !(z.ItemID != "") {
		missing = append(missing, entryItemID)
	}

	z.SubscriptionID = values[entrySubscriptionID]

	if !(z.SubscriptionID != "") {
		missing = append(missing, entrySubscriptionID)
	}

	z.Updated = func(fieldName string, values map[string]string, errors []error) time.Time {
		result := time.Time{}
		if value, ok := values[fieldName]; ok {
			t, err := time.Parse(fmtTime, value)
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

// serializeIndexes returns all index records
func (z *Entry) serializeIndexes() map[string]Record {

	result := make(map[string]Record)

	data := z.serialize(true)

	var keys []string

	keys = []string{}

	keys = append(keys, data[entryUserID])

	keys = append(keys, data[entryIsRead])

	keys = append(keys, data[entryUpdated])

	keys = append(keys, data[entryID])

	result[entryIndexRead] = Record{string(kvKeyEncode(keys...)): data[entryID]}

	keys = []string{}

	keys = append(keys, data[entryUserID])

	keys = append(keys, data[entryIsStar])

	keys = append(keys, data[entryUpdated])

	keys = append(keys, data[entryID])

	result[entryIndexStar] = Record{string(kvKeyEncode(keys...)): data[entryID]}

	keys = []string{}

	keys = append(keys, data[entryUserID])

	keys = append(keys, data[entryID])

	result[entryIndexUser] = Record{string(kvKeyEncode(keys...)): data[entryID]}

	return result
}
