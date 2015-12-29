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
	EntryIndexGUID = "GUID"
)

const (
	entryID      = "ID"
	entryGUID    = "GUID"
	entryFeedID  = "FeedID"
	entryCreated = "Created"
	entryUpdated = "Updated"
	entryURL     = "URL"
	entryAuthor  = "Author"
	entryTitle   = "Title"
	entryContent = "Content"
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
	z.GUID = ""
	z.FeedID = 0
	z.Created = time.Time{}
	z.Updated = time.Time{}
	z.URL = ""
	z.Author = ""
	z.Title = ""
	z.Content = ""

}

// Serialize serializes an object to a list of key-values.
func (z *Entry) Serialize(flags ...bool) map[string]string {
	flagNoZeroCheck := len(flags) > 0 && flags[0]
	result := make(map[string]string)

	if flagNoZeroCheck || z.ID != 0 {
		result[entryID] = fmt.Sprintf("%05d", z.ID)
	}

	if flagNoZeroCheck || z.GUID != "" {
		result[entryGUID] = z.GUID
	}

	if flagNoZeroCheck || z.FeedID != 0 {
		result[entryFeedID] = fmt.Sprintf("%05d", z.FeedID)
	}

	if flagNoZeroCheck || !z.Created.IsZero() {
		result[entryCreated] = z.Created.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || !z.Updated.IsZero() {
		result[entryUpdated] = z.Updated.UTC().Format(time.RFC3339)
	}

	if flagNoZeroCheck || z.URL != "" {
		result[entryURL] = z.URL
	}

	if flagNoZeroCheck || z.Author != "" {
		result[entryAuthor] = z.Author
	}

	if flagNoZeroCheck || z.Title != "" {
		result[entryTitle] = z.Title
	}

	if flagNoZeroCheck || z.Content != "" {
		result[entryContent] = z.Content
	}

	return result
}

// Deserialize serializes an object to a list of key-values.
func (z *Entry) Deserialize(values map[string]string) error {
	var errors []error

	z.ID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(entryID, values, errors)

	z.GUID = values[entryGUID]

	z.FeedID = func(fieldName string, values map[string]string, errors []error) uint64 {
		result, err := strconv.ParseUint(values[fieldName], 10, 64)
		if err != nil {
			errors = append(errors, err)
			return 0
		}
		return uint64(result)
	}(entryFeedID, values, errors)

	z.Created = func(fieldName string, values map[string]string, errors []error) time.Time {
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
	}(entryCreated, values, errors)

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

	z.URL = values[entryURL]

	z.Author = values[entryAuthor]

	z.Title = values[entryTitle]

	z.Content = values[entryContent]

	if len(errors) > 0 {
		return errors[0]
	}
	return nil
}

// IndexKeys returns the keys of all indexes for this object.
func (z *Entry) IndexKeys() map[string][]string {

	result := make(map[string][]string)

	data := z.Serialize(true)

	result[EntryIndexGUID] = []string{

		data[entryFeedID],

		data[entryGUID],
	}

	return result
}
