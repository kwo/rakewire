package model

/*
 *  CODE GENERATED AUTOMATICALLY WITH gokv.
 *  THIS FILE SHOULD NOT BE EDITED BY HAND.
 */

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// index names
const (
	UserFeedEntity    = "UserFeed"
	UserFeedIndexFeed = "Feed"
	UserFeedIndexUser = "User"
)

const (
	userfeedID        = "ID"
	userfeedUserID    = "UserID"
	userfeedFeedID    = "FeedID"
	userfeedGroupIDs  = "GroupIDs"
	userfeedDateAdded = "DateAdded"
	userfeedTitle     = "Title"
	userfeedNotes     = "Notes"
	userfeedAutoRead  = "AutoRead"
	userfeedAutoStar  = "AutoStar"
)

// GetID return the primary key of the object.
func (z *UserFeed) GetID() uint64 {
	return z.ID
}

// SetID sets the primary key of the object.
func (z *UserFeed) SetID(id uint64) {
	z.ID = id
}

// Clear reset all fields to zero/empty
func (z *UserFeed) Clear() {
	z.ID = 0
	z.UserID = 0
	z.FeedID = 0
	z.GroupIDs = []uint64{}
	z.DateAdded = time.Time{}
	z.Title = ""
	z.Notes = ""
	z.AutoRead = false
	z.AutoStar = false

}

// Serialize serializes an object to a list of key-values.
func (z *UserFeed) Serialize(flags ...bool) map[string]string {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != 0 {
		result[userfeedID] = fmt.Sprintf("%05d", z.ID)
	}

	if flagNoZeroCheck || z.UserID != 0 {
		result[userfeedUserID] = fmt.Sprintf("%05d", z.UserID)
	}

	if flagNoZeroCheck || z.FeedID != 0 {
		result[userfeedFeedID] = fmt.Sprintf("%05d", z.FeedID)
	}

	if flagNoZeroCheck || len(z.GroupIDs) > 0 {
		result[userfeedGroupIDs] = func(values []uint64) string {
			var buffer bytes.Buffer
			for i, value := range values {
				if i > 0 {
					buffer.WriteString(" ")
				}
				buffer.WriteString(fmt.Sprintf("%d", value))
			}
			return buffer.String()
		}(z.GroupIDs)
	}

	if flagNoZeroCheck || !z.DateAdded.IsZero() {
		result[userfeedDateAdded] = z.DateAdded.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || z.Title != "" {
		result[userfeedTitle] = z.Title
	}

	if flagNoZeroCheck || z.Notes != "" {
		result[userfeedNotes] = z.Notes
	}

	if flagNoZeroCheck || z.AutoRead {
		result[userfeedAutoRead] = func(value bool) string {
			if value {
				return "1"
			}
			return "0"
		}(z.AutoRead)
	}

	if flagNoZeroCheck || z.AutoStar {
		result[userfeedAutoStar] = func(value bool) string {
			if value {
				return "1"
			}
			return "0"
		}(z.AutoStar)
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *UserFeed) Deserialize(values map[string]string) error {
	var errors []error

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userfeedID, values, errors)

	z.UserID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userfeedUserID, values, errors)

	z.FeedID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(userfeedFeedID, values, errors)

	z.GroupIDs = func(fieldName string, values map[string]string, errors []error) []uint64 {
		var result []uint64
		if value, ok := values[fieldName]; ok {
			elements := strings.Fields(value)
			for _, element := range elements {
				value, err := strconv.ParseUint(element, 10, 64)
				if err != nil {
					errors = append(errors, err)
					break
				}
				result = append(result, value)
			}
		}
		return result
	}(userfeedGroupIDs, values, errors)

	z.DateAdded = func(fieldName string, values map[string]string, errors []error) time.Time {
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
	}(userfeedDateAdded, values, errors)

	z.Title = values[userfeedTitle]

	z.Notes = values[userfeedNotes]

	z.AutoRead = func(fieldName string, values map[string]string, errors []error) bool {
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(userfeedAutoRead, values, errors)

	z.AutoStar = func(fieldName string, values map[string]string, errors []error) bool {
		if value, ok := values[fieldName]; ok {
			return value == "1"
		}
		return false
	}(userfeedAutoStar, values, errors)

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *UserFeed) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize(true)

	result[UserFeedIndexFeed] = []string{

		data[userfeedFeedID],

		data[userfeedUserID],
	}

	result[UserFeedIndexUser] = []string{

		data[userfeedUserID],

		data[userfeedFeedID],
	}

	return result
}
