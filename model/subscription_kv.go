package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"sort"
	"strings"
	"time"
)

// index names
const (
	subscriptionEntity    = "Subscription"
	subscriptionIndexFeed = "Feed"
	subscriptionIndexUser = "User"
)

const (
	subscriptionID        = "ID"
	subscriptionUserID    = "UserID"
	subscriptionFeedID    = "FeedID"
	subscriptionGroupIDs  = "GroupIDs"
	subscriptionDateAdded = "DateAdded"
	subscriptionTitle     = "Title"
	subscriptionNotes     = "Notes"
	subscriptionAutoRead  = "AutoRead"
	subscriptionAutoStar  = "AutoStar"
)

var (
	subscriptionAllFields = []string{
		subscriptionID, subscriptionUserID, subscriptionFeedID, subscriptionGroupIDs, subscriptionDateAdded, subscriptionTitle, subscriptionNotes, subscriptionAutoRead, subscriptionAutoStar,
	}
	subscriptionAllIndexes = []string{
		subscriptionIndexFeed, subscriptionIndexUser,
	}
)

// Subscriptions is a collection of Subscription elements
type Subscriptions []*Subscription

func (z Subscriptions) Len() int      { return len(z) }
func (z Subscriptions) Swap(i, j int) { z[i], z[j] = z[j], z[i] }
func (z Subscriptions) Less(i, j int) bool {
	return z[i].ID < z[j].ID
}

// SortByID sort collection by ID
func (z Subscriptions) SortByID() {
	sort.Stable(z)
}

// First returns the first element in the collection
func (z Subscriptions) First() *Subscription { return z[0] }

// Reverse reverses the order of the collection
func (z Subscriptions) Reverse() {
	for left, right := 0, len(z)-1; left < right; left, right = left+1, right-1 {
		z[left], z[right] = z[right], z[left]
	}
}

// getID return the primary key of the object.
func (z *Subscription) getID() string {
	return z.ID
}

// Clear reset all fields to zero/empty
func (z *Subscription) clear() {
	z.ID = empty
	z.UserID = empty
	z.FeedID = empty
	z.GroupIDs = []string{}
	z.DateAdded = time.Time{}
	z.Title = empty
	z.Notes = empty
	z.AutoRead = false
	z.AutoStar = false

}

// Serialize serializes an object to a list of key-values.
// An optional flag, when set, will serialize all fields to the resulting map, not just the non-zero values.
func (z *Subscription) serialize(flags ...bool) Record {

	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(Record)

	if flagNoZeroCheck || z.ID != empty {
		result[subscriptionID] = z.ID
	}

	if flagNoZeroCheck || z.UserID != empty {
		result[subscriptionUserID] = z.UserID
	}

	if flagNoZeroCheck || z.FeedID != empty {
		result[subscriptionFeedID] = z.FeedID
	}

	if flagNoZeroCheck || len(z.GroupIDs) > 0 {
		result[subscriptionGroupIDs] = func(values []string) string {
			return strings.Join(values, " ")
		}(z.GroupIDs)
	}

	if flagNoZeroCheck || !z.DateAdded.IsZero() {
		result[subscriptionDateAdded] = z.DateAdded.UTC().Format(fmtTime)
	}

	if flagNoZeroCheck || z.Title != empty {
		result[subscriptionTitle] = z.Title
	}

	if flagNoZeroCheck || z.Notes != empty {
		result[subscriptionNotes] = z.Notes
	}

	if flagNoZeroCheck || z.AutoRead {
		result[subscriptionAutoRead] = func(value bool) string {
			if value {
				return "1"
			}
			return "0"
		}(z.AutoRead)
	}

	if flagNoZeroCheck || z.AutoStar {
		result[subscriptionAutoStar] = func(value bool) string {
			if value {
				return "1"
			}
			return "0"
		}(z.AutoStar)
	}

	return result

}

// Deserialize serializes an object to a list of key-values.
// An optional flag, when set, will return an error if unknown keys are contained in the values.
func (z *Subscription) deserialize(values Record, flags ...bool) error {

	flagUnknownCheck := len(flags) > 0 && flags[0]
	z.clear()

	var errors []error
	var missing []string
	var unknown []string

	z.ID = values[subscriptionID]
	if !(z.ID != empty) {
		missing = append(missing, subscriptionID)
	}

	z.UserID = values[subscriptionUserID]
	if !(z.UserID != empty) {
		missing = append(missing, subscriptionUserID)
	}

	z.FeedID = values[subscriptionFeedID]
	if !(z.FeedID != empty) {
		missing = append(missing, subscriptionFeedID)
	}

	z.GroupIDs = func(fieldName string, values map[string]string, errors []error) []string {
		var result []string
		if value, ok := values[fieldName]; ok {
			result = strings.Fields(value)
		}
		return result
	}(subscriptionGroupIDs, values, errors)

	z.DateAdded = func(fieldName string, values map[string]string, errors []error) time.Time {
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
	}(subscriptionDateAdded, values, errors)

	z.Title = values[subscriptionTitle]
	z.Notes = values[subscriptionNotes]
	z.AutoRead = func(fieldName string, values map[string]string, errors []error) bool {
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(subscriptionAutoRead, values, errors)

	z.AutoStar = func(fieldName string, values map[string]string, errors []error) bool {
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(subscriptionAutoStar, values, errors)

	if flagUnknownCheck {
		for fieldname := range values {
			if !isStringInArray(fieldname, subscriptionAllFields) {
				unknown = append(unknown, fieldname)
			}
		}
	}

	return newDeserializationError(subscriptionEntity, errors, missing, unknown)

}

// serializeIndexes returns all index records
func (z *Subscription) serializeIndexes() map[string]Record {

	result := make(map[string]Record)
	data := z.serialize(true)
	var keys []string

	keys = []string{}
	keys = append(keys, data[subscriptionFeedID])
	keys = append(keys, data[subscriptionUserID])
	result[subscriptionIndexFeed] = Record{string(kvKeyEncode(keys...)): data[subscriptionID]}

	keys = []string{}
	keys = append(keys, data[subscriptionUserID])
	keys = append(keys, data[subscriptionFeedID])
	result[subscriptionIndexUser] = Record{string(kvKeyEncode(keys...)): data[subscriptionID]}

	return result

}
